package store

import (
	"github.com/jinzhu/gorm"
	"github.com/smf8/http-monitor/model"
	"time"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

// GetUserByUserName retrieves user from database based on it's ID
// this method loads user's URLs and Requests lists
// returns error if user was not found
func (s *Store) GetUserByUserName(username string) (*model.User, error) {
	user := new(model.User)
	if err := s.db.Preload("Urls").Preload("Urls.Requests").First(user, model.User{Username: username}).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// GetAllUsers retrieves all users from database
func (s *Store) GetAllUsers() ([]model.User, error) {
	var users []model.User
	if err := s.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// AddUser add's a user to the database
func (s *Store) AddUser(user *model.User) error {
	return s.db.Create(user).Error
}

// AddURL add's a url to the database
func (s *Store) AddURL(url *model.URL) error {
	return s.db.Create(url).Error
}

// GetURLById retrieves a URL from database based on it's ID
// returns error if an URL was not fount
func (s *Store) GetURLById(id uint) (*model.URL, error) {
	url := new(model.URL)
	if err := s.db.First(url, id).Error; err != nil {
		return nil, err
	}
	return url, nil
}

// GetURLByUser retrieves urls for this user
// returns error if nothing was found
func (s *Store) GetURLsByUser(user *model.User) ([]model.URL, error) {
	var urls []model.URL
	if err := s.db.Model(user).Related(&urls).Error; err != nil {
		return nil, err
	}
	return urls, nil
}

//UpdateURL updates a URL to it's new value
func (s *Store) UpdateURL(url *model.URL) error {
	return s.db.Model(url).Update(url).Error
}

//DismissAlert sets "FailedTimes" value to 0 and updates it's record in database
// https://github.com/jinzhu/gorm/issues/202#issuecomment-52582525
func (s *Store) DismissAlert(url *model.URL) error {
	return s.db.Model(url).Update("failed_times", 0).Error
}

//IncrementFailed increments failed_times of a URL
func (s *Store) IncrementFailed(url *model.URL) error {
	url.FailedTimes += 1
	return s.UpdateURL(url)
}

// AddRequest adds a request to database
func (s *Store) AddRequest(req *model.Request) error {
	return s.db.Create(req).Error
}

// GetRequestByUrl retrieves all requests for this url
func (s *Store) GetRequestsByUrl(url *model.URL) ([]model.Request, error) {
	var requests []model.Request
	if err := s.db.Model(url).Related(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}

// GetUserRequestsInPeriod retrieves requests between 2 time intervals
func (s *Store) GetUserRequestsInPeriod(user *model.User, from, to time.Time) ([]model.URL, error) {
	var urls []model.URL
	if err := s.db.Model(user).Preload("Requests", "created_at >= ? and created_at <= ?", from, to).Related(&urls).Error; err != nil {
		return nil, err
	}
	return urls, nil
}
