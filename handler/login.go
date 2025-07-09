package handler

import (
	"encoding/json"
	"errors"
	"instagram/dto"
	"instagram/jwt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (h *instaHandler) Login(w http.ResponseWriter, r *http.Request) {

	// Parse login request data
	var loginReq dto.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	// Validate required fields
	if loginReq.Username == "" || loginReq.Password == "" {
		writeJSONError(w, http.StatusBadRequest, "Username and password are required")
		return
	}

	user, err := h.userRepo.GetByUsername(loginReq.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeJSONError(w, http.StatusUnauthorized, "Invalid username or password")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "Failed to fetch user")
		return
	}
	bycryptErr := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginReq.Password))
	if bycryptErr != nil {
		writeJSONError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	token, err := jwt.NewToken(user.Username)

	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	response := dto.LoginResponse{
		Token: token,
	}
	writeJSONResponse(w, http.StatusOK, response)
}
