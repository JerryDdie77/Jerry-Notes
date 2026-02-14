package app

import (
	"jerry-notes/internal/jwt"
	"jerry-notes/internal/service"
)

type App struct {
	NoteService *service.NoteService
	UserService *service.UserService
	AuthService *service.AuthService
	JWTManager  *jwt.Manager
}

func New(noteService *service.NoteService, userService *service.UserService, authService *service.AuthService, jwtManager *jwt.Manager) *App {
	return &App{
		NoteService: noteService,
		UserService: userService,
		AuthService: authService,
		JWTManager:  jwtManager,
	}
}
