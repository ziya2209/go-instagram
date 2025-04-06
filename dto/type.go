package dto

type CreateUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,max=72"`
	Email    string `json:"email" validate:"required,email"`
	Age      int    `json:"age" validate:"required,min=12"`
	Bio      string `json:"bio" validate:"required"`
}
