package models

import (
	"gogomddoc/config"
)

// func GetAttachmentByDocIDAndAttachmentID(user User, docID uint, attachmentID uint) (attachment Attachment) {
// 	config.DB.Debug().
// 		Joins("JOIN user_docs ON user_docs.doc_id = parts.doc_id and user_docs.user_id = ?", user.ID).
// 		Model(&Attachment{}).
// 		Where("attachments.doc_id = ? and ID = ?", docID, attachmentID).
// 		Find(&attachment)
// 	return attachment
// }

func GetAttachmentsIDByOwner(ownerID uint, ownerType string) (attachments []Attachment) {
	config.DB.Debug().Model(&Attachment{}).Select("ID").Where("owner_id = ? and owner_type = ?", ownerID, ownerType).Find(&attachments)
	return attachments
}

func GetAttachmentByIDAndOwner(attachmentID uint, ownerID uint, ownerType string) (attachment Attachment) {
	config.DB.Debug().Model(&Attachment{}).Where("ID = ? and owner_id = ? and owner_type = ?", attachmentID, ownerID, ownerType).Find(&attachment)
	return attachment
}

func DeleteAttachmentsByIDs(ids []uint, ownerID uint, ownerType string) (err error) {
	err = config.DB.Debug().Where("ID in ? and owner_id = ? and owner_type = ?", ids, ownerID, ownerType).Delete(&Attachment{}).Error
	return err
}

func (attachment *Attachment) Create() (err error) {
	return config.DB.Model(&Attachment{}).Create(&attachment).Error
}

func (attachment *Attachment) Delete() (err error) {
	return config.DB.Model(&Attachment{}).Delete(attachment).Error
}
