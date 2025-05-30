package jwt

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

var secretKey = "secret"

func NewToken(username string) (string, error) {
	// Create a new token object
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username

	// Generate encoded token and send it as response
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
func GetUsernameFromHeader(r *http.Request) (string, error) {
	tokenString := r.Header.Get("Authorization")[7:] // Remove "Bearer " prefix

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, http.ErrAbortHandler
		}
		return []byte(secretKey), nil
	})
	if err != nil || !token.Valid {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if username, ok := claims["username"].(string); ok {
			return username, nil
		}
	}
	return "", http.ErrAbortHandler
}

func VerifyToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, http.ErrAbortHandler
		}
		return []byte(secretKey), nil
	})
	if err != nil {

		return false, err
	}
	return token.Valid, nil
}
