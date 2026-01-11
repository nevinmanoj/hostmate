package user

import (
	user "github.com/nevinmanoj/hostmate/internal/domain/user"
)

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required"`
}

type LoginUserResponse struct {
	UserResponse
	Token string `json:"token"`
}

type UserResponse struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func ToUserResponse(u *user.User) UserResponse {
	return UserResponse{
		Email: u.Email,
		Name:  u.Name,
	}
}
func ToLoginUserResponse(u *user.User, token string) LoginUserResponse {
	return LoginUserResponse{
		UserResponse: UserResponse{
			Email: u.Email,
			Name:  u.Name,
		},
		Token: token,
	}
}
