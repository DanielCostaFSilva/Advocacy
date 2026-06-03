package document

import "errors"

var (
	ErrCaseNotFound        = errors.New("case not found")
	ErrInvalidCaseID       = errors.New("invalid case id")
	ErrInvalidInput        = errors.New("invalid input")
	ErrDocumentPersistence  = errors.New("document persistence error")
	ErrTimelinePersistence  = errors.New("timeline persistence error")
	ErrInvalidDocumentID   = errors.New("invalid document id")
	ErrDocumentNotFound    = errors.New("document not found")
)
