package main

import (
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"

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

	api := httpapi.AuthAPI{
		Users:  userStore,
		Tokens: tokenStore,
		JWT: auth.JWT{
			Secret: []byte(cfg.JWTSecret),
			Issuer: cfg.JWTIssuer,
		},
		AccessTTL:  15 * time.Minute,
		RefreshTTL: 14 * 24 * time.Hour,
	}

	log.Println("Auth service running on :8080")
	log.Fatal(http.ListenAndServe(":8080", httpapi.Router(api, cfg.CORSAllowedOrigins)))
}
