package user

import (
	"context"
)

type UserWriteRepository interface {
	UserReadRepository
	CreateUser(ctx context.Context, email, password, name string) (*User, error)
}
type UserReadRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id int64) (*User, error)
}
