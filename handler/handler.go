package handler

import (
	"github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/smf8/http-monitor/store"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type Handler struct {
	st *store.Store
}

// NewHandler creates a new handler with given store instance
func NewHandler(st *store.Store) *Handler {
	return &Handler{st: st}
}

// RegisterRoutes registers routes with their corresponding handler function
// functions are defined in handler package
func (h *Handler) RegisterRoutes(v *echo.Group) {
	userGroup := v.Group("/users")
	userGroup.POST("", h.SignUp)
	userGroup.POST("/login", h.Login)
}

func extractID(c echo.Context) uint {
	e := c.Get("user").(*jwt.Token)
	claims := e.Claims.(jwt.MapClaims)
	id := uint(claims["id"].(float64))
	return id
}
