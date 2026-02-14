package service

import (
	"context"
	"errors"
	"jerry-notes/internal/jwt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	storage      UserStorage
	emailService *EmailService
	jwtManager   *jwt.Manager
	codeTTL      time.Duration
}

func NewAuthService(storage UserStorage, emailService *EmailService, jwtManager *jwt.Manager, codeTTL time.Duration) *AuthService {
	return &AuthService{
		storage:      storage,
		emailService: emailService,
		jwtManager:   jwtManager,
		codeTTL:      codeTTL,
	}
}

// Saves user data to temp table and sends confirmation code to email
func (a *AuthService) StartRegistration(ctx context.Context, name, password, email string) error {

	// Validate input data
	switch {
	case !validEmail(email):
		return ErrInvalidEmail
	case !strongPassword(password):
		return ErrWeakPassword
	case !validName(name):
		return ErrInvalidName
	}

	// Check if email already exists
	existsEmail, err := a.storage.EmailExists(ctx, email)
	if err != nil {
		log.Printf("storage EmailExists failed: email=%s err=%v", email, err)
		return ErrInternal
	}

	if existsEmail {
		return ErrEmailTaken
	}

	// Hash password
	hash, err := hashPassword(password)
	if errors.Is(err, bcrypt.ErrPasswordTooLong) {
		return ErrPasswordTooLong
	}

	if err != nil {
		log.Printf("hash password failed: err=%v", err)
		return ErrInternal
	}

	// Generate confirmation code
	code, err := generateCode()
	if err != nil {
		log.Printf("generateCode failed: %v", err)
		return ErrInternal
	}

	// Save to pending users table
	err = a.storage.SavePendingUser(ctx, name, hash, email, code)
	if err != nil {
		log.Printf("storage SavePendingUser failed: name=%s email=%s err=%v", name, email, err)
		return ErrInternal
	}

	// Send confirmation email
	err = a.emailService.SendConfirmationCode(ctx, email, code)
	if err != nil {
		log.Printf("emailService SendConfirmationCode failed: email=%s", email)
		return ErrInternal
	}

	return nil
}

// Confirms user registration by code and returns a JWT access token
func (a *AuthService) VerifyCode(ctx context.Context, email, code string) (string, error) {

	// Fetch pending user data
	user, err := a.storage.GetPendingUser(ctx, email)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return "", ErrNotFound
		}
		log.Printf("storage GetPendingUser failed: email=%s err=%v", email, err)
		return "", ErrInternal
	}

	// Validate confirmation code and expiration
	if user.Code != code {
		return "", ErrInvalidCode
	}
	if time.Since(user.CreatedAt) > a.codeTTL {
		return "", ErrCodeExpired
	}

	// Create user in main table first (ensures user exists even if cleanup fails)
	id, err := a.storage.CreateUser(ctx, user.Name, user.PasswordHash, user.Email)
	if err != nil {
		log.Printf("storage CreateUser failed: name=%s email=%s err=%v", user.Name, user.Email, err)
		return "", ErrInternal
	}

	// Cleanup pending data (best effort)
	if err := a.storage.DeletePendingUser(ctx, user.Email); err != nil {
		log.Printf("storage DeletePendingUser failed: email=%s err=%v", user.Email, err)
	}

	// Generate JWT token
	jwtToken, err := a.jwtManager.GenerateToken(id)
	if err != nil {
		log.Printf("jwtManager GenerateToken failed: id=%d err=%v", id, err)
		return "", ErrInternal
	}
	return jwtToken, nil
}

// AuthenticateUser checks credentials and returns a JWT access token
func (a *AuthService) AuthenticateUser(ctx context.Context, email, password string) (string, error) {

	// Fetch user by email
	user, err := a.storage.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return "", ErrInvalidCredential
		}
		log.Printf("storage GetUserByEmail failed: email=%s err=%v", email, err)
		return "", ErrInternal
	}

	// Check if user is blocked
	if user.IsBlocked {
		return "", ErrUserBlocked
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", ErrInvalidCredential
	}

	// Generate JWT token
	jwtToken, err := a.jwtManager.GenerateToken(user.ID)
	if err != nil {
		log.Printf("jwtManager GenerateToken failed: id=%d err=%v", user.ID, err)
		return "", ErrInternal
	}

	return jwtToken, nil
}
