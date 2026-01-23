package config

import (
	"log"
	"os"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	JWTIssuer   string
}

func Load() Config {
	cfg := Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		JWTIssuer:   os.Getenv("JWT_ISSUER"),
	}

	if cfg.DatabaseURL == "" || cfg.JWTSecret == "" {
		log.Fatal("Missing required environment variables")
	}

	return cfg
}
