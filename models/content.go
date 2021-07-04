package models

import "gogomddoc/config"

func (content *Content) Create() (err error) {
	return config.DB.Model(&Content{}).Create(&content).Error
}

func (content *Content) Update() (err error) {
	return config.DB.Save(&content).Error
}

func (content *Content) Delete() (err error) {
	return config.DB.Model(&Content{}).Delete(&content).Error
}
