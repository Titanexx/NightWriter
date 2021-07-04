package models

import "gogomddoc/config"

func FilterGroupsFields(groups []Group) {
	for i := range groups {
		groups[i].PrivateKey = ""
	}
}

func GetAllGroups() (groups []Group, err error) {
	err = config.DB.Find(&groups).Error
	FilterGroupsFields(groups)
	return groups, err
}

func GetGroupsByUser(user User) (groups []Group, err error) {
	err = config.DB.Model(&user).Association("Groups").Find(&groups)
	FilterGroupsFields(groups)
	return groups, err
}

func (group *Group) Create() (err error) {
	return config.DB.Model(&Doc{}).Create(group).Error
}

func (group *Group) Update() (err error) {
	return config.DB.Save(group).Error
}

func (group *Group) Delete() (err error) {
	config.DB.Model(&group).Association("Users").Clear()
	config.DB.Model(&group).Association("Docs").Clear()
	return config.DB.Model(&Group{}).Delete(group).Error
}
