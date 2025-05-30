package middleware

import (
	"net/http"
	"strings"

	"instagram/jwt"
)

var secretKey = "secret"

// JWTAuthMiddleware validates the JWT token in the Authorization header.
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate the token
		token, err := jwt.VerifyToken(tokenString)
		if err != nil || !token {		
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}
		r = withUserName(r)

		next.ServeHTTP(w, r)
	})
}
