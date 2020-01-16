package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPasswordValidation(t *testing.T) {
	foo, err := NewUser("Foo", "Bar")
	assert.NoError(t, err, "Error creating user instance")
	assert.True(t, foo.ValidatePassword("Bar"), "Error validating password")
}
func TestHashPassword(t *testing.T) {
	_, err := HashPassword("")
	assert.Error(t, err, "error hashing error")
}

func TestUserCreation(t *testing.T) {
	_, err := NewUser("", "")
	assert.Error(t, err, "error creating user")
}

func TestURLCreation(t *testing.T) {
	url, err := NewURL(0, "google.com", 10)
	assert.NoError(t, err, "error creating url")
	assert.Equal(t, url.Address, "http://google.com")
	url, err = NewURL(0, "hppt://foo.bar", 10)
	assert.Error(t, err, "error validating url")
}

func TestAlarmTrigger(t *testing.T) {
	url, _ := NewURL(0, "google.com", 5)
	url.FailedTimes = 5
	assert.True(t, url.ShouldTriggerAlarm(), "error triggering alarm")
}

func TestURLSendRequest(t *testing.T) {
	url, _ := NewURL(0, "127.0.0.1:9999", 5)
	_, err := url.SendRequest()
	assert.Error(t, err)
	url.Address = "http://google.com"
	req, err := url.SendRequest()
	assert.NoError(t, err)
	assert.Equal(t, req.Result/100, 2)
}
