package controllers

import (
	"fmt"
	"gogomddoc/forms"
	"gogomddoc/helpers"
	"gogomddoc/models"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func loadDocAndUser(c *gin.Context, minRightNeededOptionnal ...models.DocRight) (doc models.Doc, user models.User, mustStop bool) {
	var minRightNeeded models.DocRight = 1
	if len(minRightNeededOptionnal) > 0 {
		minRightNeeded = minRightNeededOptionnal[0]
	}
	user = c.MustGet("user").(models.User)
	docID := helpers.GetUintIDInUrl(c, "docID")
	doc, _ = models.GetOnlyDocById(docID, user)

	if doc.Right < minRightNeeded {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "unauthorized")
		return models.Doc{}, user, true
	}

	return doc, user, false
}

func GetAllDocs(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	docs, _ := models.GetDocsByUser(user)
	docs2, _ := models.GetDocsByGroups(user.Groups)
	c.JSON(200, append(docs, docs2...))
}

func GetDocByID(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	docID := helpers.GetUintIDInUrl(c, "docID")
	if !models.HasRightByDocIDAndUser(user, docID, models.Reader) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	doc, _ := models.GetDocById(docID, user)

	c.JSON(200, doc)
}

func AddDoc(c *gin.Context) {
	var creatingDoc forms.AddDoc
	if err := c.ShouldBindJSON(&creatingDoc); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	user := c.MustGet("user").(models.User)
	doc := models.Doc{
		Title: creatingDoc.Title,
	}

	doc.Create()
	doc.AddUser(user, creatingDoc.User.Key, 3)

	c.JSON(200, doc)
}

func UpdateDocById(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	docID := helpers.GetUintIDInUrl(c, "docID")
	if !models.HasRightByDocIDAndUser(user, docID, models.Reader) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	doc, _ := models.GetDocById(docID, user)

	var updatedDoc forms.UpdateDoc
	if err := c.ShouldBindJSON(&updatedDoc); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	if doc.Title.ID == updatedDoc.Title.ID {
		doc.Title = updatedDoc.Title
		doc.Title.Update()
		c.JSON(200, "ok")
	} else {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, "Invalid json provided")
	}

}

func DeleteDocById(c *gin.Context) {
	doc, _, mustStop := loadDocAndUser(c, models.Editor)
	if mustStop {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, "unauthorized")
		return
	}

	err := doc.Delete()

	if err != nil {
		c.AbortWithStatusJSON(500, "Error during deleting")
		return
	}

	c.JSON(200, "ok")
}

func GetUsersByDocID(c *gin.Context) {
	doc, _, mustStop := loadDocAndUser(c, models.Editor)
	if mustStop {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, "unauthorized")
		return
	}

	users, _ := doc.GetUsers()
	c.JSON(200, users)
}

func AddUserToDoc(c *gin.Context) {
	doc, _, mustStop := loadDocAndUser(c, models.Editor)
	if mustStop {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, "unauthorized")
		return
	}

	var userKey forms.AddAccess
	if err := c.ShouldBindJSON(&userKey); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	userDocToAdd := models.UserDoc{UserID: userKey.ID, DocID: doc.ID, Key: userKey.Key, Right: userKey.Right}
	userDocToAdd.Create()

	c.JSON(200, 0)
}

func UpdateUserToDoc(c *gin.Context) {
	doc, _, mustStop := loadDocAndUser(c, models.Editor)
	if mustStop {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, "unauthorized")
		return
	}
	userID, _ := strconv.Atoi(c.Param("userID"))

	var newRight forms.UpdateAccess
	if err := c.ShouldBindJSON(&newRight); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	// Todo check if the deleted user access is the last editor

	userDoc, _ := doc.GetUserByID(userID)
	userDoc.Right = newRight.Right
	userDoc.Update()

	c.JSON(200, 0)
}

func DeleteUserToDoc(c *gin.Context) {
	doc, _, mustStop := loadDocAndUser(c, models.Editor)
	if mustStop {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, "unauthorized")
		return
	}

	userID, _ := strconv.Atoi(c.Param("userID"))
	userDoc, _ := doc.GetUserByID(userID)
	userDoc.Delete()

	c.JSON(200, "ok")
}

func GetGroupsByDocID(c *gin.Context) {
	doc, _, mustStop := loadDocAndUser(c, models.Editor)
	if mustStop {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, "unauthorized")
		return
	}

	groups, _ := doc.GetGroups()
	c.JSON(200, groups)
}

func GetPartsByDocID(c *gin.Context) {
	doc, _, mustStop := loadDocAndUser(c)
	if mustStop {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, "unauthorized")
		return
	}
	parts, _ := doc.GetParts()
	c.JSON(200, parts)
}

func AddPartToDoc(c *gin.Context) {
	doc, _, mustStop := loadDocAndUser(c, models.Writer)
	if mustStop {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, "unauthorized")
		return
	}

	var formPart forms.Part
	if err := c.ShouldBindJSON(&formPart); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	if formPart.Content.IV == "" || formPart.Title.IV == "" || formPart.Characteristics.IV == "" {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	part := models.Part{
		Level:           formPart.Level,
		Order:           formPart.Order,
		Title:           formPart.Title,
		Characteristics: formPart.Characteristics,
		Content:         formPart.Content,
		DocID:           doc.ID,
	}

	part.Create()
	c.JSON(200, part.ID)
}

