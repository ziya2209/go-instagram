package dto

type CreateUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,max=72"`
	Email    string `json:"email" validate:"required,email"`
	Age      int    `json:"age" validate:"required,min=12"`
	Bio      string `json:"bio" validate:"required"`
}
type UserDetailsResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	Bio      string `json:"bio"`
}

type LikePostRequest struct {
	PostId int `json:"post_id" validate:"required"`
}
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
type LoginResponse struct {
	Token string `json:"token"`
}
type CreatePostRequest struct {
	Caption string `json:"caption" validate:"required"`
	Url     string `json:"url" validate:"required"`
}
type FollowRequest struct {
	Username string `json:"username" validate:"required"`
}
type FollowersResponse struct {
	Followers []string `json:"followers"`
}
type LikePostResponse struct {
	PostId int `json:"post_id"`
}
type CommentRequest struct {
	ParentId   int    `json:"parent_id" validate:"required"`
	Comment    string `json:"comment" validate:"required"`
	ParentType string `json:"parent_type" validate:"required"` // e.g., "post" or "comment"
}

type HomepageResponse struct {
	Posts      []PostWithAuthor `json:"posts"`
	TotalPosts int              `json:"total_posts"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	TotalPages int              `json:"total_pages"`
	HasNext    bool             `json:"has_next"`
	HasPrev    bool             `json:"has_prev"`
}

type PostWithAuthor struct {
	Id             int    `json:"id"`
	UserId         int    `json:"user_id"`
	Url            string `json:"url"`
	Caption        string `json:"caption"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	AuthorUsername string `json:"author_username"`
	LikeCount      int    `json:"like_count"`
	CommentCount   int    `json:"comment_count"`
}
