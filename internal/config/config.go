package config

import "os"

type Config struct {
	HTTPAddr      string
	DatabaseURL   string
	AuthEnabled   bool
	OIDCIssuerURL string
	OIDCClientID  string
}

func FromEnv() Config {
	return Config{
		HTTPAddr:      getEnv("HTTP_ADDR", ":8080"),
		DatabaseURL:   getEnv("DATABASE_URL", "host=localhost user=soldat password=p dbname=soldat port=5432 sslmode=disable"),
		AuthEnabled:   getEnv("AUTH_ENABLED", "false") == "true",
		OIDCIssuerURL: os.Getenv("OIDC_ISSUER_URL"),
		OIDCClientID:  os.Getenv("OIDC_CLIENT_ID"),
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
