package app

import "jerry-notes/internal/service"

type App struct {
	NoteService *service.NoteService
	UserService *service.UserService
	AuthService *service.AuthService
}

func New(noteService *service.NoteService, userService *service.UserService, authService *service.AuthService) *App {
	return &App{
		NoteService: noteService,
		UserService: userService,
		AuthService: authService,
	}
}
