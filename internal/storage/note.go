package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"jerry-notes/internal/service"
)

type PosgtresNoteStorage struct {
	db *sql.DB
}

func (p PosgtresNoteStorage) GetUserID(ctx context.Context, noteID int) (int, error) {
	query := "SELECT user_id FROM notes WHERE id = $1"

	var userID int

	if err := p.db.QueryRowContext(ctx, query, noteID).Scan(&userID); err != nil {
		return 0, fmt.Errorf("Scan: %w", err)
	}

	return userID, nil
}

func (p PosgtresNoteStorage) GetNote(ctx context.Context, noteID int) (service.Note, error) {
	query := "SELECT id, user_id, title, content, created_at, updated_at FROM notes WHERE id = $1"

	var n service.Note

	if err := p.db.QueryRowContext(ctx, query, noteID).Scan(&n.ID, &n.UserID, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt); err != nil {
		return service.Note{}, fmt.Errorf("Scan: %w", err)
	}

	return n, nil
}

func (p *PosgtresNoteStorage) UpdateNote(ctx context.Context, noteID int, title, content string) error {
	query := "UPDATE notes SET title = $1, content = $2 WHERE id = $3"

	res, err := p.db.ExecContext(ctx, query, title, content, noteID)
	if err != nil {
		return fmt.Errorf("ExecContext: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("RowsAffected: %w", err)
	}

	if rowsAffected == 0 {
		return service.ErrNotFound
	}

	return nil
}

func (p *PosgtresNoteStorage) DeleteNote(ctx context.Context, noteID int) error {
	query := "DELETE FROM notes WHERE id = $1"
	_, err := p.db.ExecContext(ctx, query, noteID)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("Note with id %d not found", noteID)
		}

		return fmt.Errorf("Exec: %w", err)
	}

	return nil
}

func (p *PosgtresNoteStorage) CreateNote(ctx context.Context, userID int, title, content string) (int, error) {
	query := "INSERT INTO notes(user_id, title, content) VALUES($1, $2, $3) RETURNING id"

	var noteID int

	if err := p.db.QueryRowContext(ctx, query, userID, title, content).Scan(&noteID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, service.ErrNotFound
		}
		return 0, fmt.Errorf("Scan: %w", err)
	}

	return noteID, nil
}

func (p *PosgtresNoteStorage) ListNotes(ctx context.Context, userID int) ([]service.Note, error) {
	query := "SELECT id, user_id, title, content, created_at, updated_at FROM notes WHERE user_id = $1"

	rows, err := p.db.QueryContext(ctx, query, userID)
	if err != nil {
		return []service.Note{}, fmt.Errorf("QueryContext: %w", err)
	}
	notes := make([]service.Note, 0)
	for rows.Next() {
		var n service.Note
		if err := rows.Scan(&n.ID, &n.UserID, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return []service.Note{}, fmt.Errorf("Scan: %w", err)
		}
		notes = append(notes, n)
	}

	if err := rows.Err(); err != nil {
		return []service.Note{}, fmt.Errorf("rows Err: %w", err)
	}

	return notes, nil

}
