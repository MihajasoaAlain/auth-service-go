package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OAuthStore struct {
	DB *pgxpool.Pool
}

func (s OAuthStore) GetUserIDByIdentity(ctx context.Context, provider, providerUserID string) (string, error) {
	var userID string
	err := s.DB.QueryRow(ctx,
		`SELECT user_id FROM user_identities
		 WHERE provider=$1 AND provider_user_id=$2`,
		provider, providerUserID,
	).Scan(&userID)

	if err != nil {
		return "", errors.New("not found")
	}
	return userID, nil
}

func (s OAuthStore) CreateIdentity(ctx context.Context, userID, provider, providerUserID, email string) error {
	_, err := s.DB.Exec(ctx,
		`INSERT INTO user_identities (user_id, provider, provider_user_id, email)
		 VALUES ($1,$2,$3,$4)
		 ON CONFLICT (provider, provider_user_id) DO NOTHING`,
		userID, provider, providerUserID, email,
	)
	return err
}

func (s OAuthStore) GetUserIDByEmail(ctx context.Context, email string) (string, error) {
	var userID string
	err := s.DB.QueryRow(ctx,
		`SELECT id FROM users WHERE email=$1`,
		email,
	).Scan(&userID)

	if err != nil {
		return "", errors.New("not found")
	}
	return userID, nil
}

func (s OAuthStore) CreateUserBasic(ctx context.Context, userID, email, name, avatar string) error {
	_, err := s.DB.Exec(ctx,
		`INSERT INTO users (id,email,name,avatar_url)
		 VALUES ($1,$2,$3,$4)`,
		userID, email, name, avatar,
	)
	return err
}
