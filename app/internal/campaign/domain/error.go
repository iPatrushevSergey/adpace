package domain

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrBadInput      = errors.New("bad input")
	ErrConflict      = errors.New("conflict")
)
