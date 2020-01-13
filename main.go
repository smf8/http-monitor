package main

import (
	"fmt"
	"github.com/smf8/http-monitor/db"
	"github.com/smf8/http-monitor/model"
	"github.com/smf8/http-monitor/monitor"
	"github.com/smf8/http-monitor/store"
)

func main() {
	d := db.Setup("http-monitor.db")
	st := store.NewStore(d)
	mnt := monitor.NewMonitor(st, nil)
	urls := make([]model.URL, 3)
	mmd, err := model.NewUser("smf", "1234")
	if err != nil {
		panic(err)
	}
	if err = st.AddUser(mmd); err != nil {
		panic(err)
	}
	for i := range urls {
		addr := fmt.Sprintf("google.com")
		url, err := model.NewURL(mmd.ID, addr, 10)
		if err != nil {
			panic(err)
		}
		urls[i] = *url
		go func() {
			req, _ := url.SendRequest()
			if req != nil {
				fmt.Println(req.Result)
			}
		}()
	}
	mnt.AddURL(urls)
	mnt.Do()
}
