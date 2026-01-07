package user

import (
	"context"
	"fmt"

	"github.com/nevinmanoj/hostmate/internal/auth"
)

type UserService interface {
	CreateUser(ctx context.Context, email, password, name string) (*User, error)
	LoginUser(ctx context.Context, email, password string) (string, *User, error)
	GetUserByID(ctx context.Context, id int64) (*User, error)
}

type userService struct {
	repo      UserWriteRepository
	jwtSecret []byte
}

func NewUserService(repo UserWriteRepository, jwtSecret []byte) UserService {
	return &userService{repo: repo, jwtSecret: jwtSecret}
}

func (s *userService) CreateUser(ctx context.Context, email, password, name string) (*User, error) {
	return s.repo.CreateUser(ctx, email, password, name)
}
func (s *userService) LoginUser(ctx context.Context, email, password string) (string, *User, error) {
	// Implementation for user login
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", nil, err
	}
	err = auth.CheckPassword(password, user.PasswordHash)
	if err != nil {
		return "", nil, fmt.Errorf("Invalid credentials")
	}

	// create and issue JWT token
	token, err := auth.GenerateToken(user.ID, user.Email, s.jwtSecret)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

func (s *userService) GetUserByID(ctx context.Context, id int64) (*User, error) {
	return s.repo.GetUserByID(ctx, id)
}
