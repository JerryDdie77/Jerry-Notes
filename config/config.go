package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret         string
	DBURL             string
	MailToken         string
	JWTAccessTokenTTL time.Duration
	CodeTTL           time.Duration
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return &Config{}, fmt.Errorf("load: %w", err)
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return &Config{}, fmt.Errorf("JWT_SECRET is required")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		return &Config{}, fmt.Errorf("DB_URL is required")
	}

	mailToken := os.Getenv("MAIL_TOKEN")
	if mailToken == "" {
		return &Config{}, fmt.Errorf("MAIL_TOKEN is required")
	}

	ttlStr := os.Getenv("JWT_ACCESS_TOKEN_TTL")
	if ttlStr == "" {
		ttlStr = "30m"
	}

	jwtTTL, err := time.ParseDuration(ttlStr)
	if err != nil {
		return &Config{}, fmt.Errorf("parse JWT_ACCESS_TOKEN_TTL: %w", err)
	}

	codeTTLStr := os.Getenv("CODE_TTL")
	if codeTTLStr == "" {
		codeTTLStr = "5m"
	}

	codeTTL, err := time.ParseDuration(codeTTLStr)
	if err != nil {
		return &Config{}, fmt.Errorf("parse CODE_TTL: %w", err)
	}

	return &Config{
		JWTSecret:         secret,
		DBURL:             dbURL,
		MailToken:         mailToken,
		JWTAccessTokenTTL: jwtTTL,
		CodeTTL:           codeTTL,
	}, nil
}
