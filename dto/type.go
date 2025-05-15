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
