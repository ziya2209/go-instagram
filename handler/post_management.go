package handler

import (
	"encoding/json"
	"instagram/middleware"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// DeletePost handles the deletion of a post
func (h *instaHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	// Get the username from JWT context
	username, ok := r.Context().Value(middleware.UsernameKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Unauthorized",
		})
		return
	}

	// Get current user
	currentUser, err := h.userRepo.GetByUsername(username)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "User not found",
		})
		return
	}

	// Get post ID from URL parameters
	vars := mux.Vars(r)
	postIdStr, exists := vars["postId"]
	if !exists {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Post ID is required",
		})
		return
	}

	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid post ID format",
		})
		return
	}

	// Check if the post exists and if the current user is the owner
	isOwner, err := h.postRepo.IsPostOwner(postId, currentUser.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "Failed to verify post ownership",
			"details": err.Error(),
		})
		return
	}

	if !isOwner {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "You can only delete your own posts",
		})
		return
	}

	// Delete the post
	err = h.postRepo.Delete(postId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "Failed to delete post",
			"details": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Post deleted successfully",
		"post_id": postId,
	})
}
