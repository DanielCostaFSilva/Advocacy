package client

import "errors"

var (
	ErrClientAlreadyExists = errors.New("client already exists")
	ErrInvalidInput        = errors.New("invalid input")
	ErrClientNotFound      = errors.New("client not found")
	ErrInvalidClientID     = errors.New("invalid client id")
)
