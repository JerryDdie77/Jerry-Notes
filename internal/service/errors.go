package service

import "errors"

var ErrForbidden = errors.New("forbidden")
var ErrNotFound = errors.New("not found")
var ErrEmptyTitle = errors.New("title can not be empty")
