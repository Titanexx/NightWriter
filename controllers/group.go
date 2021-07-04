package controllers

import (
	"gogomddoc/models"

	"github.com/gin-gonic/gin"
)

func GetAllGroups(c *gin.Context) {
	// user, _ := helpers.LoadUserFromContext(c)
	// docs, _ := models.GetDocsByUser(user)
	groups, _ := models.GetAllGroups()

	c.JSON(200, groups)
}

func GetGroupsByID(c *gin.Context) {
	// user, _ := helpers.LoadUserFromContext(c)
	// id, _ := strconv.Atoi(c.Param("id"))
	// Check if admin or user is in group

	c.JSON(200, "TODO")
}

func AddGroup(c *gin.Context) {
	// user, _ := helpers.LoadUserFromContext(c)
	// metadata, _ := c.Get("token-metadata")
	// user, _ := models.GetUserById(metadata.(*auth.TokenMetadata).UserId)
	// docs, _ := models.GetDocsByUser(user)

	c.JSON(200, "TODO")
}

func UpdateGroupById(c *gin.Context) {
	c.JSON(200, "TODO")
}

func DeleteGroupById(c *gin.Context) {
	c.JSON(200, "TODO")
}
