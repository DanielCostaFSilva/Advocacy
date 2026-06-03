package contract

import "errors"

var (
	ErrInvalidClientID     = errors.New("client id is required")
	ErrInvalidCaseID       = errors.New("case id is required")
	ErrInvalidTitle        = errors.New("title must be at least 3 characters")
	ErrInvalidContractType = errors.New("invalid contract type")
	ErrInvalidAmount       = errors.New("amount must be greater than zero")
	ErrInvalidStartDate    = errors.New("start date is required")
	ErrInvalidEndDate      = errors.New("end date must be after start date")
)
