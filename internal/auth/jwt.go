package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	Secret []byte
	Issuer string
}

func (j JWT) SignAccess(userID string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"iss": j.Issuer,
		"exp": time.Now().Add(ttl).Unix(),
		"iat": time.Now().Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(j.Secret)
}

func (j JWT) Verify(tokenStr string) (string, error) {
	t, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) { return j.Secret, nil })
	if err != nil || !t.Valid {
		return "", err
	}
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return "", jwt.ErrTokenInvalidClaims
	}
	sub, _ := claims["sub"].(string)
	if sub == "" {
		return "", jwt.ErrTokenInvalidClaims
	}
	return sub, nil
}
