package postgres

import (
	"context"
	"errors"
	"time"

	"auth-service/internal/store"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserStore struct {
	DB *pgxpool.Pool
}

func (s UserStore) CreateUser(ctx context.Context, u store.User) error {
	_, err := s.DB.Exec(ctx,
		`INSERT INTO users (id,email,password_hash) VALUES ($1,$2,$3)`,
		u.ID, u.Email, u.PasswordHash,
	)
	return err
}

func (s UserStore) GetUserByEmail(ctx context.Context, email string) (store.User, error) {
	var u store.User
	err := s.DB.QueryRow(ctx,
		`SELECT id,email,password_hash,email_verified_at FROM users WHERE email=$1`,
		email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.EmailVerifiedAt)

	if err != nil {
		return u, errors.New("user not found")
	}
	return u, nil
}

func (s UserStore) SetEmailVerified(ctx context.Context, userID string, at time.Time) error {
	_, err := s.DB.Exec(ctx,
		`UPDATE users SET email_verified_at=$1 WHERE id=$2`,
		at, userID,
	)
	return err
}

func (s UserStore) SetPasswordHash(ctx context.Context, userID string, hash string) error {
	_, err := s.DB.Exec(ctx,
		`UPDATE users SET password_hash=$1 WHERE id=$2`,
		hash, userID,
	)
	return err
}
