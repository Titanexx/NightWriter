package controllers

import (
	"log"
	"net/http"

	"gogomddoc/config"
	"gogomddoc/forms"
	"gogomddoc/helpers"
	"gogomddoc/helpers/argon2"
	"gogomddoc/middleware/auth"
	"gogomddoc/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

func Login(c *gin.Context) {
	var userLogin forms.Login
	if err := c.ShouldBindJSON(&userLogin); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		c.Abort()
		return
	}

	u, err := models.GetUserByUsername(userLogin.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Bad login informations")
		c.Abort()
		return
	}

	if !u.EmailVerified {
		c.JSON(http.StatusUnauthorized, "Bad login informations")
		c.Abort()
		return
	}

	match, err := argon2.ComparePasswordAndHash(userLogin.Password, u.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Bad login informations")
		c.Abort()
		return
	}
	if !match {
		c.JSON(http.StatusUnauthorized, "Bad login informations")
		c.Abort()
		return
	}

	ts, err := auth.CreateTokens(u.ID, u.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Bad login informations")
		c.Abort()
		return
	}

	saveErr := config.Redis.Auth.CreateAuth(u.ID, ts)
	if saveErr != nil {
		c.JSON(http.StatusUnauthorized, "Bad login informations")
		c.Abort()
		return
	}

	tokens := map[string]string{
		"access_token": ts.AccessToken,
	}
	c.SetSameSite(http.SameSiteStrictMode) //TODO check if prod
	// c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("auth_token", ts.RefreshToken, int(ts.RtExpires), "/api/refresh", "", false, true)
	c.SetCookie("auth_token", ts.WsToken, int(ts.WsExpires), "/ws", "", false, true)

	// http.SetCookie(c.Writer, &http.Cookie{
	// 	Name:  "auth_token",
	// 	Value: url.QueryEscape(ts.RefreshToken),
	// 	// MaxAge:   int(ts.RtExpires),
	// 	Path:     "/api/refresh",
	// 	SameSite: http.SameSiteNoneMode,
	// 	Secure:   false,
	// 	HttpOnly: true,
	// })
	// http.SetCookie(c.Writer, &http.Cookie{
	// 	Name:  "auth_token",
	// 	Value: url.QueryEscape(ts.WsToken),
	// 	// MaxAge:   int(ts.WsExpires),
	// 	Path:     "/ws",
	// 	SameSite: http.SameSiteNoneMode,
	// 	Secure:   false,
	// 	HttpOnly: true,
	// })

	c.JSON(http.StatusOK, tokens)
}

func Logout(c *gin.Context) {
	metadata, _ := auth.ExtractMetadataFromAccess(c.Request)
	if metadata != nil {
		deleteErr := config.Redis.Auth.DeleteTokens(metadata)
		if deleteErr != nil {
			log.Print(deleteErr)
			c.JSON(http.StatusBadRequest, deleteErr.Error())
			c.Abort()
			return
		}
	}
	c.JSON(http.StatusOK, "Successfully logged out")
}

func Refresh(c *gin.Context) {
	//verify the token
	metadata, err := auth.ExtractMetadataFromRefresh(c.Request)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusUnauthorized, "unauthorized")
		c.Abort()
		return
	}

	_, err = config.Redis.Auth.FetchUuid(metadata.Uuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		c.Abort()
		return
	}

	// Delete the old ones
	// TODO remove the old access token
	delErr := config.Redis.Auth.DeleteRefresh(metadata.Uuid)
	if delErr != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		c.Abort()
		return
	}

	u, err := models.GetUserById(metadata.UserId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		c.Abort()
		return
	}

	ts, err := auth.CreateTokens(u.ID, u.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		c.Abort()
		return
	}

	saveErr := config.Redis.Auth.CreateAuth(u.ID, ts)
	if saveErr != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		c.Abort()
		return
	}

	tokens := map[string]string{
		"access_token": ts.AccessToken,
	}
	c.SetSameSite(http.SameSiteStrictMode) //TODO check if prod
	c.SetCookie("refresh_token", ts.RefreshToken, int(ts.RtExpires), "/refresh", "localhost", false, true)
	c.SetCookie("auth_token", ts.WsToken, int(ts.WsExpires), "/ws", "", false, true)
	c.JSON(http.StatusOK, tokens)
}

func Register(c *gin.Context) {
	// var u user.User
	var registeringUser forms.Register
	if err := c.ShouldBindJSON(&registeringUser); err != nil {
		log.Print(err)
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	var u models.User
	if err := copier.Copy(&u, &registeringUser); err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, "Server error")
		return
	}
	u.SetPassword(registeringUser.Password)

	admins, err := models.GetUsersByRole("admin")
	if err != nil {
		helpers.ServerError(c, err)
		return
	}
	if len(admins) != 0 {
		u.EmailVerified = false
	} else {
		u.EmailVerified = true
		u.Role = "admin"
	}

	if err := u.Create(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "The user already exist")
		return
	}

	c.JSON(http.StatusCreated, "User is created")
}
