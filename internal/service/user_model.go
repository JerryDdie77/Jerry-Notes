package service

import "time"

type User struct {
	ID           int64
	Name         string
	CreatedAt    time.Time
	PasswordHash string
	Email        string
	IsBlocked    bool
}

type PendingUser struct {
	Email        string
	Name         string
	PasswordHash string
	Code         string
	CreatedAt    time.Time
}
