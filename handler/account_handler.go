package handler

import (
	"encoding/json"
	"instagram/dto"
	"instagram/models"
	"instagram/repo"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type instaHandler struct {
	userRepo   repo.UserRepo
	postRepo   repo.PostRepo
	followRepo repo.FollowRepo // Assuming followRepo is part of userRepo
}

// NewInstaHandler creates a new instance of InstaHandler
func NewInstaHandler(ur repo.UserRepo, pr repo.PostRepo, fr repo.FollowRepo) InstaHandler {
	return &instaHandler{
		userRepo:   ur,
		postRepo:   pr,
		followRepo: fr, // Assuming followRepo is part of userRepo
	}
}

// CreateAcc handles the creation of a new user account
func (h *instaHandler) CreateAcc(w http.ResponseWriter, r *http.Request) {
	// Parse createUserReq data from request
	var createUserReq dto.CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&createUserReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	// Validate required fields
	var emptyFields []string
	createUserReq.Username = strings.TrimSpace(createUserReq.Username)
	createUserReq.Email = strings.TrimSpace(createUserReq.Email)
	createUserReq.Password = strings.TrimSpace(createUserReq.Password)

	if createUserReq.Username == "" {
		emptyFields = append(emptyFields, "username")
	}
	if createUserReq.Email == "" {
		emptyFields = append(emptyFields, "email")
	}
	if createUserReq.Password == "" {
		emptyFields = append(emptyFields, "password")
	}

	if len(emptyFields) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":  "Required fields are empty",
			"fields": emptyFields,
		})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(createUserReq.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to hash password",
		})
		return
	}

	// Create user in database
	user := &models.User{
		Username:     createUserReq.Username,
		Email:        createUserReq.Email,
		PasswordHash: string(hashedPassword),
		Age:          createUserReq.Age,
		Bio:          createUserReq.Bio,
	}

	err = h.userRepo.Create(user)
	if err != nil {
		if err.Error() == "email already exists" {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Email already exists",
			})
			return
		}
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
func (h *instaHandler) AddComment(w http.ResponseWriter, r *http.Request)      {}
func (h *instaHandler) LikePost(w http.ResponseWriter, r *http.Request)        {}
