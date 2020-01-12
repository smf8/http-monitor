package monitor

import (
	"errors"
	"fmt"
	"github.com/gammazero/workerpool"
	"github.com/smf8/http-monitor/model"
	"github.com/smf8/http-monitor/store"
	"net/http"
)

type Monitor struct {
	store *store.Store
	URLs  []model.URL
	wp    *workerpool.WorkerPool
}

// NewMonitor creates a Monitor instance with 'store' and 'url'
// it also creates a worker pool of size 'workerSize'
// if 'urls' is set to nil it will be initialized with an empty slice
func NewMonitor(store *store.Store, urls []model.URL) *Monitor {
	mnt := new(Monitor)
	if urls == nil {
		mnt.URLs = make([]model.URL, 0)
	}
	mnt.URLs = urls
	mnt.store = store

	// max number of workers
	mnt.wp = workerpool.New(20)
	return mnt
}

// RemoveURL removes a URL from current list of monitor's urls
// returns error if the URL to be deleted was not found
func (mnt *Monitor) RemoveURL(url model.URL) error {
	var index = -1
	for i := range mnt.URLs {
		if mnt.URLs[i].ID == url.ID {
			index = i
		}
	}
	if index == -1 {
		return errors.New("url to be deleted was not found in the slice")
	}
	// deleting from list efficiently
	mnt.URLs[index], mnt.URLs[len(mnt.URLs)-1] = mnt.URLs[len(mnt.URLs)-1], mnt.URLs[index]
	mnt.URLs = mnt.URLs[:len(mnt.URLs)-1]
	return nil
}

// AddURL appends a slice of urls to the current list of urls
func (mnt *Monitor) AddURL(urls []model.URL) {
	mnt.URLs = append(mnt.URLs, urls...)
}

// Cancel stops all tasks of fetching urls
// it will wait for current running jobs to finish
func (mnt *Monitor) Cancel() error {
	mnt.wp.Stop()
	if !mnt.wp.Stopped() {
		return errors.New("could not stop monitor")
	}
	return nil
}

// DoURL checks a single URL's response and saves it's request into database
func (mnt *Monitor) DoURL(url model.URL) error {
	var err error
	mnt.wp.Submit(func() {
		req, err := url.SendRequest()
		if err != nil {
			fmt.Println(err, "could not make request")
			req.Result = http.StatusBadRequest
		}
		// add request to database
		if err = mnt.store.AddRequest(req); err != nil {
			fmt.Println(err, "could not save request to database")
		}
		// status code was other than 2XX
		if req.Result/100 != 2 {
			if err = mnt.store.IncrementFailed(&url); err != nil {
				fmt.Println(err, "could not increment failed times for url")
			}
		}
	})
	return err
}

// Do ranges over URLs currently inside Monitor instance
// and save each one's request inside database
// this function does not block
func (mnt *Monitor) Do() error {
	var err error
	for urlIndex := range mnt.URLs {
		url := mnt.URLs[urlIndex]
		mnt.wp.Submit(func() {
			// sending request
			req, err := url.SendRequest()
			if err != nil {
				fmt.Println(err, "could not make request")
				req = new(model.Request)
				req.UrlId = url.ID
				req.Result = http.StatusBadRequest
			}
			// add request to database
			if err = mnt.store.AddRequest(req); err != nil {
				fmt.Println(err, "could not save request to database")
			}
			// status code was other than 2XX
			if req.Result/100 != 2 {
				if err = mnt.store.IncrementFailed(&url); err != nil {
					fmt.Println(err, "could not increment failed times for url")
				}
			}
		})
	}
	mnt.wp.StopWait()
	return err
}
