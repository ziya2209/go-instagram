package handler

import (
	"encoding/json"
	"instagram/dto"
	"net/http"
)

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