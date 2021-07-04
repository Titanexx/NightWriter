package models

import "gogomddoc/config"

func (userDoc *UserDoc) Create() (err error) {
	return config.DB.Model(&UserDoc{}).Create(&userDoc).Error
}

func (userDoc *UserDoc) Update() (err error) {
	return config.DB.Save(&userDoc).Error
}

func (userDoc *UserDoc) Delete() (err error) {
	return config.DB.Model(&UserDoc{}).Delete(&userDoc).Error
}
