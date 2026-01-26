package service

import "time"

type Note struct {
	ID        int
	UserID    int
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type NoteInput struct {
	Title   string
	Content string
}

type NoteOutput struct {
	ID        int
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type NoteList struct {
	ID        int
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
