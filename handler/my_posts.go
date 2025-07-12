package handler

import (
	"encoding/json"
	"instagram/dto"
	"instagram/middleware"
	"math"
	"net/http"
	"strconv"
)

// MyPosts returns all posts created by the current user with pagination
func (h *instaHandler) MyPosts(w http.ResponseWriter, r *http.Request) {
	// Get username from context (set by JWT middleware)
	username, ok := r.Context().Value(middleware.UsernameKey).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get the user details to get user ID
	user, err := h.userRepo.GetByUsername(username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Parse pagination parameters
	page := 1
	limit := 10 // Default limit

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Get posts for the current user
	posts, totalPosts, err := h.postRepo.GetPostsByUserId(user.Id, page, limit)
	if err != nil {
		http.Error(w, "Failed to get posts", http.StatusInternalServerError)
		return
	}

	// Convert posts to PostWithAuthor format
	var postsWithAuthor []dto.PostWithAuthor
	for _, post := range posts {
		// Get like count
		likeCount, err := h.postRepo.GetLikeCount(post.Id)
		if err != nil {
			likeCount = 0 // Continue with 0 if error
		}

		// Get comment count
		commentCount, err := h.postRepo.GetCommentCount(post.Id)
		if err != nil {
			commentCount = 0 // Continue with 0 if error
		}

		postWithAuthor := dto.PostWithAuthor{
			Id:             post.Id,
			UserId:         post.UserId,
			Url:            post.Url,
			Caption:        post.Caption,
			CreatedAt:      post.CreatedAt,
			UpdatedAt:      post.UpdatedAt,
			AuthorUsername: username, // Current user is the author
			LikeCount:      likeCount,
			CommentCount:   commentCount,
		}
		postsWithAuthor = append(postsWithAuthor, postWithAuthor)
	}

	// Calculate pagination info
	totalPages := int(math.Ceil(float64(totalPosts) / float64(limit)))
	hasNext := page < totalPages
	hasPrev := page > 1

	response := dto.HomepageResponse{
		Posts:      postsWithAuthor,
		TotalPosts: totalPosts,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
