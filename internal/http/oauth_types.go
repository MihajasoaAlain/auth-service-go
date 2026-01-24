package httpapi

import "context"

type OAuthStore interface {
	GetUserIDByIdentity(ctx context.Context, provider, providerUserID string) (string, error)
	CreateIdentity(ctx context.Context, userID, provider, providerUserID, email string) error
	GetUserIDByEmail(ctx context.Context, email string) (string, error)
	CreateUserBasic(ctx context.Context, userID, email, name, avatar string) error
}
