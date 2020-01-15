package main

import (
	"github.com/labstack/echo/v4"
	middleware2 "github.com/labstack/echo/v4/middleware"
	"github.com/smf8/http-monitor/common"
	"github.com/smf8/http-monitor/db"
	"github.com/smf8/http-monitor/handler"
	"github.com/smf8/http-monitor/middleware"
	"github.com/smf8/http-monitor/monitor"
	"github.com/smf8/http-monitor/store"
)

func main() {
	d := db.Setup("http-monitor.db")
	st := store.NewStore(d)
	mnt := monitor.NewMonitor(st, nil, 10)
	mnt.Do()

	e := echo.New()
	v1 := e.Group("/api")
	h := handler.NewHandler(st)
	v1.Use(middleware.JWT(common.JWTSecret))
	v1.Use(middleware2.RemoveTrailingSlash())
	h.RegisterRoutes(v1)
	middleware.AddToWhiteList("/api/users/login", "POST")

	e.HTTPErrorHandler = common.CustomHTTPErrorHandler
	e.Logger.Fatal(e.Start(":8080"))
}
