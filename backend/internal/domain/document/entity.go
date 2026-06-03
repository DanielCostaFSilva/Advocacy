package document

import (
	"time"

	"github.com/google/uuid"
)

type Document struct {
	ID          uuid.UUID
	CaseID      uuid.UUID
	Name        string
	Description string
	Type        DocumentType
	FileName    string
	MimeType    string
	FileSize    int64
	StorageKey  string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewDocument(
	caseID uuid.UUID,
	name string,
	description string,
	docType DocumentType,
	fileName string,
	mimeType string,
	fileSize int64,
	storageKey string,
) (*Document, error) {
	if caseID == uuid.Nil {
		return nil, ErrInvalidCaseID
	}

	if len(name) < 3 {
		return nil, ErrInvalidDocumentName
	}

	switch docType {
	case DocumentTypePetition, DocumentTypeContract, DocumentTypeEvidence, DocumentTypeCourtDecision, DocumentTypePowerOfAttorney, DocumentTypeOther:
	default:
		return nil, ErrInvalidDocumentType
	}

	if fileName == "" {
		return nil, ErrInvalidFileName
	}

	if mimeType == "" {
		return nil, ErrInvalidMimeType
	}

	if fileSize <= 0 {
		return nil, ErrInvalidFileSize
	}

	if storageKey == "" {
		return nil, ErrInvalidStorageKey
	}

	now := time.Now()

	return &Document{
		ID:          uuid.New(),
		CaseID:      caseID,
		Name:        name,
		Description: description,
		Type:        docType,
		FileName:    fileName,
		MimeType:    mimeType,
		FileSize:    fileSize,
		StorageKey:  storageKey,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}
