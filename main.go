package main

import (
	"gogomddoc/config"
	c "gogomddoc/controllers"
	m "gogomddoc/middleware"
	"gogomddoc/models"
	"net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDb()
	models.DbInit()
	config.ConnectRedis()

	r := gin.Default()
	api := r.Group("api/")
	api.POST("/login", c.Login)
	api.POST("/register", c.Register)
	api.GET("/refresh", c.Refresh)

	authenticated := api.Group("/")
	authenticated.Use(m.Auth())
	{
		authenticated.GET("/logout", c.Logout)
		authenticated.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"data": "hello world"})
		})

		users := authenticated.Group("/users")
		users.GET("/", c.GetAllUsers)
		users.GET("/verified", c.GetAllVerifiedUser)
		users.GET("/me", c.Me)
		users.GET("/:id", c.GetUserByID)
		users.POST("/:id", c.UpdateUserById)
		users.DELETE("/:id", c.DeleteUserById)

		groups := authenticated.Group("/groups")
		groups.GET("", c.GetAllGroups)
		groups.PUT("", c.AddGroup)
		groups.GET("/:id", c.GetGroupsByID)
		groups.POST("/:id", c.UpdateGroupById)
		groups.DELETE("/:id", c.DeleteGroupById)

		docs := authenticated.Group("/docs")
		docs.Use(m.LoadUser())
		{
			docs.GET("/", c.GetAllDocs)
			docs.PUT("/", c.AddDoc)
			doc := docs.Group("/:docID")
			doc.GET("", c.GetDocByID)
			doc.POST("", c.UpdateDocById)
			doc.DELETE("", c.DeleteDocById)
			// Rights management
			docUsers := doc.Group("/users")
			docUsers.GET("", c.GetUsersByDocID)
			docUsers.PUT("", c.AddUserToDoc)
			docUsers.POST("/:userID", c.UpdateUserToDoc)
			docUsers.DELETE("/:userID", c.DeleteUserToDoc)
			docGroups := doc.Group("/groups")
			docGroups.GET("", c.GetGroupsByDocID)
			docGroups.PUT("", c.Todo)
			docGroups.DELETE("", c.Todo)
			// attachmentsGroups := doc.Group("/attachments")
			// attachmentsGroups.GET("", c.Todo)
			// attachmentsGroups.POST("", c.AddAttachmentToDoc)
			// attachmentsGroups.DELETE("/:attachmentID", c.Todo)
			// Parts management
			docParts := doc.Group("/parts")
			docParts.GET("", c.GetPartsByDocID)
			docParts.PUT("", c.AddPartToDoc)
			docPart := docParts.Group("/:partID")
			docPart.POST("", c.UpdatePartToDoc)
			docPart.DELETE("", c.DeletePartToDoc)
			docPartAttachments := docPart.Group("/attachments")
			docPartAttachments.GET("", c.GetAttachmentsIDByDocPart)
			docPartAttachments.POST("", c.AddAttachmentToDocPart)
			docPartAttachments.DELETE("", c.DelAttachmentsByDocPart)
			docPartAttachments.GET("/:attachmentID", c.GetAttachmentByIDAndDocPart)
		}
	}

	ws := r.Group("/ws")
	authenticatedWS := ws.Group("/")
	authenticatedWS.Use(m.WsAuth())
	authenticatedWS.Use(m.LoadUser())
	{
		authenticatedWS.GET("docs/:docID/parts/:partID", c.WsHandlerPartByDocIDAndPartID)
	}

	r.Use(static.Serve("/", static.LocalFile("./front/dist", true)))
	r.NoRoute(func(c *gin.Context) {
		c.File("./front/dist/index.html")
	})

	r.Run()
}
