package timeline

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	documentDomain "legalflow/internal/domain/document"
)

func TestNewDocumentUploadedEvent_ShouldGenerateDOCUMENTUPLOADED(t *testing.T) {
	doc := &documentDomain.Document{
		ID: uuid.New(), CaseID: uuid.New(), Name: "Contrato de Honorários",
		Type: documentDomain.DocumentTypeContract, FileName: "contrato.pdf",
		MimeType: "application/pdf", FileSize: 1024, StorageKey: "key",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}

	event := NewDocumentUploadedEvent(doc.CaseID, doc)

	require.NotNil(t, event)
	assert.Equal(t, EventDocumentUploaded, event.Type)
	assert.Equal(t, doc.CaseID, event.CaseID)
	assert.Contains(t, event.Description, "Contrato de Honorários")
}

func TestNewDocumentUploadedEvent_ShouldPopulateMetadata(t *testing.T) {
	doc := &documentDomain.Document{
		ID: uuid.New(), CaseID: uuid.New(), Name: "Contrato",
		Type: documentDomain.DocumentTypeContract, FileName: "contrato.pdf",
		MimeType: "application/pdf", FileSize: 1024, StorageKey: "key",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}

	event := NewDocumentUploadedEvent(doc.CaseID, doc)

	require.NotNil(t, event)
	require.NotEmpty(t, event.Metadata)

	var meta DocumentMetadata
	err := json.Unmarshal(event.Metadata, &meta)
	require.NoError(t, err)
	assert.Equal(t, doc.ID.String(), meta.DocumentID)
	assert.Equal(t, string(doc.Type), meta.DocumentType)
	assert.Equal(t, doc.FileName, meta.FileName)
}

func TestNewDocumentDownloadedEvent_ShouldGenerateDOCUMENTDOWNLOADED(t *testing.T) {
	doc := &documentDomain.Document{
		ID: uuid.New(), CaseID: uuid.New(), Name: "Contrato",
		Type: documentDomain.DocumentTypeContract, FileName: "contrato.pdf",
		MimeType: "application/pdf", FileSize: 1024, StorageKey: "key",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}

	event := NewDocumentDownloadedEvent(doc.CaseID, doc)

	require.NotNil(t, event)
	assert.Equal(t, EventDocumentDownloaded, event.Type)
	assert.Equal(t, doc.CaseID, event.CaseID)
	assert.Contains(t, event.Description, "Contrato")
}

func TestNewDocumentDownloadedEvent_ShouldPopulateMetadata(t *testing.T) {
	doc := &documentDomain.Document{
		ID: uuid.New(), CaseID: uuid.New(), Name: "Relatório",
		Type: documentDomain.DocumentTypePetition, FileName: "peticao.pdf",
		MimeType: "application/pdf", FileSize: 2048, StorageKey: "key2",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}

	event := NewDocumentDownloadedEvent(doc.CaseID, doc)

	require.NotNil(t, event)
	require.NotEmpty(t, event.Metadata)

	var meta DocumentMetadata
	err := json.Unmarshal(event.Metadata, &meta)
	require.NoError(t, err)
	assert.Equal(t, doc.ID.String(), meta.DocumentID)
	assert.Equal(t, string(doc.Type), meta.DocumentType)
	assert.Equal(t, doc.FileName, meta.FileName)
}

func TestNewDocumentEvent_DescriptionShouldContainDocumentName(t *testing.T) {
	doc := &documentDomain.Document{
		ID: uuid.New(), CaseID: uuid.New(), Name: "Procuração",
		Type: documentDomain.DocumentTypePowerOfAttorney, FileName: "procuracao.pdf",
		MimeType: "application/pdf", FileSize: 512, StorageKey: "key",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}

	t.Run("uploaded", func(t *testing.T) {
		event := NewDocumentUploadedEvent(doc.CaseID, doc)
		assert.Contains(t, event.Description, "Procuração")
		assert.Contains(t, event.Description, "enviado")
	})

	t.Run("downloaded", func(t *testing.T) {
		event := NewDocumentDownloadedEvent(doc.CaseID, doc)
		assert.Contains(t, event.Description, "Procuração")
		assert.Contains(t, event.Description, "acessado")
	})
}
