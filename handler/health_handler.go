package handler

import (
	"encoding/json"
	"net/http"
)

func (i *instaHandler) Health(w http.ResponseWriter, r *http.Request) {
	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write a simple JSON response indicating that the service is healthy
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "all ok",
	})
}
