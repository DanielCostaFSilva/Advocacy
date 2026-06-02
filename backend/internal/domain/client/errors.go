package client

import "errors"

var (
	ErrInvalidName  = errors.New("name must be at least 3 characters")
	ErrInvalidCPF   = errors.New("cpf must have exactly 11 digits")
	ErrInvalidEmail = errors.New("invalid email format")
)
