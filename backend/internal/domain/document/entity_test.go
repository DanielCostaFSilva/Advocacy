package document

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDocument_ShouldCreateSuccessfully(t *testing.T) {
	caseID := uuid.New()

	doc, err := NewDocument(caseID, "Contrato Social", "Descrição do contrato", DocumentTypeContract, "contrato.pdf", "application/pdf", 1024, "cases/"+caseID.String()+"/contrato.pdf")

	require.NoError(t, err)
	require.NotNil(t, doc)
	assert.NotEmpty(t, doc.ID)
	assert.Equal(t, caseID, doc.CaseID)
	assert.Equal(t, "Contrato Social", doc.Name)
	assert.Equal(t, "Descrição do contrato", doc.Description)
	assert.Equal(t, DocumentTypeContract, doc.Type)
	assert.Equal(t, "contrato.pdf", doc.FileName)
	assert.Equal(t, "application/pdf", doc.MimeType)
	assert.Equal(t, int64(1024), doc.FileSize)
	assert.Equal(t, "cases/"+caseID.String()+"/contrato.pdf", doc.StorageKey)
	assert.False(t, doc.CreatedAt.IsZero())
	assert.False(t, doc.UpdatedAt.IsZero())
}

func TestNewDocument_ShouldReturnErrInvalidCaseID(t *testing.T) {
	_, err := NewDocument(uuid.Nil, "Contrato", "Desc", DocumentTypeContract, "file.pdf", "application/pdf", 1024, "key")
	assert.ErrorIs(t, err, ErrInvalidCaseID)
}

func TestNewDocument_ShouldReturnErrInvalidDocumentNameWhenEmpty(t *testing.T) {
	_, err := NewDocument(uuid.New(), "", "Desc", DocumentTypeContract, "file.pdf", "application/pdf", 1024, "key")
	assert.ErrorIs(t, err, ErrInvalidDocumentName)
}

func TestNewDocument_ShouldReturnErrInvalidDocumentNameWhenLessThan3(t *testing.T) {
	_, err := NewDocument(uuid.New(), "AB", "Desc", DocumentTypeContract, "file.pdf", "application/pdf", 1024, "key")
	assert.ErrorIs(t, err, ErrInvalidDocumentName)
}

func TestNewDocument_ShouldReturnErrInvalidDocumentType(t *testing.T) {
	_, err := NewDocument(uuid.New(), "Contrato", "Desc", "invalid_type", "file.pdf", "application/pdf", 1024, "key")
	assert.ErrorIs(t, err, ErrInvalidDocumentType)
}

func TestNewDocument_ShouldReturnErrInvalidFileName(t *testing.T) {
	_, err := NewDocument(uuid.New(), "Contrato", "Desc", DocumentTypeContract, "", "application/pdf", 1024, "key")
	assert.ErrorIs(t, err, ErrInvalidFileName)
}

func TestNewDocument_ShouldReturnErrInvalidMimeType(t *testing.T) {
	_, err := NewDocument(uuid.New(), "Contrato", "Desc", DocumentTypeContract, "file.pdf", "", 1024, "key")
	assert.ErrorIs(t, err, ErrInvalidMimeType)
}

func TestNewDocument_ShouldReturnErrInvalidFileSizeWhenZero(t *testing.T) {
	_, err := NewDocument(uuid.New(), "Contrato", "Desc", DocumentTypeContract, "file.pdf", "application/pdf", 0, "key")
	assert.ErrorIs(t, err, ErrInvalidFileSize)
}

func TestNewDocument_ShouldReturnErrInvalidFileSizeWhenNegative(t *testing.T) {
	_, err := NewDocument(uuid.New(), "Contrato", "Desc", DocumentTypeContract, "file.pdf", "application/pdf", -1, "key")
	assert.ErrorIs(t, err, ErrInvalidFileSize)
}

func TestNewDocument_ShouldReturnErrInvalidStorageKey(t *testing.T) {
	_, err := NewDocument(uuid.New(), "Contrato", "Desc", DocumentTypeContract, "file.pdf", "application/pdf", 1024, "")
	assert.ErrorIs(t, err, ErrInvalidStorageKey)
}

func TestNewDocument_ShouldAcceptAllTypes(t *testing.T) {
	caseID := uuid.New()
	types := []DocumentType{DocumentTypePetition, DocumentTypeContract, DocumentTypeEvidence, DocumentTypeCourtDecision, DocumentTypePowerOfAttorney, DocumentTypeOther}
	for _, dt := range types {
		doc, err := NewDocument(caseID, "Documento", "Desc", dt, "file.pdf", "application/pdf", 1024, "key")
		require.NoError(t, err)
		assert.Equal(t, dt, doc.Type)
	}
}

func TestNewDocument_ShouldAllowEmptyDescription(t *testing.T) {
	caseID := uuid.New()

	doc, err := NewDocument(caseID, "Contrato Social", "", DocumentTypeContract, "contrato.pdf", "application/pdf", 1024, "key")

	require.NoError(t, err)
	assert.Empty(t, doc.Description)
}
