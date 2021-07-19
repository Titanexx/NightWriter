package forms

import (
	"nightwriter/models"
)

type Part struct {
	Level           int            `json:"level" binding:"omitempty,numeric"`
	Order           int            `json:"order" binding:"omitempty,numeric"`
	Title           models.Content `json:"title" binding:"omitempty"`
	Characteristics models.Content `json:"characteristics" binding:"omitempty"`
	Content         models.Content `json:"content" binding:"omitempty"`
}
