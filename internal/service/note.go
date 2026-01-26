package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type NoteStorage interface {
	GetUserID(ctx context.Context, noteID int) (int, error)
	GetNote(ctx context.Context, noteID int) (Note, error)
	UpdateNote(ctx context.Context, noteID int, title, content string) error
	DeleteNote(ctx context.Context, noteID int) error
	CreateNote(ctx context.Context, userID int, title, content string) (int, error)
	ListNotes(ctx context.Context, userID int) ([]Note, error)
}

type NoteService struct {
	storage NoteStorage
}

func NewNoteService(s NoteStorage) *NoteService {
	return &NoteService{storage: s}
}

func (n *NoteService) GetNote(ctx context.Context, noteID int, userID int) (NoteOutput, error) {
	note, err := n.storage.GetNote(ctx, noteID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NoteOutput{}, ErrNotFound
		}
		return NoteOutput{}, fmt.Errorf("storage GetNote: %w", err)
	}

	if note.UserID != userID {
		return NoteOutput{}, ErrForbidden
	}

	return NoteOutput{
		ID:        noteID,
		Title:     note.Title,
		Content:   note.Content,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}, nil

}

func (n *NoteService) UpdateNote(ctx context.Context, userID, noteID int, content, title string) error {

	if title == "" {
		return ErrEmptyTitle
	}

	userIDFromNote, err := n.storage.GetUserID(ctx, noteID)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return ErrNotFound
	case err != nil:
		return fmt.Errorf("storage GetUserID: %w", err)
	}

	if userIDFromNote != userID {
		return ErrForbidden
	}

	err = n.storage.UpdateNote(ctx, noteID, title, content)
	switch {
	case errors.Is(err, ErrNotFound):
		return ErrNotFound
	case err != nil:
		return fmt.Errorf("storage UpdateNote: %w", err)
	}

	return nil
}

func (n *NoteService) DeleteNote(ctx context.Context, userID, noteID int) error {

	userIDFromNote, err := n.storage.GetUserID(ctx, noteID)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return ErrNotFound

	case err != nil:
		return fmt.Errorf("storage GetUserID: %w", err)
	}

	if userIDFromNote != userID {
		return ErrForbidden
	}

	err = n.storage.DeleteNote(ctx, noteID)
	if err != nil {
		return fmt.Errorf("storage DeleteNote: %w", err)
	}

	return nil
}

func (n *NoteService) CreateNote(ctx context.Context, userID int, note NoteInput) (int, error) {

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
		return 0, fmt.Errorf("storage CreateNote: %w", err)
	}

	return noteID, nil

}

func (n *NoteService) ListNotes(ctx context.Context, userID int) ([]NoteList, error) {

	notes := make([]NoteList, 0)
	fullNotes, err := n.storage.ListNotes(ctx, userID)
	if err != nil {
		return []NoteList{}, fmt.Errorf("storage ListNotes: %w", err)
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
