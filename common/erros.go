package common

// got error handling idea from https://medium.com/ki-labs-engineering/rest-api-error-handling-in-go-behavioral-type-assertion-509d93636afd

import (
	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
	"net/http"
)

const (
	BindingError = "error parsing request, check request format and try again."
)

// ClientError has An additional ErrorBody() function which returns map[string]interface{}
// implementing Error() function is for server-side error logging
// implementing ErrorBody() is for json response to clients
type ClientError interface {
	Error() string
	ErrorBody() Error
}

type Error struct {
	Body map[string]interface{} `json:"errors"`
}

func NewError() *Error {
	body := make(map[string]interface{})
	return &Error{Body: body}
}

// RequestError implements ClientError interface
// call ErrorBody() on it to provide a map version of errors
type RequestError struct {
	Cause  error
	Detail string
	Status int
}

// NewBindingError creates an error with given data. "bindingType" is the type that is being bound
// provide short description for "detail" to display to user as error message
func NewRequestError(detail string, err error, statusCode int) *RequestError {
	if e, ok := err.(*RequestError); ok {
		return e
	}
	return &RequestError{err, detail, statusCode}
}
func (se *RequestError) Error() string {
	if se.Cause == nil {
		return se.Detail
	}
	return se.Detail + " : " + se.Cause.Error()
}
func (se *RequestError) ErrorBody() Error {
	body := NewError()
	body.Body["body"] = se.Detail
	return *body
}

type ValidationError struct {
	Cause    error
	ErrorMap map[string]string
	Detail   string
}

func NewValidationError(err error, detail string) *ValidationError {
	m := govalidator.ErrorsByField(err)
	return &ValidationError{err, m, detail}
}

func (e *ValidationError) Error() string {
	if e.Cause == nil {
		return e.Detail
	}
	return e.Detail + " : " + e.Cause.Error()
}

func (e *ValidationError) ErrorBody() Error {
	err := new(Error)
	m := make(map[string]interface{})
	for key, value := range e.ErrorMap {
		m[key] = value
	}
	err.Body = m
	return *err
}

// CustomHTTPErrorHandler, here we define what to do with each types of errors
func CustomHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	switch err.(type) {
	case ClientError:
		ce := err.(ClientError)
		e := c.JSON(code, ce.ErrorBody())
		if e != nil {
			c.Logger().Error(e)
		}

		c.Logger().Error(ce)
		switch e := ce.(type) {
		case *RequestError:
			code = e.Status
		case *ValidationError:
			code = http.StatusBadRequest
		}
	case *echo.HTTPError:
		e := c.JSON(err.(*echo.HTTPError).Code, err.(*echo.HTTPError).Message)
		if e != nil {
			c.Logger().Error(e)
		}
		c.Logger().Error(err)
	}
}
