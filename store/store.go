package store

import (
	"errors"
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
	// remove pre loading in the future if necessary
	if err := s.db.Preload("Urls").Preload("Urls.Requests").First(user, model.User{Username: username}).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserById retrieves a user from database with given id
// returns error if user was not found
func (s *Store) GetUserByID(id uint) (*model.User, error) {
	usr := &model.User{}
	usr.ID = id
	if err := s.db.Model(usr).Preload("Urls").Find(usr).Error; err != nil {
		return nil, err
	}
	return usr, nil
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

func (s *Store) GetAllURLs() ([]model.URL, error) {
	var urls []model.URL
	if err := s.db.Model(&model.URL{}).Find(&urls).Error; err != nil {
		return nil, err
	}
	return urls, nil
}

// GetURLById retrieves a URL from database based on it's ID
// returns error if an URL was not fount
func (s *Store) GetURLById(id uint) (*model.URL, error) {
	url := new(model.URL)
	if err := s.db.Preload("Requests").First(url, id).Error; err != nil {
		return nil, err
	}
	return url, nil
}

// GetURLByUser retrieves urls for this user
// returns error if nothing was found
func (s *Store) GetURLsByUser(userID uint) ([]model.URL, error) {
	var urls []model.URL
	if err := s.db.Model(&model.URL{}).Where("user_id == ?", userID).Find(&urls).Error; err != nil {
		return nil, err
	}
	return urls, nil
}

//UpdateURL updates a URL to it's new value
func (s *Store) UpdateURL(url *model.URL) error {
	return s.db.Model(url).Update(url).Error
}

// DeleteURL deletes a url with it's requests from database
// returns an error if url was not found
func (s *Store) DeleteURL(urlID uint) error {
	url := &model.URL{}
	url.ID = urlID
	// for hard deleting user s.db.Unscoped()
	q := s.db.Model(url).Preload("Requests").Delete(&model.Request{}, "url_id == ?", urlID).Delete(url)
	if q.Error != nil {
		return q.Error
	}
	if q.RowsAffected == 0 {
		return errors.New("no rows found to delete at delete url")
	}
	return nil
}

//DismissAlert sets "FailedTimes" value to 0 and updates it's record in database
// https://github.com/jinzhu/gorm/issues/202#issuecomment-52582525
func (s *Store) DismissAlert(urlID uint) error {
	url := &model.URL{}
	url.ID = urlID
	return s.db.Model(url).Update("failed_times", 0).Error
}

// FetchAlerts retrieves urls which "failed_times" is greater than it's "threshold" for given userID
// TODO: write tests for this function
func (s *Store) FetchAlerts(userID uint) ([]model.URL, error) {
	var urls []model.URL
	if err := s.db.Model(&model.URL{}).Where("user_id == ? and failed_times >= threshold", userID).Find(urls).Error; err != nil {
		return nil, err
	}
	return urls, nil
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
func (s *Store) GetRequestsByUrl(urlID uint) ([]model.Request, error) {
	var requests []model.Request
	if err := s.db.Model(&model.Request{UrlId: urlID}).Where("url_id == ?", urlID).Find(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}

// GetUserRequestsInPeriod retrieves requests between 2 time intervals
func (s *Store) GetUserRequestsInPeriod(urlID uint, from, to time.Time) (*model.URL, error) {
	url := &model.URL{}
	url.ID = urlID
	if err := s.db.Model(url).Preload("Requests", "created_at >= ? and created_at <= ?", from, to).First(url).Error; err != nil {
		return nil, err
	}
	return url, nil
}
