package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
	"github.com/smf8/http-monitor/model"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var database *gorm.DB
var store *model.Store

func TestMain(m *testing.M) {
	//initializing database
	database = Setup("test.db")
	database.LogMode(false)
	store = model.NewStore(database)

	returnCode := m.Run()
	// removing file and closing database after all tests are done
	if err := database.Close(); err != nil {
		log.Error(err)
	}
	if err := os.Remove("test.db"); err != nil {
		log.Error(err)
	}
	os.Exit(returnCode)
}

//TestUsers tests user insertion / reading
func TestUsers(t *testing.T) {
	user1, err := model.NewUser("TestUser", "TestPassword")
	assert.NoError(t, err, "error creating user instance")
	err = store.AddUser(user1)
	assert.NoError(t, nil, "error adding user to database")

	user1.Username = "TestUser1"
	user1.ID = 0
	_ = store.AddUser(user1)

	dbUser, err := store.GetUserByUserName("TestUser")
	assert.NoError(t, err, "error reading user from database")
	assert.Equal(t, dbUser.Username, "TestUser")

	users, err := store.GetAllUsers()
	assert.NoError(t, err, "error reading all users from database")
	assert.Equal(t, 2, len(users))
}

func TestUrls(t *testing.T) {
	urls := make([]*model.Url, 7)
	user, _ := store.GetAllUsers()
	// Url insertion
	for i := range urls {
		urls[i] = new(model.Url)
		urls[i].UserId = user[0].ID
		urls[i].Address = fmt.Sprintf("www.foo%d.bar", i)
		urls[i].Threshold = 10
		err := store.AddURL(urls[i])

		assert.NoError(t, err, "Error inserting url into database")
	}
	// Url reading
	u, err := store.GetURLById(1)
	assert.NoError(t, err, "Error reading url with id 1 from database")

	assert.Equal(t, u.Address, "www.foo0.bar", "Mismatch url in database")

	// Updating Url

	err = store.IncrementFailed(u)
	err = store.IncrementFailed(u)
	assert.NoError(t, err, "Error incrementing failed times")

	u, _ = store.GetURLById(1)
	assert.Equal(t, 2, u.FailedTimes, "Increment failed_times didn't work")

	err = store.DismissAlert(u)
	assert.NoError(t, err, "Error resetting failed times in database")

	u, _ = store.GetURLById(1)
	assert.Equal(t, 0, u.FailedTimes, "Resetting failed times didn't work")
}

func TestRequests(t *testing.T) {
	user, _ := store.GetAllUsers()
	urls, _ := store.GetURLsByUser(&user[0])

	// test url insertion
	for i := range urls {
		req := new(model.Request)
		req.Result = 300
		req.UrlId = urls[i/3].ID
		err := store.AddRequest(req)
		assert.NoError(t, err)
	}
	// test request retrieval
	reqs, err := store.GetRequestsByUrl(&urls[0])
	assert.NoError(t, err, "Error retrieving requests from database")
	assert.Equal(t, 3, len(reqs), "Mismatch between number of inserted and retrieved requests")
}
