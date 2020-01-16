package handler

import (
	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
	"github.com/smf8/http-monitor/common"
	"github.com/smf8/http-monitor/model"
	"net/http"
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
	URLID uint `json:"url_id"`
}

func (r *alertDismissRequest) bind(c echo.Context, url *model.URL) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if _, err := govalidator.ValidateStruct(r); err != nil {
		e := common.NewValidationError(err, "Error validating alert dismiss request")
		return e
	}
	url.ID = r.URLID
	return nil
}

type urlStatusRequest struct {
	FromTime int64 `valid:"optional, time~Provide time as unix timestamp before current time" json:"from_time, omitempty" query:"from_time"`
	ToTime   int64 `valid:"optional, time~Provide time as unix timestamp before current time" json:"to_time, omitempty" query:"to_time"`
}

func (r *urlStatusRequest) parse(c echo.Context) error {
	if err := c.Bind(r); err != nil {
		return common.NewRequestError("error parsing url status request, if you want to specify time, use unix timestamp", err, http.StatusBadRequest)
	}
	govalidator.CustomTypeTagMap.Set("time", govalidator.CustomTypeValidator(func(i interface{}, context interface{}) bool {
		if _, ok := context.(urlStatusRequest); ok {
			if t, ok := i.(int64); ok {
				if time.Now().Unix() > t {
					return true
				}
			}
		}
		return false
	}))
	if _, err := govalidator.ValidateStruct(r); err != nil {
		e := common.NewValidationError(err, "error validating url status request")
		return e
	}
	if r.FromTime > r.ToTime && r.ToTime != 0 {
		return common.NewRequestError("end of time interval must be later than it's start", nil, http.StatusBadRequest)
	}
	return nil
}
