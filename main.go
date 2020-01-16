package main

import (
	"github.com/labstack/echo/v4"
	"github.com/smf8/http-monitor/common"
	"github.com/smf8/http-monitor/db"
	"github.com/smf8/http-monitor/handler"
	"github.com/smf8/http-monitor/monitor"
	"github.com/smf8/http-monitor/store"
	"log"
	"time"
)

func main() {
	d := db.Setup("http-monitor.db")
	st := store.NewStore(d)
	mnt := monitor.NewMonitor(st, nil, 10)
	//mnt.Do()
	sch, _ := monitor.NewScheduler(mnt)
	sch.DoWithIntervals(time.Minute * 5)

	err := mnt.LoadFromDatabase()
	if err != nil {
		log.Println(err)
	}
	e := echo.New()
	v1 := e.Group("/api")
	h := handler.NewHandler(st, sch)
	h.RegisterRoutes(v1)

	e.HTTPErrorHandler = common.CustomHTTPErrorHandler
	e.Logger.Fatal(e.Start(":8080"))
}
