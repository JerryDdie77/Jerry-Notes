package service

import (
	"context"
	"errors"
	"log"
	"time"
)

type NoteStorage interface {
	GetNote(ctx context.Context, noteID int64) (Note, error)
	UpdateNote(ctx context.Context, noteID int64, title, content string) error
	DeleteNote(ctx context.Context, noteID int64) error
	CreateNote(ctx context.Context, userID int64, title, content string) (int64, error)
	ListNotes(ctx context.Context, userID int64) ([]Note, error)
}

type NoteService struct {
	storage NoteStorage
}

func NewNoteService(s NoteStorage) *NoteService {
	return &NoteService{storage: s}
}

type Note struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NoteInput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type NoteOutput struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NoteList struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (n *NoteService) GetNote(ctx context.Context, noteID int64, userID int64) (NoteOutput, error) {
	note, err := n.storage.GetNote(ctx, noteID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return NoteOutput{}, ErrNotFound
		}
		log.Printf("storage GetNote failed: noteID=%d userID=%d err=%v\n", noteID, userID, err)
		return NoteOutput{}, ErrInternal
	}

	if note.UserID != userID {
		return NoteOutput{}, ErrForbidden
	}

	return NoteOutput{
		ID:        note.ID,
		Title:     note.Title,
		Content:   note.Content,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}, nil

}

func (n *NoteService) UpdateNote(ctx context.Context, userID, noteID int64, content, title string) error {

	if title == "" {
		return ErrEmptyTitle
	}

	note, err := n.storage.GetNote(ctx, noteID)

	switch {
	case errors.Is(err, ErrNotFound):
		return ErrNotFound
	case err != nil:

		log.Printf("storage GetNote failed: noteID=%d userID=%d err=%v\n", noteID, userID, err)

		return ErrInternal
	}

	userIDFromNote := note.UserID

	if userIDFromNote != userID {
		return ErrForbidden
	}

	err = n.storage.UpdateNote(ctx, noteID, title, content)

	switch {
	case errors.Is(err, ErrNotFound):
		return ErrNotFound
	case err != nil:
		log.Fatalf("storage UpdateNote failed: noteID=%d title=%s content=%s err=%v", noteID, title, content, err)
		return ErrInternal
	}

	return nil
}

func (n *NoteService) DeleteNote(ctx context.Context, userID, noteID int64) error {

	note, err := n.storage.GetNote(ctx, noteID)
	switch {
	case errors.Is(err, ErrNotFound):
		return ErrNotFound

	case err != nil:
		log.Printf("storage GetNote failed: noteID=%d userID=%d err=%v\n", noteID, userID, err)
		return ErrInternal
	}

	userIDFromNote := note.UserID

	if userIDFromNote != userID {
		return ErrForbidden
	}

	err = n.storage.DeleteNote(ctx, noteID)
	if err != nil {
		log.Printf("storage DeleteNote failed: noteID=%d err=%v\n", noteID, err)
		return ErrInternal
	}

	return nil
}

func (n *NoteService) CreateNote(ctx context.Context, userID int64, note NoteInput) (int64, error) {

	title := note.Title
	content := note.Content

	if title == "" {
		return 0, ErrEmptyTitle
	}

	noteID, err := n.storage.CreateNote(ctx, userID, title, content)

	switch {
	case errors.Is(err, ErrNotFound):
		return 0, ErrNotFound
	case err != nil:
		log.Printf("storage CreateNote failed: userID=%d title=%s content=%s err=%v\n", userID, title, content, err)
		return 0, ErrInternal
	}

	return noteID, nil

}

func (n *NoteService) ListNotes(ctx context.Context, userID int64) ([]NoteList, error) {

	notes := make([]NoteList, 0)
	fullNotes, err := n.storage.ListNotes(ctx, userID)
	if err != nil {
		log.Printf("storage ListNotes failed: userID=%d", userID)
		return nil, ErrInternal
	}

	for _, v := range fullNotes {
		nl := NoteList{
			ID:        v.ID,
			Title:     v.Title,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		}
		notes = append(notes, nl)
	}

	return notes, nil
}
