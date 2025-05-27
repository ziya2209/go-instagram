package middleware

import (
	"context"
	"instagram/jwt"
	"net/http"
)

func ContextUpdater(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = WithNewRequestId(r)
		next.ServeHTTP(w, r)
	})
}

type ContextKey string

const (
	// RequestIDKey is the key for the request ID in the context
	RequestIDKey ContextKey = "requestID"
)
func withUserName(r *http.Request) *http.Request {
	ctx := r.Context()
	username, _ := jwt.GetUsernameFromRequest(r)

	ctx = context.WithValue(ctx, "username", username)
	return r.WithContext(ctx)
}

func WithNewRequestId(r *http.Request) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, RequestIDKey, "12345")
	return r.WithContext(ctx)
}
