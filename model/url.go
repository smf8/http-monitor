package model

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
	"net/http"
	"strings"
)

type URL struct {
	gorm.Model
	UserId      uint   `gorm:"unique_index:index_addr_user"` // for preventing url duplication for a single user
	Address     string `gorm:"unique_index:index_addr_user"`
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
func NewURL(userID uint, address string, threshold int) (*URL, error) {
	url := new(URL)
	url.UserId = userID
	url.Threshold = threshold
	url.FailedTimes = 0

	isValid := govalidator.IsURL(address)
	if !strings.HasPrefix("http://", address) {
		address = "http://" + address
	}
	if isValid {
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

// SendRequest sends a HTTP GET request to the url
// returns a *Request with result status code
func (url *URL) SendRequest() (*Request, error) {
	resp, err := http.Get(url.Address)
	req := new(Request)
	req.UrlId = url.ID
	if err != nil {
		return req, err
	}
	req.Result = resp.StatusCode
	return req, nil
}
