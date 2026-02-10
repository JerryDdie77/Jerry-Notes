package jwt_test

import (
	"jerry-notes/internal/jwt"
	"testing"
	"time"
)

func TestJWT(t *testing.T) {
	const id = 12
	jwtManager := jwt.NewManager("jx8Equ5m8JQ3GRN2Mvpt69738qqLeLE7U6yXBDw59dEnHjmY18Uwm6McgF8I8EeY6LJCG2R3lozIxEeuYTvLdz", 1*time.Second)
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
