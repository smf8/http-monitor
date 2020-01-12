package model

import (
	"errors"
	"github.com/jinzhu/gorm"
	"regexp"
)

type URL struct {
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

//NewURL creates a URL instance if it's address is a valid URL address
func NewURL(userID uint, address string, threshold, failedTimes int) (*URL, error) {
	url := new(URL)
	url.UserId = userID
	url.Threshold = threshold
	url.FailedTimes = failedTimes

	reg, err := regexp.Compile(`^(?:http(s)?:\/\/)?[\w.-]+(?:\.[\w\.-]+)+[\w\-\._~:/?#[\]@!\$&'\(\)\*\+,;=.]+$`)
	if err != nil {
		return nil, err
	}
	if reg.MatchString(address) {
		//valid URL address
		url.Address = address
		return url, nil
	}
	return nil, errors.New("not a valid URL address")
}

// ShouldTriggerAlarm checks if current url's failed times is greater than it's threshold
//
// Use this function to check alarm and trigger an alarm with other functions
func (url *URL) ShouldTriggerAlarm() bool {
	return url.FailedTimes >= url.Threshold
}
