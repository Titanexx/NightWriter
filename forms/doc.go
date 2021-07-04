package forms

import (
	"gogomddoc/models"
)

type AddAccess struct {
	ID    uint            `json:"ID" binding:"required,numeric"`
	Key   string          `json:"key" binding:"required,base64"`
	Right models.DocRight `json:"right" binding:"required,numeric,min=0,max=3"`
}

type UpdateAccess struct {
	Right models.DocRight `json:"right" binding:"required,numeric,min=1,max=3"`
}

type AddDoc struct {
	Title models.Content `json:"title" binding:"required"`
	User  AddAccess      `json:"user" binding:"required"`
}

type UpdateDoc struct {
	ID    uint           `json:"ID"`
	Title models.Content `json:"title" binding:"required"`
}
