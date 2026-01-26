package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"auth-service/internal/auth"
	"auth-service/internal/store"

	"github.com/google/uuid"
)

type AuthAPI struct {
	Users  store.Store
	Tokens interface {
		SaveRefreshToken(ctx context.Context, userID string, rawToken string, exp time.Time) error
		ValidateRefreshToken(ctx context.Context, rawToken string) (string, bool)
		RevokeRefreshToken(ctx context.Context, rawToken string) error
	}
	JWT        auth.JWT
	AccessTTL  time.Duration
	RefreshTTL time.Duration
	OAuth      OAuthStore
	Google     GoogleAuth
	Github     GitHubAuth
}

type creds struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

func (a AuthAPI) Register(w http.ResponseWriter, r *http.Request) {
	var c creds
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	hash, err := auth.HashPassword(c.Password)
	if err != nil {
		http.Error(w, "server error", 500)
		return
	}

	user := store.User{
		ID:           uuid.NewString(),
		Email:        c.Email,
		PasswordHash: hash,
	}

	if err := a.Users.CreateUser(r.Context(), user); err != nil {
		http.Error(w, "email already exists", 409)
		return
	}

	writeJSON(w, 201, map[string]string{
		"id":    user.ID,
		"email": user.Email,
	})
}

func (a AuthAPI) Login(w http.ResponseWriter, r *http.Request) {
	var c creds
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "invalid body", 400)
		return
	}

	user, err := a.Users.GetUserByEmail(r.Context(), c.Email)
	if err != nil || !auth.CheckPassword(user.PasswordHash, c.Password) {
		http.Error(w, "invalid credentials", 401)
		return
	}

	access, err := a.JWT.SignAccess(user.ID, a.AccessTTL)
	if err != nil {
		http.Error(w, "server error", 500)
		return
	}
	refresh, err := auth.NewRefreshToken()
	if err != nil {
		http.Error(w, "server error", 500)
		return
	}
	exp := time.Now().Add(a.RefreshTTL)
	if err := a.Tokens.SaveRefreshToken(r.Context(), user.ID, refresh, exp); err != nil {
		http.Error(w, "server error", 500)
		return
	}

	writeJSON(w, 200, map[string]any{
		"access_token":  access,
		"expires_in":    int(a.AccessTTL.Seconds()),
		"refresh_token": refresh,
	})
}

func (a AuthAPI) Me(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	writeJSON(w, 200, map[string]string{"user_id": userID})
}

func (a AuthAPI) Refresh(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.RefreshToken == "" {
		http.Error(w, "bad request", 400)
		return
	}

	userID, ok := a.Tokens.ValidateRefreshToken(r.Context(), body.RefreshToken)
	if !ok {
		http.Error(w, "invalid refresh token", 401)
		return
	}

	_ = a.Tokens.RevokeRefreshToken(r.Context(), body.RefreshToken)

	newRefresh, err := auth.NewRefreshToken()
	if err != nil {
		http.Error(w, "server error", 500)
		return
	}
	_ = a.Tokens.SaveRefreshToken(r.Context(), userID, newRefresh, time.Now().Add(a.RefreshTTL))

	access, err := a.JWT.SignAccess(userID, a.AccessTTL)
	if err != nil {
		http.Error(w, "server error", 500)
		return
	}

	writeJSON(w, 200, map[string]any{
		"access_token":  access,
		"expires_in":    int(a.AccessTTL.Seconds()),
		"refresh_token": newRefresh,
	})
}

func (a AuthAPI) Logout(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.RefreshToken == "" {
		http.Error(w, "bad request", 400)
		return
	}

	_ = a.Tokens.RevokeRefreshToken(r.Context(), body.RefreshToken)
	writeJSON(w, 200, map[string]string{"status": "ok"})
}
