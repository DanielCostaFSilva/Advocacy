package legalcase

import "errors"

var (
	ErrClientNotFound    = errors.New("client not found")
	ErrCaseAlreadyExists = errors.New("case already exists")
	ErrInvalidClientID   = errors.New("invalid client id")
	ErrInvalidInput      = errors.New("invalid input")
)
