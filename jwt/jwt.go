package jwt

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
)




var secretKey = "secret"
func NewToken() {
}


func VerifyToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, http.ErrAbortHandler
		}
		return []byte(secretKey), nil
	})
	return token.Valid, err
}