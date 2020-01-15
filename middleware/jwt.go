// codes from https://github.com/labstack/echo/blob/master/middleware/jwt.go
package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type whiteList struct {
	path   string
	method string
}

// authWhiteList specifies paths to be skipped by jwt authentication middleware
var authWhiteList []whiteList

// AddToWhiteList is used to add a path to skipper white list
// provide path relative to api version like /api/your/path/here as skipper uses strings.Contains to find whether
// it is in context path or not
func AddToWhiteList(path string, method string) {
	if authWhiteList == nil {
		authWhiteList = make([]whiteList, 0)
	}
	authWhiteList = append(authWhiteList, whiteList{path, method})
}

func skipper(c echo.Context) bool {
	for _, v := range authWhiteList {
		if c.Path() == v.path && c.Request().Method == v.method {
			return true
		}
	}
	return false
}

func JWT(key interface{}) echo.MiddlewareFunc {
	c := middleware.DefaultJWTConfig
	c.SigningKey = key
	c.Skipper = skipper
	return middleware.JWTWithConfig(c)
}
