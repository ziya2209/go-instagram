package handler

import (
	"encoding/json"
	"instagram/dto"
	"instagram/jwt"
	"instagram/models"
	"net/http"
	"time"
)

func (h *instaHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var createPostReq dto.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&createPostReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	// Validate required fields
	if createPostReq.Caption == "" || createPostReq.Url == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Caption and URL are required",
		})
		return
	}

	// Get username from JWT token
	username, err := jwt.GetUsernameFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid token",
		})
		return
	}

	// Get user ID from username
	user, err := h.userRepo.GetByUsername(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to get user information",
		})
		return
	}

	// Create new post
	post := &models.Post{
		UserId:    user.Id,
		Url:       createPostReq.Url,
		Caption:   createPostReq.Caption,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	// Save post to database
	if err := h.postRepo.Create(post); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to create post",
		})
		return
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Post created successfully",
		"postId":  post.Id,
	})
}
