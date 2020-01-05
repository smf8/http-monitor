package model

import "github.com/jinzhu/gorm"

type Url struct {
	gorm.Model
	UserId      uint
	Address     string
	Threshold   int
	FailedTimes int
	Requests    []Request `gorm:"foreignkey:url_id"`
}

type Request struct {
	gorm.Model
	UrlId  uint
	Result int
}
