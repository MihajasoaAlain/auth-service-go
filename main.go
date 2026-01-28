package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"

	"auth-service/internal/auth"
	"auth-service/internal/config"
	"auth-service/internal/db"
	httpapi "auth-service/internal/http"
	"auth-service/internal/store/postgres"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found (using system env)")
	}

	cfg := config.Load()
	dbPool := db.Connect(cfg.DatabaseURL)

	userStore := postgres.UserStore{DB: dbPool}
	tokenStore := postgres.TokenStore{DB: dbPool}
	oauthStore := postgres.OAuthStore{DB: dbPool}

	oidcProvider, err := oidc.NewProvider(context.Background(), "https://accounts.google.com")
	if err != nil {
		log.Fatal(err)
	}

	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	googleRedirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	webOrigin := os.Getenv("WEB_ORIGIN")

	googleAuth := httpapi.GoogleAuth{
		OAuthConfig: oauth2.Config{
			ClientID:     googleClientID,
			ClientSecret: googleClientSecret,
			RedirectURL:  googleRedirectURL,
			Endpoint:     google.Endpoint,
			Scopes:       []string{oidc.ScopeOpenID, "email", "profile"},
		},
		Verifier:  oidcProvider.Verifier(&oidc.Config{ClientID: googleClientID}),
		WebOrigin: webOrigin,
	}
	githubAuth := httpapi.GitHubAuth{
		OAuthConfig: oauth2.Config{
			ClientID:     cfg.GithubClientID,
			ClientSecret: cfg.GithubClientSecret,
			RedirectURL:  cfg.GithubRedirectURL,
			Endpoint:     github.Endpoint,
			Scopes:       []string{"read:user", "user:email"},
		},
		WebOrigin: webOrigin,
	}

	api := httpapi.AuthAPI{
		Users:  userStore,
		Tokens: tokenStore,
		OAuth:  oauthStore,
		Google: googleAuth,
		Github: githubAuth,

		JWT: auth.JWT{
			Secret: []byte(cfg.JWTSecret),
			Issuer: cfg.JWTIssuer,
		},
		AccessTTL:  15 * time.Minute,
		RefreshTTL: 14 * 24 * time.Hour,
		VerifyTTL:  24 * time.Hour,
		ResetTTL:   1 * time.Hour,
	}

	log.Println("Auth service running on :8080")
	log.Fatal(http.ListenAndServe(":8080", httpapi.Router(api, cfg.CORSAllowedOrigins)))
}
