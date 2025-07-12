package handler

import (
	"encoding/json"
	"instagram/dto"
	"instagram/middleware"
	"net/http"
	"strconv"
)

// ShowHomePage handles the homepage request with pagination
func (h *instaHandler) ShowHomePage(w http.ResponseWriter, r *http.Request) {
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

	// Parse pagination parameters
	page := 1
	limit := 10 // Default posts per page

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 50 { // Max 50 posts per page
			limit = l
		}
	}

	// Get users that current user follows
	follows, err := h.followRepo.GetFollowers(currentUser.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "Failed to get followed users",
			"details": err.Error(),
		})
		return
	}

	// Extract followed user IDs
	var followedUserIds []int
	for _, follow := range follows {
		followedUserIds = append(followedUserIds, follow.FollowedId)
	}

	// If user doesn't follow anyone, return empty response
	if len(followedUserIds) == 0 {
		response := dto.HomepageResponse{
			Posts:      []dto.PostWithAuthor{},
			TotalPosts: 0,
			Page:       page,
			Limit:      limit,
			TotalPages: 0,
			HasNext:    false,
			HasPrev:    false,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Get posts from followed users with pagination
	posts, totalPosts, err := h.postRepo.GetPostsByUserIds(followedUserIds, page, limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "Failed to get posts",
			"details": err.Error(),
		})
		return
	}

	// Get additional information for each post (author, likes, comments)
	var postsWithInfo []dto.PostWithAuthor
	for _, post := range posts {
		// Get author information
		author, err := h.userRepo.GetById(post.UserId)
		if err != nil {
			// Skip this post if we can't get author info
			continue
		}

		// Get like count
		likeCount, err := h.postRepo.GetLikeCount(post.Id)
		if err != nil {
			likeCount = 0 // Default to 0 if error
		}

		// Get comment count
		commentCount, err := h.postRepo.GetCommentCount(post.Id)
		if err != nil {
			commentCount = 0 // Default to 0 if error
		}

		postWithInfo := dto.PostWithAuthor{
			Id:             post.Id,
			UserId:         post.UserId,
			Url:            post.Url,
			Caption:        post.Caption,
			CreatedAt:      post.CreatedAt,
			UpdatedAt:      post.UpdatedAt,
			AuthorUsername: author.Username,
			LikeCount:      likeCount,
			CommentCount:   commentCount,
		}
		postsWithInfo = append(postsWithInfo, postWithInfo)
	}

	// Calculate pagination info
	totalPages := (totalPosts + limit - 1) / limit // Ceiling division
	hasNext := page < totalPages
	hasPrev := page > 1

	response := dto.HomepageResponse{
		Posts:      postsWithInfo,
		TotalPosts: totalPosts,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
