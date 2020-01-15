package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/smf8/http-monitor/common"
	"github.com/smf8/http-monitor/model"
	"net/http"
)

// Login handler function
func (h *Handler) Login(c echo.Context) error {
	req := &userAuthRequest{}
	user := &model.User{}
	if err := req.bind(c, user); err != nil {
		return err
	}
	// retrieving user from database
	u, err := h.st.GetUserByUserName(user.Username)
	if err != nil || !u.ValidatePassword(user.Password) {
		return common.NewRequestError("Invalid username or password", err, http.StatusUnauthorized)
	}
	// TODO: generate and provide user with JWT
	token, err := common.GenerateJWT(u.ID)
	if err != nil {
		return common.NewRequestError("Error generating JWT token", err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, token)
}

// SignUp user handler to handle sign up requests
func (h *Handler) SignUp(c echo.Context) error {
	req := &userAuthRequest{}
	user := &model.User{}
	if err := req.bind(c, user); err != nil {
		return err
	}
	user.Password, _ = model.HashPassword(user.Password)
	// saving user
	if err := h.st.AddUser(user); err != nil {
		return common.NewRequestError("could not save user in database", err, http.StatusInternalServerError)
	}
	// TODO: generate and provide user with JWT
	token, err := common.GenerateJWT(user.ID)
	if err != nil {
		return common.NewRequestError("Error generating JWT token", err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusCreated, token)
}
