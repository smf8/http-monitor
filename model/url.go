package model

import "github.com/jinzhu/gorm"

type URL struct {
	gorm.Model
	UserId      uint
	Address     string
	Threshold   int
	FailedTimes int
	Requests    []Request `gorm:"foreignkey:url_id"`
}

func (url *URL) ShouldTriggerAlarm() bool {
	return url.FailedTimes >= url.Threshold
}

type Request struct {
	gorm.Model
	UrlId  uint
	Result int
}
