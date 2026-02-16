package jwt_test

import (
	"jerry-notes/config"
	"jerry-notes/internal/jwt"
	"testing"
	"time"
)

func TestJWT(t *testing.T) {
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	const id = 12
	jwtManager := jwt.NewManager(cfg.JWTSecret, 1*time.Second)
	token, err := jwtManager.GenerateToken(id)
	if err != nil {
		t.Fatal(err)
	}

	userID, err := jwtManager.ParseUserID(token)
	if err != nil {
		t.Fatal(err)
	}

	if userID != id {
		t.Fatalf("ID %d is not equal to %d", userID, id)
	}

}
