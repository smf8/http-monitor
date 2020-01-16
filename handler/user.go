package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/smf8/http-monitor/common"
	"github.com/smf8/http-monitor/model"
	"net/http"
)

// Login handler function, if login is successful a response with JWT and username is returned followed by a 200 status code
// json request format:
//
//{
//	"username": "foo" [alpha numeric, len > 3],
//	"password": "bar1" [len > 3]
//}
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
	return c.JSON(http.StatusOK, NewResponseData(NewUserResponse(u)))
}

// SignUp user handler to handle sign up request, if successful it returns JWT with username followed by 201 status code
// json request format:
//
//{
//	"username": "foo" [alpha numeric, len > 3],
//	"password": "bar1" [len > 3]
//}
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

	return c.JSON(http.StatusCreated, NewResponseData(NewUserResponse(user)))
}

// FetchAlerts retrieves all alerts for the user, returns a list of urls with alert
func (h *Handler) FetchAlerts(c echo.Context) error {
	userID := extractID(c)
	alerts, err := h.st.FetchAlerts(userID)
	if err != nil {
		return common.NewRequestError("coult not get alerts from database", err, http.StatusBadRequest)
	}
	return c.JSON(http.StatusOK, NewResponseData(alerts))
}
