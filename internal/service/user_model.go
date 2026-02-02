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
