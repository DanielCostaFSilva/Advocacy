package client

import "errors"

var (
	ErrClientAlreadyExists = errors.New("client already exists")
	ErrInvalidInput        = errors.New("invalid input")
)
