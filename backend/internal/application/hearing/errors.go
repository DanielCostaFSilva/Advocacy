package hearing

import "errors"

var (
	ErrCaseNotFound  = errors.New("case not found")
	ErrInvalidCaseID = errors.New("invalid case id")
	ErrInvalidInput  = errors.New("invalid input")
)
