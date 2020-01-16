package handler

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
	"github.com/smf8/http-monitor/common"
	"github.com/smf8/http-monitor/model"
	"net/http"
	"strconv"
	"time"
)

type userAuthRequest struct {
	Username string `valid:"stringlength(4|32), alphanum" json:"username"`
	Password string `valid:"stringlength(4|32)" json:"password"`
}

// binding user auth request with model.User instance
func (r *userAuthRequest) bind(c echo.Context, user *model.User) error {
	if err := c.Bind(r); err != nil {
		return common.NewRequestError("error binding user request", err, http.StatusBadRequest)
	}
	if _, err := govalidator.ValidateStruct(r); err != nil {
		e := common.NewValidationError(err, "Error validating sign-up request")
		return e
	}
	user.Username = r.Username
	user.Password = r.Password
	return nil
}

type urlCreateRequest struct {
	Address   string `json:"address" valid:"url"`
	Threshold int    `json:"threshold" valid:"int"`
}

func (r *urlCreateRequest) bind(c echo.Context, url *model.URL) error {
	if err := c.Bind(r); err != nil {
		return common.NewRequestError("error binding url create request, check json structure and try again", err, http.StatusBadRequest)
	}
	if _, err := govalidator.ValidateStruct(r); err != nil {
		e := common.NewValidationError(err, "Error validating create url request")
		return e
	}
	url.Address = r.Address
	url.Threshold = r.Threshold
	url.FailedTimes = 0
	return nil
}

type alertDismissRequest struct {
	Address string `json:"address" valid:"url"`
}

func (r *alertDismissRequest) bind(c echo.Context, url *model.URL) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if _, err := govalidator.ValidateStruct(r); err != nil {
		e := common.NewValidationError(err, "Error validating alert dismiss request")
		return e
	}
	url.Address = r.Address
	return nil
}

type urlStatusRequest struct {
	Address  string    `json:"address" valid:"url"`
	FromTime time.Time `valid:"optional, time~Provide time as unix timestamp" json:"from_time, omitempty"`
	ToTime   time.Time `valid:"optional, time~Provide time as unix timestamp" json:"to_time, omitempty"`
}

func (r *urlStatusRequest) bind(c echo.Context, url *model.URL) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	govalidator.CustomTypeTagMap.Set("time", govalidator.CustomTypeValidator(func(i interface{}, context interface{}) bool {
		if _, ok := context.(urlStatusRequest); ok {
			if t, ok := i.(time.Time); ok {
				if govalidator.IsUnixTime(fmt.Sprintf("%d", t.Unix())) {
					return true
				}
			}
		}
		return false
	}))
	if _, err := govalidator.ValidateStruct(r); err != nil {
		if !govalidator.IsUnixTime(strconv.Itoa(int(r.FromTime.Unix()))) || govalidator.IsUnixTime(strconv.Itoa(int(r.ToTime.Unix()))) {
			e := common.NewValidationError(err, "error validating url status request")
			return e
		}
	}
	url.Address = r.Address
	return nil
}
