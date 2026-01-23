package config

import (
	"log"
	"os"
	"strings"
)

type Config struct {
	DatabaseURL        string
	JWTSecret          string
	JWTIssuer          string
	CORSAllowedOrigins []string
}

func Load() Config {
	cfg := Config{
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		JWTSecret:          os.Getenv("JWT_SECRET"),
		JWTIssuer:          os.Getenv("JWT_ISSUER"),
		CORSAllowedOrigins: parseEnvList(os.Getenv("CORS_ALLOWED_ORIGINS")),
	}

	if cfg.DatabaseURL == "" || cfg.JWTSecret == "" {
		log.Fatal("Missing required environment variables")
	}

	return cfg
}

func parseEnvList(value string) []string {
	if value == "" {
		return []string{"*"}
	}

	raw := strings.Split(value, ",")
	out := make([]string, 0, len(raw))
	for _, item := range raw {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		out = append(out, item)
	}
	if len(out) == 0 {
		return []string{"*"}
	}
	return out
}
