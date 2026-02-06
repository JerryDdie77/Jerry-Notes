package service

import "context"

type UserStorage interface {
	GetUserByID(ctx context.Context, id int64) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	EmailExists(ctx context.Context, email string) (bool, error)
	CreateUser(ctx context.Context, name, passwordHash, email string) (int64, error)
	ChangeNameByID(ctx context.Context, id int64, newName string) error
	DeleteUserByID(ctx context.Context, id int64) error
}
