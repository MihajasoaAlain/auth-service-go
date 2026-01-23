package httpapi

import (
	"encoding/json"
	"net/http"
	"time"

	"auth-service/internal/auth"
	"auth-service/internal/store"

	"github.com/google/uuid"
)

type AuthAPI struct {
	Users     store.Store
	JWT       auth.JWT
	AccessTTL time.Duration
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

	token, _ := a.JWT.SignAccess(user.ID, a.AccessTTL)

	writeJSON(w, 200, map[string]string{
		"access_token": token,
	})
}

func (a AuthAPI) Me(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	writeJSON(w, 200, map[string]string{"user_id": userID})
}
