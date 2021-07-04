package helpers

import (
	"gogomddoc/middleware/auth"
	"gogomddoc/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ServerError(c *gin.Context, err error) {
	log.Print(err)
	c.JSON(http.StatusInternalServerError, "Server error")
}

func LoadUserFromContext(c *gin.Context) (user models.User, err error) {
	metadata, _ := c.Get("token-metadata")
	user, err = models.GetUserById(metadata.(*auth.TokenMetadata).UserId)

	return user, err
}

func GetUintIDInUrl(c *gin.Context, name string) (id uint) {
	_id, _ := strconv.Atoi(c.Param(name))
	return uint(_id)
}

func GetIntIDInUrl(c *gin.Context, name string) (id int) {
	id, _ = strconv.Atoi(c.Param(name))
	return id
}
