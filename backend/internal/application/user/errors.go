package user

import "errors"

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidInput       = errors.New("invalid input")
	ErrHashPassword       = errors.New("failed to hash password")
)
