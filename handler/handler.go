package handler

import (
	"github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/smf8/http-monitor/monitor"
	"github.com/smf8/http-monitor/store"
)

// require validator to add "required" tag to every struct field in the package
func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type Handler struct {
	st  *store.Store
	sch *monitor.Scheduler
}

// NewHandler creates a new handler with given store instance
func NewHandler(st *store.Store, sch *monitor.Scheduler) *Handler {
	return &Handler{st: st, sch: sch}
}

func extractID(c echo.Context) uint {
	e := c.Get("user").(*jwt.Token)
	claims := e.Claims.(jwt.MapClaims)
	id := uint(claims["id"].(float64))
	return id
}
