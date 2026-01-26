package httpapi

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"auth-service/internal/auth"
)

func randStateGH() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b), err
}

func (a AuthAPI) GitHubStart(w http.ResponseWriter, r *http.Request) {
	state, err := randStateGH()
	if err != nil {
		http.Error(w, "server error", 500)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state_gh",
		Value:    state,
		HttpOnly: true,
		MaxAge:   300,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	u := a.Github.OAuthConfig.AuthCodeURL(state)
	http.Redirect(w, r, u, http.StatusFound)
}

func (a AuthAPI) GitHubCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	if code == "" || state == "" {
		http.Error(w, "missing code/state", 400)
		return
	}

	c, err := r.Cookie("oauth_state_gh")
	if err != nil || c.Value != state {
		http.Error(w, "invalid state", 400)
		return
	}

	tok, err := a.Github.OAuthConfig.Exchange(ctx, code)
	if err != nil {
		http.Error(w, "exchange failed", 401)
		return
	}

	req, _ := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user", nil)
	req.Header.Set("Authorization", "Bearer "+tok.AccessToken)
	req.Header.Set("Accept", "application/vnd.github+json")

	res, err := http.DefaultClient.Do(req)
	if err != nil || res.StatusCode >= 300 {
		http.Error(w, "github user fetch failed", 401)
		return
	}
	defer res.Body.Close()

	var ghUser struct {
		ID        int64  `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
	}
	_ = json.NewDecoder(res.Body).Decode(&ghUser)

	email := ""
	req2, _ := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user/emails", nil)
	req2.Header.Set("Authorization", "Bearer "+tok.AccessToken)
	req2.Header.Set("Accept", "application/vnd.github+json")

	res2, err := http.DefaultClient.Do(req2)
	if err == nil && res2.StatusCode < 300 {
		defer res2.Body.Close()
		var emails []struct {
			Email    string `json:"email"`
			Primary  bool   `json:"primary"`
			Verified bool   `json:"verified"`
		}
		_ = json.NewDecoder(res2.Body).Decode(&emails)
		for _, e := range emails {
			if e.Primary && e.Verified {
				email = e.Email
				break
			}
		}
	}

	providerUserID := url.QueryEscape((func() string {
		return jsonNumberToString(ghUser.ID)
	})())

	displayName := ghUser.Name
	if displayName == "" {
		displayName = ghUser.Login
	}

	userID, err := a.UpsertOAuthUser(ctx, "github", providerUserID, email, displayName, ghUser.AvatarURL)
	if err != nil {
		http.Error(w, "server error", 500)
		return
	}

	access, _ := a.JWT.SignAccess(userID, a.AccessTTL)
	refresh, _ := auth.NewRefreshToken()
	_ = a.Tokens.SaveRefreshToken(ctx, userID, refresh, time.Now().Add(a.RefreshTTL))

	redirect := a.Github.WebOrigin + "/oauth/success?access=" +
		url.QueryEscape(access) + "&refresh=" + url.QueryEscape(refresh)
	http.Redirect(w, r, redirect, http.StatusFound)
}

func jsonNumberToString(n int64) string {
	return fmt.Sprintf("%d", n)
}
