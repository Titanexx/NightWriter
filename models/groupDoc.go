package models

import "gogomddoc/config"

func (groupDoc *GroupDoc) Create() (err error) {
	return config.DB.Model(&GroupDoc{}).Create(&groupDoc).Error
}

func (groupDoc *GroupDoc) Update() (err error) {
	return config.DB.Save(&groupDoc).Error
}

func (groupDoc *GroupDoc) Delete() (err error) {
	return config.DB.Model(&GroupDoc{}).Delete(&groupDoc).Error
}
