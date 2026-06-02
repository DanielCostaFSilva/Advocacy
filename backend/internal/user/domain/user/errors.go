package user

import "errors"

var (
	ErrInvalidName         = errors.New("name must be at least 3 characters")
	ErrInvalidEmail        = errors.New("invalid email format")
	ErrInvalidPasswordHash = errors.New("password hash cannot be empty")
)
