package httpapi

import (
	"context"

	"github.com/google/uuid"
)

func (a AuthAPI) UpsertOAuthUser(
	ctx context.Context,
	provider, providerUserID, email, name, avatar string,
) (string, error) {

	if uid, err := a.OAuth.GetUserIDByIdentity(ctx, provider, providerUserID); err == nil {
		return uid, nil
	}

	if email != "" {
		if uid, err := a.OAuth.GetUserIDByEmail(ctx, email); err == nil {
			_ = a.OAuth.CreateIdentity(ctx, uid, provider, providerUserID, email)
			return uid, nil
		}
	}

	userID := uuid.NewString()
	if err := a.OAuth.CreateUserBasic(ctx, userID, email, name, avatar); err != nil {
		return "", err
	}

	if err := a.OAuth.CreateIdentity(ctx, userID, provider, providerUserID, email); err != nil {
		return "", err
	}

	return userID, nil
}
