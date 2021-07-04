package models

import (
	"gogomddoc/config"

	"gorm.io/gorm"
)

func GetPartByDocIDAndPartID(user User, docID uint, partID uint) (part Part) {
	config.DB.Debug().
		Joins("JOIN user_docs ON user_docs.doc_id = parts.doc_id and user_docs.user_id = ?", user.ID).
		Model(&Part{}).
		Where("parts.doc_id = ? and ID = ?", docID, partID).
		Find(&part)
	return part
}

func (part *Part) Create() (err error) {
	return config.DB.Session(&gorm.Session{FullSaveAssociations: true}).Model(&Part{}).Create(&part).Error
}

func (part *Part) Update() (err error) {
	return config.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(&part).Error
}

func (part *Part) Delete() (err error) {
	return config.DB.Model(&Part{}).Delete(part).Error
}
