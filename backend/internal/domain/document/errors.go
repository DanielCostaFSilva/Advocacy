package document

import "errors"

var (
	ErrInvalidCaseID      = errors.New("case id is required")
	ErrInvalidDocumentName = errors.New("document name must be at least 3 characters")
	ErrInvalidDocumentType = errors.New("invalid document type")
	ErrInvalidFileName     = errors.New("file name is required")
	ErrInvalidMimeType     = errors.New("mime type is required")
	ErrInvalidFileSize     = errors.New("file size must be greater than zero")
	ErrInvalidStorageKey   = errors.New("storage key is required")
)
