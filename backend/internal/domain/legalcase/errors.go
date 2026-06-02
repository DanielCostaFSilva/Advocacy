package legalcase

import "errors"

var (
	ErrInvalidClientID    = errors.New("client id is required")
	ErrInvalidCaseNumber  = errors.New("case number must be at least 5 characters")
	ErrInvalidCaseTitle   = errors.New("case title must be at least 5 characters")
	ErrInvalidCourt       = errors.New("court is required")
	ErrCaseAlreadyExists  = errors.New("case already exists")
)
