package monitor

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
	"github.com/smf8/http-monitor/db"
	"github.com/smf8/http-monitor/model"
	"github.com/smf8/http-monitor/store"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var mnt *Monitor
var st *store.Store
var d *gorm.DB

func TestMain(m *testing.M) {
	setupDB()

	exitCode := m.Run()

	tearDown()

	os.Exit(exitCode)

}

func setupDB() {
	d = db.Setup("test-monitor")
	st = store.NewStore(d)
	user, _ := model.NewUser("foo", "bar")
	_ = st.AddUser(user)

	mnt = NewMonitor(st, nil)

}

func tearDown() {
	// removing file and closing database after all tests are done
	if err := d.Close(); err != nil {
		log.Error(err)
	}
	if err := os.Remove("test-monitor.db"); err != nil {
		log.Error(err)
	}
}
func TestMonitor_Do(t *testing.T) {
	tearDown()
	setupDB()
	urls := []model.URL{
		{UserId: 1, Address: "http://127.0.0.1", Threshold: 10, FailedTimes: 0},
		{UserId: 2, Address: "http://127.0.0.1", Threshold: 10, FailedTimes: 0},
	}
	st.AddURL(&urls[0])
	st.AddURL(&urls[1])
	mnt.AddURL(urls)
	mnt.Do()
	req, _ := st.GetRequestsByUrl(&urls[0])
	assert.Len(t, req, 1)
}
func TestMonitor_DoURL(t *testing.T) {
	tearDown()
	setupDB()
	url, _ := model.NewURL(1, "http://128.0.0.1", 5)
	st.AddURL(url)

	mnt.DoURL(*url)
	req, _ := st.GetRequestsByUrl(url)
	assert.Len(t, req, 1)
	assert.Equal(t, req[0].Result, 400)

	url, err := st.GetURLById(1)
	assert.NoError(t, err)

	assert.Equal(t, url.FailedTimes, 1)
}

func TestMonitor_Cancel(t *testing.T) {
	tearDown()
	setupDB()
	mnt.Do()
	res := mnt.Cancel()
	assert.NoError(t, res)
}

func TestMonitor_RemoveURL(t *testing.T) {
	urls := []model.URL{
		{UserId: 1, Address: "http://127.0.0.1", Threshold: 10, FailedTimes: 0},
		{UserId: 2, Address: "http://127.0.0.1", Threshold: 10, FailedTimes: 0},
	}
	urls[0].ID = 1
	urls[1].ID = 2
	mnt.AddURL(urls)
	err := mnt.RemoveURL(urls[0])
	assert.NoError(t, err)

	u := model.URL{}
	u.ID = 4
	err = mnt.RemoveURL(u)
	assert.Error(t, err)
}
