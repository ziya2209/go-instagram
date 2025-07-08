package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
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
	// UsernameKey is the key for the username in the context
	UsernameKey ContextKey = "username"
)



func WithNewRequestId(r *http.Request) *http.Request {
	ctx := r.Context()
	requestID := uuid.New().String()
	ctx = context.WithValue(ctx, RequestIDKey, requestID)
	return r.WithContext(ctx)
}
