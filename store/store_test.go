package store

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
	"github.com/smf8/http-monitor/db"
	"github.com/smf8/http-monitor/model"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var database *gorm.DB
var st *Store
var usersList []*model.User
var urlsList []*model.URL

func TestMain(m *testing.M) {
	//initializing database
	database = db.Setup("test.db")
	st = NewStore(database)

	setup()

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

func setup() {
	usersList = make([]*model.User, 2)
	usersList[0], _ = model.NewUser("TestUser", "TestPassword")
	usersList[1], _ = model.NewUser("TestUser1", "TestPassword1")

	urlsList = make([]*model.URL, 10)
	for i := range urlsList {
		urlsList[i] = new(model.URL)
		urlsList[i].UserId = usersList[0].ID
		urlsList[i].Address = fmt.Sprintf("www.foo%d.bar", i)
		urlsList[i].Threshold = 10
	}
}

//TestUsers tests user insertion / reading
func TestUsers(t *testing.T) {
	err := st.AddUser(usersList[0])
	assert.NoError(t, err, "error adding user to database")
	_ = st.AddUser(usersList[1])
	dbUser, err := st.GetUserByUserName("TestUser")
	assert.NoError(t, err, "error reading user from database")
	assert.Equal(t, dbUser.Username, "TestUser")
	_, err = st.GetUserByUserName("invalid-username")
	assert.Error(t, err)
	users, err := st.GetAllUsers()
	assert.NoError(t, err, "error reading all users from database")
	assert.Equal(t, 2, len(users))
	// Changing usersList so that they have valid ID value from database
	usersList[0], usersList[1] = &users[0], &users[1]
}

func TestUrls(t *testing.T) {
	// URL insertion
	for i := range urlsList {
		urlsList[i].UserId = usersList[0].ID
		err := st.AddURL(urlsList[i])
		assert.NoError(t, err, "Error inserting url into database")
	}
	// URL reading
	u, err := st.GetURLById(1)
	assert.NoError(t, err, "Error reading url with id 1 from database")

	assert.Equal(t, u.Address, "www.foo0.bar", "Mismatch url in database")

	_, err = st.GetURLById(1000)
	assert.Error(t, err)
	// Updating URL

	_, err = st.GetURLsByUser(usersList[0].ID)
	assert.NoError(t, err)

	err = st.IncrementFailed(u)
	err = st.IncrementFailed(u)
	assert.NoError(t, err, "Error incrementing failed times")

	u, _ = st.GetURLById(1)
	assert.Equal(t, 2, u.FailedTimes, "Increment failed_times didn't work")

	err = st.DismissAlert(u.ID)
	assert.NoError(t, err, "Error resetting failed times in database")

	u, _ = st.GetURLById(1)
	assert.Equal(t, 0, u.FailedTimes, "Resetting failed times didn't work")
}

func TestRequests(t *testing.T) {
	// test url insertion
	for i := range urlsList {
		req := new(model.Request)
		req.Result = 300
		req.UrlId = urlsList[i/3].ID
		err := st.AddRequest(req)
		assert.NoError(t, err)
	}
	// test request retrieval
	reqs, err := st.GetRequestsByUrl(urlsList[0].ID)
	assert.NoError(t, err, "Error retrieving requests from database")
	assert.Equal(t, 3, len(reqs), "Mismatch between number of inserted and retrieved requests")

	urlsByTime, err := st.GetUserRequestsInPeriod(urlsList[0].ID, time.Now().Add(-time.Minute*3), time.Now())
	assert.NoError(t, err)
	assert.Equal(t, 3, len(urlsByTime.Requests), "error getting urls filtered by time")
}
