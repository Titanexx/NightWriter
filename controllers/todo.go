package controllers

import (
	"github.com/gin-gonic/gin"
)

func Todo(c *gin.Context) {
	c.JSON(200, "todo")
}
