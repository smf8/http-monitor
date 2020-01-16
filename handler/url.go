package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/smf8/http-monitor/common"
	"github.com/smf8/http-monitor/model"
	"net/http"
)

func (h *Handler) FetchURLs(c echo.Context) error {
	userID := extractID(c)
	urls, err := h.st.GetURLsByUser(userID)
	if err != nil {
		return common.NewRequestError("Error retrieving urls from database, maybe check your token again", err, http.StatusBadRequest)
	}
	resp := newURLListResponse(urls)
	return c.JSON(http.StatusOK, NewResponseData(resp))
}

func (h *Handler) CreateURL(c echo.Context) error {
	userID := extractID(c)
	req := &urlCreateRequest{}
	url := &model.URL{}

	if err := req.bind(c, url); err != nil {
		return err
	}
	url.UserId = userID
	// adding url to database
	if err := h.st.AddURL(url); err != nil {
		// internal error
		return common.NewRequestError("error adding url to database", err, http.StatusInternalServerError)
	}
	// adding url to monitor scheduler
	h.sch.Mnt.AddURL([]model.URL{*url})
	return c.JSON(http.StatusCreated, NewResponseData("URL created successfully"))
}
