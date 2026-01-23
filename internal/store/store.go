package store

import "context"

type User struct {
	ID           string
	Email        string
	PasswordHash string
}

type Store interface {
	CreateUser(ctx context.Context, u User) error
	GetUserByEmail(ctx context.Context, email string) (User, error)
}
