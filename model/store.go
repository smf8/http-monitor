package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

// User store functions

func (s *Store) GetUserByUserName(username string) (*User, error) {
	user := new(User)
	if err := s.db.Preload("Urls").Preload("Urls.Requests").First(user, User{Username: username}).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Store) GetAllUsers() ([]User, error) {
	var users []User
	if err := s.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (s *Store) AddUser(user *User) error {
	return s.db.Create(user).Error
}

// Url store functions

func (s *Store) AddURL(url *Url) error {
	return s.db.Create(url).Error
}
func (s *Store) GetURLById(id uint) (*Url, error) {
	url := new(Url)
	if err := s.db.First(url, id).Error; err != nil {
		return nil, err
	}
	return url, nil
}
func (s *Store) GetURLsByUser(user *User) ([]Url, error) {
	var urls []Url
	if err := s.db.Model(user).Related(&urls).Error; err != nil {
		return nil, err
	}
	return urls, nil
}
func (s *Store) UpdateUrl(url *Url) error {
	return s.db.Model(url).Update(url).Error
}

//DismissAlert sets "FailedTimes" value to 0 and updates it's record in database
// https://github.com/jinzhu/gorm/issues/202#issuecomment-52582525
func (s *Store) DismissAlert(url *Url) error {
	return s.db.Model(url).Update("failed_times", 0).Error
}

func (s *Store) IncrementFailed(url *Url) error {
	url.FailedTimes += 1
	return s.UpdateUrl(url)
}

// Request store functions

func (s *Store) AddRequest(req *Request) error {
	return s.db.Create(req).Error
}
func (s *Store) GetRequestsByUrl(url *Url) ([]Request, error) {
	var requests []Request
	if err := s.db.Model(url).Related(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}
func (s *Store) GetUserRequestsInPeriod(user *User, from, to time.Time) ([]Url, error) {
	var urls []Url
	if err := s.db.Model(user).Preload("Requests", "created_at >= ? and created_at <= ?", from, to).Related(&urls).Error; err != nil {
		return nil, err
	}
	return urls, nil
}
