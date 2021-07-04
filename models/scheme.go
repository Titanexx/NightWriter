package models

import (
	"gogomddoc/config"
	"log"
	"time"

	"gorm.io/gorm"
)

type DocRight int

const (
	Nothing DocRight = iota
	Reader
	Writer
	Editor
)

type Group struct {
	gorm.Model
	Name       string
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
	Users      []User `gorm:"many2many:user_groups;"`
	Docs       []Doc  `gorm:"many2many:group_docs;"`
}

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"unique;type:varchar(32);not null;" copier:"must"`
	Password string `json:"-" gorm:"type:varchar(255)"`
	Name     string `json:"name" gorm:"type:varchar(32);not null" copier:"must"`
	// Picture        string `json:"picture"`
	Email         string   `json:"email" gorm:"type:varchar(100);not null;unique" copier:"must"`
	EmailVerified bool     `json:"email_verified"`
	Role          string   `json:"role" gorm:"type:varchar(32);index"`
	PublicKey     string   `json:"public_key"`
	PrivateKey    string   `json:"private_key,omitempty"`
	Groups        []Group  `gorm:"many2many:user_groups;"`
	Docs          []Doc    `gorm:"many2many:user_docs;"`
	Right         DocRight `json:"right" gorm:"->;-:migration"`
}

type UserGroup struct {
	UserID    uint `gorm:"primaryKey"`
	GroupID   uint `gorm:"primaryKey"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
	Key       string
}

type Attachment struct {
	ID        uint   `gorm:"primaryKey"`
	Data      []byte `json:"data,omitempty" gorm:"type:bytea"`
	OwnerID   uint   `json:"-"`
	OwnerType string `json:"-" gorm:"type:varchar(5);"`
}

type Content struct {
	ID      uint   ``
	Content string `json:"content,omitempty" binding:"omitempty,base64"`
	State   string `json:"state,omitempty" binding:"omitempty,base64"`
	IV      string `json:"iv" binding:"omitempty,hexadecimal"`
}

type Part struct {
	gorm.Model
	Template          bool         `json:"-" gorm:"type:boolean;default:false;index:,where:template=true"`
	Level             int          `json:"level"`
	Order             int          `json:"order"`
	TitleID           uint         `json:"-"`
	Title             Content      `json:"title" gorm:"constraint:OnDelete:CASCADE;"`
	CharacteristicsID uint         `json:"-"`
	Characteristics   Content      `json:"characteristics" gorm:"constraint:OnDelete:CASCADE;"`
	ContentID         uint         `json:"-"`
	Content           Content      `json:"content" gorm:"constraint:OnDelete:CASCADE;"`
	Attachments       []Attachment `json:"Attachements" gorm:"polymorphic:Owner;constraint:OnDelete:CASCADE"`
	DocID             uint         `json:"doc_id"`
}

type Doc struct {
	gorm.Model
	Template    bool         `json:"-" gorm:"index:,where:template=true"`
	TitleID     uint         `json:"-"`
	Title       Content      `json:"title" gorm:"constraint:OnDelete:CASCADE;"`
	Parts       []Part       `json:"parts" gorm:"constraint:OnDelete:CASCADE;"`
	Attachments []Attachment `json:"Attachements" gorm:"polymorphic:Owner;constraint:OnDelete:CASCADE"`
	Users       []User       `json:"users,omitempty" gorm:"many2many:user_docs;"`
	Groups      []Group      `json:"groups,omitempty" gorm:"many2many:group_docs;"`
	Key         string       `json:"key" gorm:"->;-:migration"`
	Right       DocRight     `json:"right" gorm:"->;-:migration"`
}

type UserDoc struct {
	UserID    uint     `json:"user_id" gorm:"primaryKey"`
	DocID     uint     `json:"doc_id" gorm:"primaryKey"`
	Key       string   `json:"key"`
	Right     DocRight `json:"right" gorm:"default:0"`
	CreatedAt time.Time
}

type GroupDoc struct {
	GroupID   uint     `json:"group_id" gorm:"primaryKey"`
	DocID     uint     `json:"doc_id" gorm:"primaryKey"`
	Key       string   `json:"key"`
	Right     DocRight `json:"right" gorm:"default:0"`
	CreatedAt time.Time
}

func DbInit() {
	err := config.DB.SetupJoinTable(&User{}, "Groups", &UserGroup{})
	if err != nil {
		panic(err)
	}
	err = config.DB.SetupJoinTable(&Group{}, "Users", &UserGroup{})
	if err != nil {
		panic(err)
	}

	err = config.DB.SetupJoinTable(&User{}, "Docs", &UserDoc{})
	if err != nil {
		panic(err)
	}
	err = config.DB.SetupJoinTable(&Doc{}, "Users", &UserDoc{})
	if err != nil {
		panic(err)
	}

	err = config.DB.SetupJoinTable(&Group{}, "Docs", &GroupDoc{})
	if err != nil {
		panic(err)
	}
	err = config.DB.SetupJoinTable(&Doc{}, "Groups", &GroupDoc{})
	if err != nil {
		panic(err)
	}

	config.DB.Debug().AutoMigrate(&User{}, &Group{}, &Content{}, &Attachment{}, &Part{}, &Doc{})
	log.Print("DB is initialized.")
}
