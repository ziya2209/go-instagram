package handler

import (
	"instagram/dto"
	"net/http"
)

func (h *instaHandler) GetAllUser(w http.ResponseWriter, r *http.Request) {
	// Fetch all users from the database
	users, err := h.userRepo.GetAll()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	// Return the list of users
	udr := make([]dto.UserDetailsResponse, len(users))
	for i, user := range users {
		udr[i] = dto.UserDetailsResponse{
			Username: user.Username,
			Email:    user.Email,
			Age:      user.Age,
			Bio:      user.Bio,
		}
	}
	writeJSONResponse(w, http.StatusOK, udr)

}
