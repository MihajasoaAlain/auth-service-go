package httpapi

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"
	"time"

	"auth-service/internal/auth"
)

func randState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b), err
}

func (a AuthAPI) GoogleStart(w http.ResponseWriter, r *http.Request) {
	state, _ := randState()

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		HttpOnly: true,
		MaxAge:   300,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	url := a.Google.OAuthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusFound)
}

func (a AuthAPI) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	c, err := r.Cookie("oauth_state")
	if err != nil || c.Value != state {
		http.Error(w, "invalid state", 400)
		return
	}

	tok, err := a.Google.OAuthConfig.Exchange(ctx, code)
	if err != nil {
		http.Error(w, "exchange failed", 401)
		return
	}

	rawID, ok := tok.Extra("id_token").(string)
	if !ok {
		http.Error(w, "missing id_token", 401)
		return
	}

	idToken, err := a.Google.Verifier.Verify(ctx, rawID)
	if err != nil {
		http.Error(w, "invalid id_token", 401)
		return
	}

	var claims struct {
		Sub     string `json:"sub"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	_ = idToken.Claims(&claims)

	userID, err := a.UpsertOAuthUser(ctx, "google", claims.Sub, claims.Email, claims.Name, claims.Picture)
	if err != nil {
		http.Error(w, "server error", 500)
		return
	}

	access, _ := a.JWT.SignAccess(userID, a.AccessTTL)
	refresh, _ := auth.NewRefreshToken()
	_ = a.Tokens.SaveRefreshToken(ctx, userID, refresh, time.Now().Add(a.RefreshTTL))

	redirect := a.Google.WebOrigin + "/oauth/success?access=" +
		url.QueryEscape(access) + "&refresh=" + url.QueryEscape(refresh)

	http.Redirect(w, r, redirect, http.StatusFound)
}
