package hearing

import "errors"

var (
	ErrInvalidCaseID       = errors.New("case id is required")
	ErrInvalidTitle        = errors.New("title must be at least 5 characters")
	ErrInvalidHearingType  = errors.New("invalid hearing type")
	ErrInvalidLocation     = errors.New("location is required")
	ErrInvalidScheduledAt  = errors.New("scheduled date must be in the future")
)
