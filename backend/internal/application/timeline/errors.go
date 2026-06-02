package timeline

import "errors"

var (
	ErrInvalidCaseID = errors.New("invalid case id")
	ErrCaseNotFound  = errors.New("case not found")
)
