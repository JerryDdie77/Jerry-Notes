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

func (p PostgresUserStorage) GetUserByEmail(ctx context.Context, email string) (service.User, error) {
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

		if errors.Is(err, sql.ErrNoRows) {
			return service.User{}, service.ErrNotFound
		}

		return service.User{}, fmt.Errorf("scan: %w", err)
	}

	return u, nil
}

func (p PostgresUserStorage) GetUserByID(ctx context.Context, id int64) (service.User, error) {
	query := "SELECT id, user_name, created_at, password_hash, email, is_blocked FROM users WHERE id = $1"
	var u service.User
	err := p.db.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.Name, &u.CreatedAt, &u.PasswordHash, &u.Email, &u.IsBlocked)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return service.User{}, service.ErrNotFound
		}

		return service.User{}, fmt.Errorf("Scan: %w", err)
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

		if errors.Is(err, sql.ErrNoRows) {
			return 0, service.ErrNotFound
		}

		return 0, fmt.Errorf("scan: %w", err)
	}

	return id, nil
}

func (p *PostgresUserStorage) SavePendingUser(ctx context.Context, name, passwordHash, email, code string) error {
	query := "INSERT INTO pending_users(email, user_name, password_hash, code) VALUES($1, $2, $3, $4)"

	res, err := p.db.ExecContext(ctx, query, email, name, passwordHash, code)
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

func (p PostgresUserStorage) GetPendingUser(ctx context.Context, email string) (service.PendingUser, error) {
	query := "SELECT email, user_name, password_hash, code, created_at FROM pending_users WHERE email = $1"

	var u service.PendingUser

	row := p.db.QueryRowContext(ctx, query, email)
	err := row.Scan(&u.Email, &u.Name, &u.PasswordHash, &u.Code, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return service.PendingUser{}, service.ErrNotFound
		}

		return service.PendingUser{}, fmt.Errorf("scan: %w", err)
	}

	return u, nil
}

func (p *PostgresUserStorage) DeletePendingUser(ctx context.Context, email string) error {
	query := "DELETE FROM pending_users WHERE email = $1"
	res, err := p.db.ExecContext(ctx, query, email)
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

func (p *PostgresUserStorage) ChangeNameByID(ctx context.Context, id int64, newName string) error {
	query := "UPDATE users SET user_name = $1 WHERE id = $2"
	res, err := p.db.ExecContext(ctx, query, newName, id)
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

func (p *PostgresUserStorage) DeleteUserByID(ctx context.Context, id int64) error {
	query := "DELETE FROM users WHERE id = $1"
	res, err := p.db.ExecContext(ctx, query, id)
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
