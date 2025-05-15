package middleware

import (
	"log"
	"net/http"
	"time"
)
// Logger is a middleware handler that logs HTTP requests
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a custom response writer to capture status code
		lrw := NewLoggingResponseWriter(w)

		// Call the next handler
		next.ServeHTTP(lrw, r)

		// Calculate request duration
		duration := time.Since(start)

		// Log the request details
		log.Printf(
			"[%s] %s %s | Status: %d | Duration: %v",
			r.RemoteAddr,
			r.Method,
			r.URL.Path,
			lrw.statusCode,
			duration,
		)
	})
}

// LoggingResponseWriter stores the status code of the response
type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// NewLoggingResponseWriter creates a new LoggingResponseWriter
func NewLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{w, http.StatusOK}
}

// WriteHeader captures the status code
func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
