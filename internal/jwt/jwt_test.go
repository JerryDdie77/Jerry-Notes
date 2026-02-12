package jwt_test

import (
	"jerry-notes/internal/jwt"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func TestJWT(t *testing.T) {
	_ = godotenv.Load()
	const id = 12
	jwtManager := jwt.NewManager(os.Getenv("JWT_SECRET"), 1*time.Second)
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
