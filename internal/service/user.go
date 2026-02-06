package service

import (
	"context"
	"errors"
	"log"
)

type UserService struct {
	storage UserStorage
}

func NewUserService(storage UserStorage) *UserService {
	return &UserService{
		storage: storage,
	}
}

func (u *UserService) ChangeName(ctx context.Context, id int64, newName string) error {

	if newName == "" {
		return ErrEmptyName
	}

	user, err := u.storage.GetUserByID(ctx, id)

	// Check existing
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		}

		log.Printf("storage GetUserByID failed: id=%d err=%v", id, err)
		return ErrInternal
	}

	// Check if user changing HIS name, not others

	if id != user.ID {
		return ErrForbidden
	}

	err = u.storage.ChangeNameByID(ctx, id, newName)
	if err != nil {
		log.Printf("storage ChangeNameByID failed:")
	}

	return nil
}

func (u *UserService) DeleteUser(ctx context.Context, id int64) error {
	user, err := u.storage.GetUserByID(ctx, id)

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		}

		log.Printf("storage GetUserByID failed: id=%d err=%v", id, err)
		return ErrInternal
	}

	if id != user.ID {
		return ErrForbidden
	}

	// In this point we know that user exists and got the rules to delete yourself. So he can do this
	err = u.storage.DeleteUserByID(ctx, id)
	if err != nil {
		log.Printf("storage: DeleteUserByID failed: id=%d err=%v", id, err)
		return ErrInternal
	}
	return nil
}
