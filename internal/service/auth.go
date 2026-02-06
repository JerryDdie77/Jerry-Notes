package service

import (
	"context"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

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

		log.Printf("storage EmailExists failed: email=%s err=%v", email, err)

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

		log.Printf("hash password failed: err=%v", err)

		return 0, ErrInternal
	}

	// After all the checks, we are going to write the user data to the DB
	id, err := a.storage.CreateUser(ctx, name, hash, email)

	if err != nil {

		log.Printf("storage CreateUser failed: name=%s email=%s err=%v", name, email, err)

		return 0, ErrInternal
	}

	// Returning the id of the newly created user
	return id, nil

}

// This func will be check login and password and return the user_id
func (a *AuthService) AuthentificateUser(ctx context.Context, email, password string) (int64, error) {
	user, err := a.storage.GetUserByEmail(ctx, email)

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return 0, ErrInvalidCredential
		}

		log.Printf("storage GetUser failed: email=%s err=%v", email, err)

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
