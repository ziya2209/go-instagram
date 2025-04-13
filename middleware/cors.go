package middleware

import (
	"net/http"

	"github.com/gorilla/handlers"
)

// q: what is CORS
// a: CORS (Cross-Origin Resource Sharing) is a security feature implemented by web browsers that allows or restricts web applications running at one origin to make requests to resources from a different origin.
// It is used to prevent malicious websites from making unauthorized requests to other domains.

func CORS(next http.Handler) http.Handler {
	origins := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"})
	headers := handlers.AllowedHeaders([]string{"Origin", "Content-Type", "Accept", "Authorization"})
	maxAge := handlers.MaxAge(3600)

	middleware := handlers.CORS(origins, methods, headers, maxAge)
	return middleware(next)
}
