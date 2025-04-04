package handler

import (
	"encoding/json"
	"instagram/models"
	"instagram/repo"
	"net/http"
)

// Implementation of InstaHandler
type instaHandler struct {
	userRepo repo.UserRepo
}

// NewInstaHandler creates a new instance of InstaHandler
func NewInstaHandler(ur repo.UserRepo) InstaHandler {
	// In a real implementation, you would inject the repository here
	return &instaHandler{
		userRepo: ur,
	}
}

// CreateAcc handles the creation of a new user account
func (h *instaHandler) CreateAcc(w http.ResponseWriter, r *http.Request) {
	// Parse user data from request
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	// Validate required fields
	if user.Name == "" || user.Email == "" || user.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Name, email and password are required",
		})
		return
	}

	// Create user in database
	err = h.userRepo.Create(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to create account",
		})
		return
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Account created successfully",
		"userId":  user.Id,
	})
}

// Placeholder implementations for other required methods
func (h *instaHandler) PostGetComments(w http.ResponseWriter, r *http.Request) {}
func (h *instaHandler) ShowHomePage(w http.ResponseWriter, r *http.Request)    {}
func (h *instaHandler) CreatePost(w http.ResponseWriter, r *http.Request)      {}
func (h *instaHandler) AddComment(w http.ResponseWriter, r *http.Request)      {}
func (h *instaHandler) LikePost(w http.ResponseWriter, r *http.Request)        {}
func (h *instaHandler) Login(w http.ResponseWriter, r *http.Request)           {}
