package service

import (
	"context"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserStorage interface {
	GetUser(ctx context.Context, email string) (User, error)
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

	// Check validity of the data and the strength of the password
	switch {

	case !validEmail(email):
		return 0, ErrInvalidEmail

	case !strongPassword(password):
		return 0, ErrWeakPassword

	case name == "":
		return 0, ErrEmptyName
	}

	// Checking if the email is busy
	existsEmail, err := a.storage.EmailExists(ctx, email)
	if err != nil {
		return 0, ErrInternal
	}

	if existsEmail {
		return 0, ErrEmailTaken
	}

	// Hashing the password
	hash, err := hashPassword(password)
	if errors.Is(err, bcrypt.ErrPasswordTooLong) {
		return 0, ErrPasswordTooLong
	}

	if err != nil {
		return 0, ErrInternal
	}

	// After all the checks, we are going to write the user data to the DB
	id, err := a.storage.CreateUser(ctx, name, hash, email)

	if err != nil {
		return 0, ErrInternal
	}

	// Returning the id of the newly created user
	return id, nil

}

// This func will be check login and password and return the user_id
func (a *AuthService) LoginUser(ctx context.Context, email, password string) (int64, error) {
	user, err := a.storage.GetUser(ctx, email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredential
		}

		return 0, ErrInternal
	}

	if user.IsBlocked {
		return 0, ErrUserBlocked
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return 0, ErrInvalidCredential
	}

	return user.ID, nil

}
