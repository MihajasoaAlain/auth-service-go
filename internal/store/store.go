package store

import (
	"context"
	"time"
)

type User struct {
	ID              string
	Email           string
	PasswordHash    string
	EmailVerifiedAt *time.Time
}

type Store interface {
	CreateUser(ctx context.Context, u User) error
	GetUserByEmail(ctx context.Context, email string) (User, error)
	SetEmailVerified(ctx context.Context, userID string, at time.Time) error
	SetPasswordHash(ctx context.Context, userID string, hash string) error
}
