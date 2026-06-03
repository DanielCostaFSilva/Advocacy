package contract

import "errors"

var (
	ErrClientNotFound      = errors.New("client not found")
	ErrCaseNotFound        = errors.New("case not found")
	ErrInvalidClientID     = errors.New("invalid client id")
	ErrInvalidCaseID       = errors.New("invalid case id")
	ErrInvalidInput        = errors.New("invalid input")
	ErrContractPersistence = errors.New("contract persistence error")
	ErrTimelinePersistence = errors.New("timeline persistence error")
	ErrContractNotFound    = errors.New("contract not found")
	ErrInvalidContractID   = errors.New("invalid contract id")
	ErrContractAlreadyClosed = errors.New("contract already closed")
)
