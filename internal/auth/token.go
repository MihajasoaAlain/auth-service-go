package auth

import (
	"crypto/rand"
	"encoding/base64"
)

func NewOneTimeToken() (string, error) {
	b := make([]byte, 48) // ~64 chars
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
