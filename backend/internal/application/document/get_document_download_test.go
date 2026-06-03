package document

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	documentDomain "legalflow/internal/domain/document"
	timelineDomain "legalflow/internal/domain/timeline"
)

type mockDownloadURLProvider struct {
	GenerateFunc func(document *documentDomain.Document) string
	ExpiresAtFunc func() time.Time
}

func (m *mockDownloadURLProvider) Generate(document *documentDomain.Document) string {
	if m.GenerateFunc != nil {
		return m.GenerateFunc(document)
	}
	return ""
}

func (m *mockDownloadURLProvider) ExpiresAt() time.Time {
	if m.ExpiresAtFunc != nil {
		return m.ExpiresAtFunc()
	}
	return time.Time{}
}

func TestGetDocumentDownload_ShouldReturnTemporaryURL(t *testing.T) {
	docID := uuid.New()
	now := time.Now()

	docs := &documentDomain.MockRepository{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*documentDomain.Document, error) {
			return &documentDomain.Document{
				ID: docID, CaseID: uuid.New(), Name: "Contrato de Honorários",
				Type: documentDomain.DocumentTypeContract, FileName: "contrato.pdf",
				MimeType: "application/pdf", FileSize: 1024, StorageKey: "key",
				CreatedAt: now, UpdatedAt: now,
			}, nil
		},
	}
	timeline := &timelineDomain.MockRepository{
		CreateFunc: func(ctx context.Context, event *timelineDomain.TimelineEvent) error {
			return nil
		},
	}
	provider := &mockDownloadURLProvider{
		GenerateFunc: func(document *documentDomain.Document) string {
			return "https://temporary-download.local/documents/" + docID.String()
		},
		ExpiresAtFunc: func() time.Time {
			return now.Add(30 * time.Minute)
		},
	}

	uc := NewGetDocumentDownloadUseCase(docs, timeline, provider)
	output, err := uc.Execute(context.Background(), GetDocumentDownloadInput{
		DocumentID: docID.String(),
	})

	require.NoError(t, err)
	require.NotNil(t, output)
	assert.Equal(t, docID.String(), output.DocumentID)
	assert.Equal(t, "contrato.pdf", output.FileName)
	assert.Equal(t, "https://temporary-download.local/documents/"+docID.String(), output.DownloadURL)
	assert.WithinDuration(t, now.Add(30*time.Minute), output.ExpiresAt, time.Second)
}

func TestGetDocumentDownload_ShouldReturnErrInvalidDocumentID(t *testing.T) {
	uc := NewGetDocumentDownloadUseCase(
		&documentDomain.MockRepository{},
		&timelineDomain.MockRepository{},
		&mockDownloadURLProvider{},
	)

	output, err := uc.Execute(context.Background(), GetDocumentDownloadInput{
		DocumentID: "not-a-uuid",
	})

	require.Error(t, err)
	assert.Nil(t, output)
	assert.ErrorIs(t, err, ErrInvalidDocumentID)
}

func TestGetDocumentDownload_ShouldReturnErrDocumentNotFound(t *testing.T) {
	docID := uuid.New().String()
	uc := NewGetDocumentDownloadUseCase(
		&documentDomain.MockRepository{
			FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*documentDomain.Document, error) {
				return nil, nil
			},
		},
		&timelineDomain.MockRepository{},
		&mockDownloadURLProvider{},
	)

	output, err := uc.Execute(context.Background(), GetDocumentDownloadInput{
		DocumentID: docID,
	})

	require.Error(t, err)
	assert.Nil(t, output)
	assert.ErrorIs(t, err, ErrDocumentNotFound)
}

func TestGetDocumentDownload_ShouldCreateTimelineEvent(t *testing.T) {
	docID := uuid.New()
	now := time.Now()
	var capturedEvent *timelineDomain.TimelineEvent

	docs := &documentDomain.MockRepository{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*documentDomain.Document, error) {
			return &documentDomain.Document{
				ID: docID, CaseID: uuid.New(), Name: "Contrato",
				Type: documentDomain.DocumentTypeContract, FileName: "c.pdf",
				MimeType: "application/pdf", FileSize: 100, StorageKey: "k",
				CreatedAt: now, UpdatedAt: now,
			}, nil
		},
	}
	timeline := &timelineDomain.MockRepository{
		CreateFunc: func(ctx context.Context, event *timelineDomain.TimelineEvent) error {
			capturedEvent = event
			return nil
		},
	}
	provider := &mockDownloadURLProvider{
		GenerateFunc: func(document *documentDomain.Document) string {
			return "https://temporary-download.local/documents/" + docID.String()
		},
		ExpiresAtFunc: func() time.Time {
			return now.Add(30 * time.Minute)
		},
	}

	uc := NewGetDocumentDownloadUseCase(docs, timeline, provider)
	_, err := uc.Execute(context.Background(), GetDocumentDownloadInput{
		DocumentID: docID.String(),
	})

	require.NoError(t, err)
	require.NotNil(t, capturedEvent)
	assert.Equal(t, timelineDomain.EventDocumentDownloaded, capturedEvent.Type)
	assert.Contains(t, capturedEvent.Description, "Contrato")

	var meta timelineDomain.DocumentDownloadMetadata
	err = json.Unmarshal(capturedEvent.Metadata, &meta)
	require.NoError(t, err)
	assert.Equal(t, docID.String(), meta.DocumentID)
	assert.Equal(t, "c.pdf", meta.FileName)
}

func TestGetDocumentDownload_ShouldReturnErrorWhenTimelineFails(t *testing.T) {
	docID := uuid.New()
	now := time.Now()

	docs := &documentDomain.MockRepository{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*documentDomain.Document, error) {
			return &documentDomain.Document{
				ID: docID, CaseID: uuid.New(), Name: "Doc",
				Type: documentDomain.DocumentTypeContract, FileName: "d.pdf",
				MimeType: "application/pdf", FileSize: 100, StorageKey: "k",
				CreatedAt: now, UpdatedAt: now,
			}, nil
		},
	}
	timeline := &timelineDomain.MockRepository{
		CreateFunc: func(ctx context.Context, event *timelineDomain.TimelineEvent) error {
			return errors.New("timeline error")
		},
	}
	provider := &mockDownloadURLProvider{
		GenerateFunc: func(document *documentDomain.Document) string {
			return "url"
		},
		ExpiresAtFunc: func() time.Time {
			return now.Add(30 * time.Minute)
		},
	}

	uc := NewGetDocumentDownloadUseCase(docs, timeline, provider)
	output, err := uc.Execute(context.Background(), GetDocumentDownloadInput{
		DocumentID: docID.String(),
	})

	require.Error(t, err)
	assert.Nil(t, output)
	assert.ErrorIs(t, err, ErrTimelinePersistence)
}

func TestGetDocumentDownload_ShouldReturnErrorWhenRepoFails(t *testing.T) {
	docID := uuid.New()
	docs := &documentDomain.MockRepository{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*documentDomain.Document, error) {
			return nil, errors.New("db error")
		},
	}
	uc := NewGetDocumentDownloadUseCase(
		docs,
		&timelineDomain.MockRepository{},
		&mockDownloadURLProvider{},
	)

	output, err := uc.Execute(context.Background(), GetDocumentDownloadInput{
		DocumentID: docID.String(),
	})

	require.Error(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "db error")
}
