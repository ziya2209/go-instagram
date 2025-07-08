package handler

import (
	"encoding/json"
	"instagram/dto"
	"instagram/models"
	"net/http"
	"strconv"
)

func (h *instaHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the context
	username, ok := r.Context().Value("username").(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Unauthorized",
		})
		return
	}

	var commentReq dto.CommentRequest
	err := json.NewDecoder(r.Body).Decode(&commentReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	if commentReq.ParentId == 0 || commentReq.Comment == "" || commentReq.ParentType == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Parent ID, comment, and parent type are required",
		})
		return
	}

	user, err := h.userRepo.GetByUsername(username)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Unauthorized",
		})
		return
	}

	comment := models.Comment{
		UserId:     user.Id,
		Message:    commentReq.Comment,
		ParentId:   commentReq.ParentId,
		ParentType: commentReq.ParentType,
	}

	err = h.commentRepo.Create(&comment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "Failed to add comment",
			"details": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

func (h *instaHandler) PostGetComments(w http.ResponseWriter, r *http.Request) {
	// Get the post ID from the request
	postId := r.URL.Query().Get("postId")
	if postId == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Post ID is required",
		})
		return
	}
	// Convert postId to int
	postIdInt, err := strconv.Atoi(postId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid Post ID format",
		})
		return
	}

	comments, err := h.commentRepo.GetByPostId(postIdInt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "Failed to retrieve comments",
			"details": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comments)
}
