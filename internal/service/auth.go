package service

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int64
	Name         string
	PasswordHash string
	Email        string
	IsBlocked    bool
	RegisteredAt time.Time
}

type UserStorage interface {
	GetUserByName(ctx context.Context, name string) (User, error)
	EmailExists(ctx context.Context, email string) (bool, error)
	CreateUser(ctx context.Context, name, password_hash, email string) (int64, error)
}

type AuthService struct {
	storage UserStorage
}

func NewAuthService(storage UserStorage) *AuthService {
	return &AuthService{
		storage: storage,
	}
}

func (a *AuthService) RegisterUser(ctx context.Context, name, password, email string) (int64, error) {

	switch {

	case !validEmail(email):
		return 0, ErrInvalidEmail

	case !strongPassword(password):
		return 0, ErrWeakPassword

	case name == "":
		return 0, ErrEmptyName
	}

	existsEmail, err := a.storage.EmailExists(ctx, email)
	if err != nil {
		return 0, ErrInternal
	}

	if existsEmail {
		return 0, ErrEmailTaken
	}

	hash, err := hashPassword(password)
	if errors.Is(err, bcrypt.ErrPasswordTooLong) {
		return 0, ErrPasswordTooLong
	}

	if err != nil {
		return 0, ErrInternal
	}

	id, err := a.storage.CreateUser(ctx, name, hash, email)

	if err != nil {
		return 0, ErrInternal
	}

	return id, nil

}
