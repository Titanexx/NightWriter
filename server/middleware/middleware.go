package middleware

import (
	"net/http"

	"nightwriter/config"
	"nightwriter/middleware/auth"
	"nightwriter/models"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		metadata, err := auth.ExtractMetadataFromAccess(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}

		// Fetch if is the accessuuid is a good one.
		// Todo check the signature
		_, err = config.Redis.Auth.FetchUuid(metadata.Uuid)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}

		c.Set("token-metadata", metadata)

		c.Next()
	}
}

func WsAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		metadata, err := auth.ExtractMetadataFromWs(c.Request)
		if err != nil {
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

		c.Set("token-metadata", metadata)

		c.Next()
	}
}

func LoadUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		metadata, _ := c.Get("token-metadata")
		user, err := models.GetUserById(metadata.(*auth.TokenMetadata).UserId)

		if err != nil {
			c.JSON(http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}

		c.Set("user", user)
	}
}
