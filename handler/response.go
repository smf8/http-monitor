package handler

import (
	"github.com/smf8/http-monitor/common"
	"github.com/smf8/http-monitor/model"
	"time"
)

type responseData struct {
	Data interface{} `json:"data"`
}

func NewResponseData(data interface{}) *responseData {
	return &responseData{data}
}

type userResponse struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

func NewUserResponse(user *model.User) *userResponse {
	token, _ := common.GenerateJWT(user.ID)
	ur := &userResponse{Username: user.Username, Token: token}
	return ur
}

// TODO : as model.url struct does not have an inner User instance, create one for it and update urlResponse to send username instead of it's id
type urlResponse struct {
	ID          int       `json:"id"`
	URL         string    `json:"url"`
	UserID      uint      `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	Threshold   int       `json:"threshold"`
	FailedTimes int       `json:"failed_times"`
}

func newURLResponse(url *model.URL) *urlResponse {
	u := new(urlResponse)
	u.URL = url.Address
	u.UserID = url.UserId
	u.CreatedAt = url.CreatedAt
	u.Threshold = url.Threshold
	u.FailedTimes = url.FailedTimes
	return u
}

type urlListResponse struct {
	URLs     []*urlResponse `json:"urls"`
	UrlCount int            `json:"url_count"`
}

func newURLListResponse(list []model.URL) *urlListResponse {
	resp := new(urlListResponse)
	resp.URLs = make([]*urlResponse, 0)
	for i := range list {
		resp.URLs = append(resp.URLs, newURLResponse(&list[i]))
	}
	resp.UrlCount = len(list)
	return resp
}

type requestResponse struct {
	URL        string    `json:"url"`
	ResultCode string    `json:"result_code"`
	CreatedAt  time.Time `json:"created_at"`
}

type requestListResponse struct {
	Requests []*requestResponse `json:"requests"`
}

type alertResponse struct {
	URL         string `json:"url"`
	FailedTimes int    `json:"failed_times"`
}
type alertListResponse struct {
	Alarms []*alertResponse `json:"alarms"`
}
