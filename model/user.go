package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique_index;not null"`
	Password string `gorm:"not null"`
	Urls     []Url  `gorm:"foreignkey:user_id"`
}
