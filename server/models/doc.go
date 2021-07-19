package models

import (
	"nightwriter/config"

	"gorm.io/gorm"
)

func GetDocsByUser(user User) (docs []Doc, err error) {
	err = config.DB.Preload("Title").Omit("Parts.*").Model(&user).Select("*").Association("Docs").Find(&docs)
	return docs, err
}

func GetDocsByGroups(groups []Group) (docs []Doc, err error) {
	err = config.DB.Preload("Title").Omit("Parts.*").Model(&groups).Select("*").Association("Docs").Find(&docs)
	return docs, err
}

func GetDocById(id uint, user User) (doc Doc, err error) {
	err = config.DB.
		Preload("Title").
		Preload("Parts").
		Preload("Parts.Title").
		Preload("Parts.Characteristics").
		Preload("Parts.Content").
		Model(&user).Select("*").Where("ID = ?", id).Association("Docs").Find(&doc)
	return doc, err
}

func GetOnlyDocById(id uint, user User) (doc Doc, err error) {
	err = config.DB.Model(&user).Select("*").Where("ID = ?", id).Association("Docs").Find(&doc)
	return doc, err
}

func IsDocAndPartExistByID(user User, docID uint, partID uint) bool {
	var count int64
	config.DB.Joins("JOIN user_docs ON user_docs.doc_id = parts.doc_id and user_docs.user_id = ?", user.ID).Model(&Part{}).Where("parts.doc_id = ? and ID = ?", docID, partID).Count(&count)
	return count == 1
}

func HasRightByDocIDAndUser(user User, docID uint, right DocRight) bool {
	var count int64
	config.DB.Model(&UserDoc{}).Where("user_id = ? and doc_id = ? and user_docs.right >= ?", user.ID, docID, right).Count(&count)
	return count == 1
}

func (doc *Doc) Create() (err error) {
	return config.DB.Model(&Doc{}).Omit("Users.*").Omit("Groups.*").Omit("Parts.*").Create(&doc).Error
}

func (doc *Doc) Update() (err error) {
	return config.DB.Session(&gorm.Session{FullSaveAssociations: true}).Omit("Users.*").Omit("Groups.*").Omit("Parts.*").Save(&doc).Error
}

func (doc *Doc) Delete() (err error) {
	config.DB.Model(&doc).Association("Users").Clear()
	config.DB.Model(&doc).Association("Groups").Clear()
	return config.DB.Model(&Doc{}).Delete(doc).Error
}

func (doc *Doc) GetUsers() (users []User, err error) {
	err = config.DB.Model(&doc).Select("*").Association("Users").Find(&users)
	FilterUsersFields(users)
	return users, err
}

func (doc *Doc) GetUserByID(userID int) (userDoc UserDoc, err error) {
	err = config.DB.Model(&userDoc).Where("user_id = ? and doc_id = ?", userID, doc.ID).First(&userDoc).Error
	return userDoc, err
}

func (doc *Doc) GetKeyByUser(user User) (key string, err error) {
	userDoc := UserDoc{}
	err = config.DB.Model(&userDoc).Select("key").Where("user_id = ? and doc_id = ?", user.ID, doc.ID).First(&key).Error
	return key, err
}

func (doc *Doc) AddUser(user User, key string, right DocRight) (err error) {
	return config.DB.Create(&UserDoc{
		UserID: user.ID,
		DocID:  doc.ID,
		Key:    key,
		Right:  right,
	}).Error
}

func (doc *Doc) DeleteUser(user User) (err error) {
	return config.DB.Model(&doc).Association("Users").Delete(&user)
}

func (doc *Doc) GetGroups() (groups []Group, err error) {
	err = config.DB.Model(&doc).Association("Groups").Find(&groups)
	FilterGroupsFields(groups)
	return groups, err
}

func (doc *Doc) GetGroupByID(group Group) (groupDoc GroupDoc, err error) {
	err = config.DB.Model(&groupDoc).Where("group_id = ? and doc_id = ?", group.ID, doc.ID).First(&groupDoc).Error
	return groupDoc, err
}

func (doc *Doc) GetKeyByGroup(group Group) (key string, err error) {
	groupDoc := GroupDoc{}
	err = config.DB.Model(&groupDoc).Select("key").Where("group_id = ? and doc_id = ?", group.ID, doc.ID).First(&key).Error
	return key, err
}

func (doc *Doc) AddGroup(group Group, key string, right DocRight) (err error) {
	return config.DB.Model(&doc).Create(&GroupDoc{
		GroupID: group.ID,
		DocID:   doc.ID,
		Key:     key,
		Right:   right,
	}).Error
}

func (doc *Doc) DeleteGroup(group Group) (err error) {
	return config.DB.Model(&doc).Association("Groups").Delete(&group)
}

func (doc *Doc) GetParts() (parts []Part, err error) {
	err = config.DB.Model(&doc).Preload("Title").Preload("Characteristics").Preload("Content").Association("Parts").Find(&parts)
	return parts, err
}

func (doc *Doc) AddPart(part Part) (err error) {
	return config.DB.Model(&doc).Association("Parts").Append(&part)
}

func (doc *Doc) GetPart(partID uint) (part Part, err error) {
	err = config.DB.Model(&doc).Where("ID = ?", partID).Association("Parts").Find(&part)
	return part, err
}

func (doc *Doc) DelPart(part Part) (err error) {
	return config.DB.Model(&doc).Association("Parts").Delete(&part)
}

func (doc *Doc) AddAttachment(attachment *Attachment) (err error) {
	return config.DB.Model(&doc).Association("Attachments").Append(&attachment)
}
