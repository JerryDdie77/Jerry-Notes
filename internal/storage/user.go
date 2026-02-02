package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"jerry-notes/internal/service"
)

type PostgresUserStorage struct {
	db *sql.DB
}

func NewPostgresUserStorage(db *sql.DB) *PostgresUserStorage {
	return &PostgresUserStorage{
		db: db,
	}
}

func (p PosgtresNoteStorage) GetUser(ctx context.Context, email string) (service.User, error) {
	query := "SELECT id, user_name, created_at, password_hash, email, is_blocked FROM users WHERE email = $1"

	var u service.User

	if err := p.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID,
		&u.Name,
		&u.CreatedAt,
		&u.PasswordHash,
		&u.Email,
		&u.IsBlocked,
	); err != nil {
		return service.User{}, fmt.Errorf("scan: %w", err)
	}
	return u, nil
}

func (p PostgresUserStorage) EmailExists(ctx context.Context, email string) (bool, error) {
	var tmp int64

	query := "SELECT id FROM users WHERE email = $1"
	if err := p.db.QueryRowContext(ctx, query, email).Scan(&tmp); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, fmt.Errorf("scan: %w", err)
	}

	return true, nil
}

func (p *PostgresUserStorage) CreateUser(ctx context.Context, name, passwordHash, email string) (int64, error) {
	query := "INSERT INTO users(user_name, password_hash, email) VALUES($1, $2, $3) RETURNING id"

	var id int64

	if err := p.db.QueryRowContext(ctx, query, name, passwordHash, email).Scan(&id); err != nil {
		return 0, fmt.Errorf("scan: %w", err)
	}

	return id, nil
}
