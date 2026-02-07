package configs

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret  string
	DBURL      string
	GmailToken string
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

	gmailToken := os.Getenv("GMAIL_TOKEN")

	if gmailToken == "" {
		return &Config{}, fmt.Errorf("GMAIL_TOKEN is required")
	}

	return &Config{JWTSecret: secret, DBURL: dbURL, GmailToken: gmailToken}, nil
}
