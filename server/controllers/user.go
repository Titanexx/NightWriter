package controllers

import (
	"nightwriter/middleware/auth"
	"nightwriter/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Me(c *gin.Context) {
	metadata, _ := c.Get("token-metadata")
	u, _ := models.GetUserById(metadata.(*auth.TokenMetadata).UserId)
	c.JSON(200, u)
}

func GetAllUsers(c *gin.Context) {
	u, _ := models.GetAllUsers()
	c.JSON(200, u)
}

func GetAllVerifiedUser(c *gin.Context) {
	u, _ := models.GetAllVerifiedUser()
	c.JSON(200, u)
}

func GetUserByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	u, _ := models.GetUserById(uint(id))
	c.JSON(200, u)
}

func UpdateUserById(c *gin.Context) {
	c.JSON(200, "TODO")
}

func DeleteUserById(c *gin.Context) {
	c.JSON(200, "TODO")
}
