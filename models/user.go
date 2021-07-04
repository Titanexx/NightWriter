package models

import (
	"gogomddoc/config"
	"gogomddoc/helpers/argon2"
)

func FilterUsersFields(users []User) {
	for i := range users {
		users[i].PrivateKey = ""
	}
}

func GetAllUsers() (users []User, err error) {
	err = config.DB.Find(&users).Error
	FilterUsersFields(users)
	return users, err
}

func GetAllVerifiedUser() (users []User, err error) {
	err = config.DB.Find(&users).Where("email_verified = true").Error
	FilterUsersFields(users)
	return users, err
}

func GetUsersByRole(role string) (users []User, err error) {
	err = config.DB.Where("role= ?", role).Find(&users).Error
	FilterUsersFields(users)
	return users, err
}

func GetUserById(id uint) (user User, err error) {
	err = config.DB.Preload("Groups").First(&user, id).Error
	return user, err
}

func GetUserByEmail(email string) (user User, err error) {
	err = config.DB.Preload("Groups").Where("email = ?", email).First(&user).Error
	return user, err
}

func GetUserByUsername(username string) (user User, err error) {
	err = config.DB.Preload("Groups").Where("username = ?", username).First(&user).Error
	return user, err
}

func (user *User) Create() (err error) {
	return config.DB.Model(&User{}).Create(user).Error
}

func (user *User) Update() (err error) {
	return config.DB.Save(user).Error
}

func (user *User) Delete() (err error) {
	config.DB.Model(&user).Association("Groups").Clear()
	config.DB.Model(&user).Association("Docs").Clear()
	return config.DB.Model(&User{}).Delete(user).Error
}

// SetPassword sets a new password stored as hash.
func (m *User) SetPassword(password string) error {
	hash, err := argon2.CreateHash(password, argon2.DefaultParams)
	if err != nil {
		return err
	}

	m.Password = hash
	return nil
}
