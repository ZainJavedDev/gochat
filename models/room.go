package models

import "github.com/jinzhu/gorm"

type Room struct {
	gorm.Model
	Name   string `gorm:"unique;not null"`
	UserId uint   `gorm:"not null"`

	User User `gorm:"foreignkey:UserId" json:"-"`
}