func UpdatePartToDoc(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	docID := helpers.GetUintIDInUrl(c, "docID")
	partID := helpers.GetUintIDInUrl(c, "partID")

	if !models.HasRightByDocIDAndUser(user, docID, models.Writer) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	var formPart forms.Part
	if err := c.ShouldBindJSON(&formPart); err != nil {
		log.Print(err)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	part := models.GetPartByDocIDAndPartID(user, docID, partID)
	part.Level = formPart.Level
	part.Order = formPart.Order

	if formPart.Title.ID != 0 {
		if part.TitleID == formPart.Title.ID {
			part.Title = formPart.Title
		} else {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, "Invalid json provided")
			return
		}
	}
	if formPart.Characteristics.ID != 0 {
		if part.CharacteristicsID == formPart.Characteristics.ID {
			part.Characteristics = formPart.Characteristics
		} else {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, "Invalid json provided")
			return
		}
	}
	if formPart.Content.ID != 0 {
		if part.ContentID == formPart.Content.ID {
			part.Content = formPart.Content
		} else {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, "Invalid json provided")
			return
		}
	}

	part.Update()

	c.JSON(200, part)
}

func DeletePartToDoc(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	docID := helpers.GetUintIDInUrl(c, "docID")
	partID := helpers.GetUintIDInUrl(c, "partID")

	if !models.HasRightByDocIDAndUser(user, docID, models.Writer) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	part := models.GetPartByDocIDAndPartID(user, docID, partID)
	part.Delete()
	c.JSON(200, "ok")
}

// func AddAttachmentToDoc(c *gin.Context) {
// 	user := c.MustGet("user").(models.User)
// 	docID := helpers.GetUintIDInUrl(c, "docID")
// 	if !models.HasRightByDocIDAndUser(user, docID, models.Writer) {
// 		c.AbortWithStatusJSON(http.StatusUnauthorized, "unauthorized")
// 		return
// 	}

// 	attachmentFile, err := c.FormFile("attachment")
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, "Invalid data provided")
// 		return
// 	}

// 	attachmentContent, _ := attachmentFile.Open()
// 	attachmentByte, err := ioutil.ReadAll(attachmentContent)
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, "No file is received")
// 		return
// 	}

// 	attachment := models.Attachment{Data: attachmentByte, OwnerID: docID, OwnerType: "docs"}
// 	attachment.Create()
// 	// File saved successfully. Return proper result
// 	c.JSON(http.StatusOK, attachment.ID)
// }

func AddAttachmentToDocPart(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	docID := helpers.GetUintIDInUrl(c, "docID")
	partID := helpers.GetUintIDInUrl(c, "partID")

	if !models.HasRightByDocIDAndUser(user, docID, models.Writer) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	fmt.Printf("COUCOUCOPUICOUCOUC")

	part := models.GetPartByDocIDAndPartID(user, docID, partID)
	if part.ID == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	attachmentFile, err := c.FormFile("attachment")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, "Invalid data provided")
		return
	}

	attachmentContent, _ := attachmentFile.Open()
	attachmentByte, err := ioutil.ReadAll(attachmentContent)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "No file is received")
		return
	}

	attachment := models.Attachment{Data: attachmentByte, OwnerID: partID, OwnerType: "parts"}
	attachment.Create()
	// File saved successfully. Return proper result
	c.JSON(http.StatusOK, attachment.ID)
}

func GetAttachmentsIDByDocPart(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	docID := helpers.GetUintIDInUrl(c, "docID")
	partID := helpers.GetUintIDInUrl(c, "partID")

	if !models.HasRightByDocIDAndUser(user, docID, models.Reader) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	part := models.GetPartByDocIDAndPartID(user, docID, partID)
	if part.ID == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	attachmentIDs := models.GetAttachmentsIDByOwner(partID, "parts")

	c.JSON(http.StatusOK, attachmentIDs)
}

func GetAttachmentByIDAndDocPart(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	docID := helpers.GetUintIDInUrl(c, "docID")
	partID := helpers.GetUintIDInUrl(c, "partID")
	attachmentID := helpers.GetUintIDInUrl(c, "attachmentID")

	if !models.HasRightByDocIDAndUser(user, docID, models.Reader) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	part := models.GetPartByDocIDAndPartID(user, docID, partID)
	if part.ID == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	attachment := models.GetAttachmentByIDAndOwner(attachmentID, partID, "parts")

	c.Data(http.StatusOK, "application/octet-stream", attachment.Data)
}

func DelAttachmentsByDocPart(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	docID := helpers.GetUintIDInUrl(c, "docID")
	partID := helpers.GetUintIDInUrl(c, "partID")

	if !models.HasRightByDocIDAndUser(user, docID, models.Writer) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	part := models.GetPartByDocIDAndPartID(user, docID, partID)
	if part.ID == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	var ids []uint
	if err := c.ShouldBindJSON(&ids); err != nil {
		log.Print(err)
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	if err := models.DeleteAttachmentsByIDs(ids, partID, "parts"); err != nil {
		log.Print(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, "Bug in delete")
		return
	}

	c.JSON(http.StatusOK, "ok")
}
