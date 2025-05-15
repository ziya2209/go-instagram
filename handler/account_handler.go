package handler

import (
	"encoding/json"
	"instagram/dto"
	"instagram/models"
	"instagram/repo"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

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
	if createUserReq.Username == "" || createUserReq.Email == "" || createUserReq.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Name, email and password are required",
		})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(createUserReq.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to process password",
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

func (h *instaHandler)GetAllUser(w http.ResponseWriter, r *http.Request) {	
	// Fetch all users from the database
	users, err := h.userRepo.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to fetch users",
		})
		return
	}
	

	// Return the list of users
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
  udr := make([]dto.UserDetailsResponse, len(users))
	for i, user := range users {
		udr[i] = dto.UserDetailsResponse{
			Username: user.Username,
			Email:    user.Email,
			Age:      user.Age,
			Bio:      user.Bio,
		}
	}
	json.NewEncoder(w).Encode(udr)

}

// Placeholder implementations for other required methods
func (h *instaHandler) PostGetComments(w http.ResponseWriter, r *http.Request) {}
func (h *instaHandler) ShowHomePage(w http.ResponseWriter, r *http.Request)    {}
func (h *instaHandler) CreatePost(w http.ResponseWriter, r *http.Request)      {}
func (h *instaHandler) AddComment(w http.ResponseWriter, r *http.Request)      {}
func (h *instaHandler) LikePost(w http.ResponseWriter, r *http.Request)        {}
func (h *instaHandler) Login(w http.ResponseWriter, r *http.Request)           {}
