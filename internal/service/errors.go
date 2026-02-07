package service

import "errors"

var ErrForbidden = errors.New("forbidden")
var ErrNotFound = errors.New("not found")
var ErrEmptyTitle = errors.New("title can not be empty")
var ErrInvalidEmail = errors.New("email is invalid")
var ErrWeakPassword = errors.New("password is weak")
var ErrInternal = errors.New("internal server error")
var ErrInvalidName = errors.New("name is invalid")
var ErrEmailTaken = errors.New("email is already taken")
var ErrPasswordTooLong = errors.New("password can not be longer than 72 symbols")
var ErrInvalidCredential = errors.New("password isn't correct")
var ErrUserBlocked = errors.New("user is blocked")
var ErrInvalidCode = errors.New("code isn't correct")
var ErrCodeExpired = errors.New("code is expired")
