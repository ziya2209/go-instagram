package handler

import (
	"encoding/json"
	"fmt"
	"instagram/dto"
	"instagram/repo"
	"net/http"
)

func (h *instaHandler) Follow(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the context
	var fr dto.FollowRequest
	err := json.NewDecoder(r.Body).Decode(&fr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body",
		})
		return
	}
	// Assuming the user ID is stored in the context
	username, ok := r.Context().Value("username").(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Unauthorized",
		})
		return
	}
	// Validate the request
	if fr.Username == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Followed username is required",
		})
		return
	}
	// Get the user ID of the current user
	user, err := h.userRepo.GetByUsername(username)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Unauthorized",
		})
		return
	}

	// Get the user ID of the user to be followed
	followedUser, err := h.userRepo.GetByUsername(fr.Username)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "User not found",
		})
		return
	}
	// Check if the user is trying to follow themselves
	if followedUser.Id == user.Id {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "You cannot follow yourself",
		})
		return
	}
	// Call the follow repository to follow the user
	err = h.followRepo.FollowUser(user.Id, followedUser.Id)
	if err != nil {
		if err == repo.ErrDuplicateFollow {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "You are already following this user",
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to follow user",
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Successfully followed user",
	})
}

func (h *instaHandler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	// Get username from context
	username, ok := r.Context().Value("username").(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Unauthorized",
		})
		return
	}

	// Get current user
	user, err := h.userRepo.GetByUsername(username)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Unauthorized",
		})
		return
	}

	// Get followers for the current user
	followers, err := h.followRepo.GetFollowers(user.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to fetch followers",
		})
		return
	}

	// Get usernames for all followed IDs
	followedUsernames := make([]string, 0)
	for _, follower := range followers {
		user, err := h.userRepo.GetById(follower.FollowedId)
		if err != nil {
			continue // Skip if user not found
		}
		followedUsernames = append(followedUsernames, user.Username)
		
		
	}
	fmt.Println(followedUsernames)

	// Return response
	w.WriteHeader(http.StatusOK)
	response := dto.FollowersResponse{
		Followers: followedUsernames,
	}
	json.NewEncoder(w).Encode(response)
}
