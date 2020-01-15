package common

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

var JWTSecret = []byte("SuperSecret")

func GenerateJWT(id int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	t, err := token.SignedString(JWTSecret)
	if err != nil {
		return "", NewRequestError("Error generating JWT token", err, http.StatusInternalServerError)
	}
	return t, nil
}
