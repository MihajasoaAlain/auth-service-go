package postgres

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TokenStore struct {
	DB *pgxpool.Pool
}

func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func (s TokenStore) SaveRefreshToken(ctx context.Context, userID string, rawToken string, exp time.Time) error {
	_, err := s.DB.Exec(ctx,
		`INSERT INTO refresh_tokens (token_hash,user_id,expires_at) VALUES ($1,$2,$3)`,
		hashToken(rawToken), userID, exp,
	)
	return err
}

func (s TokenStore) ValidateRefreshToken(ctx context.Context, rawToken string) (userID string, ok bool) {
	err := s.DB.QueryRow(ctx,
		`SELECT user_id
		   FROM refresh_tokens
		  WHERE token_hash=$1
		    AND revoked_at IS NULL
		    AND expires_at > now()`,
		hashToken(rawToken),
	).Scan(&userID)

	return userID, err == nil
}

func (s TokenStore) RevokeRefreshToken(ctx context.Context, rawToken string) error {
	_, err := s.DB.Exec(ctx,
		`UPDATE refresh_tokens SET revoked_at=now() WHERE token_hash=$1`,
		hashToken(rawToken),
	)
	return err
}
