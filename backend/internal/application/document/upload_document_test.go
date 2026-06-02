package document

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	caseDomain "legalflow/internal/domain/legalcase"
	documentDomain "legalflow/internal/domain/document"
	timelineDomain "legalflow/internal/domain/timeline"
)

func TestUploadDocument_ShouldSucceed(t *testing.T) {
	caseID := uuid.New()

	uc := NewUploadDocumentUseCase(
		&documentDomain.MockRepository{
			CreateFunc: func(ctx context.Context, doc *documentDomain.Document) error {
				doc.ID = uuid.New()
				return nil
			},
		},
		&caseDomain.MockCaseRepository{
			FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*caseDomain.Case, error) {
				return &caseDomain.Case{ID: caseID}, nil
			},
		},
		&timelineDomain.MockRepository{
			CreateFunc: func(ctx context.Context, event *timelineDomain.TimelineEvent) error {
				return nil
			},
		},
	)

	output, err := uc.Execute(context.Background(), UploadDocumentInput{
		CaseID:     caseID.String(),
		Name:       "Contrato de Honorários",
		Type:       "contract",
		FileName:   "contrato.pdf",
		MimeType:   "application/pdf",
		FileSize:   2048,
		StorageKey: "cases/" + caseID.String() + "/contrato.pdf",
	})

	require.NoError(t, err)
	require.NotNil(t, output)
	assert.NotEmpty(t, output.ID)
	assert.Equal(t, caseID.String(), output.CaseID)
	assert.Equal(t, "Contrato de Honorários", output.Name)
	assert.Equal(t, "contract", output.Type)
	assert.Equal(t, "contrato.pdf", output.FileName)
	assert.Equal(t, int64(2048), output.FileSize)
}

func TestUploadDocument_ShouldReturnErrInvalidCaseID(t *testing.T) {
	uc := NewUploadDocumentUseCase(
		&documentDomain.MockRepository{},
		&caseDomain.MockCaseRepository{},
		&timelineDomain.MockRepository{},
	)

	output, err := uc.Execute(context.Background(), UploadDocumentInput{
		CaseID: "not-a-uuid",
	})

	require.Error(t, err)
	assert.Nil(t, output)
	assert.ErrorIs(t, err, ErrInvalidCaseID)
}

func TestUploadDocument_ShouldReturnErrCaseNotFound(t *testing.T) {
	caseID := uuid.New().String()

	uc := NewUploadDocumentUseCase(
		&documentDomain.MockRepository{},
		&caseDomain.MockCaseRepository{
			FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*caseDomain.Case, error) {
				return nil, nil
			},
		},
		&timelineDomain.MockRepository{},
	)

	output, err := uc.Execute(context.Background(), UploadDocumentInput{
		CaseID: caseID,
	})

	require.Error(t, err)
	assert.Nil(t, output)
	assert.ErrorIs(t, err, ErrCaseNotFound)
}

func TestUploadDocument_ShouldReturnErrInvalidInputWhenInvalidData(t *testing.T) {
	caseID := uuid.New()

	uc := NewUploadDocumentUseCase(
		&documentDomain.MockRepository{},
		&caseDomain.MockCaseRepository{
			FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*caseDomain.Case, error) {
				return &caseDomain.Case{ID: caseID}, nil
			},
		},
		&timelineDomain.MockRepository{},
	)

	output, err := uc.Execute(context.Background(), UploadDocumentInput{
		CaseID:   caseID.String(),
		Name:     "AB",
		Type:     "contract",
		FileName: "file.pdf",
		MimeType: "application/pdf",
		FileSize: 1024,
	})

	require.Error(t, err)
	assert.Nil(t, output)
	assert.ErrorIs(t, err, ErrInvalidInput)
}

func TestUploadDocument_ShouldReturnErrorWhenDocRepoFails(t *testing.T) {
	caseID := uuid.New()

	uc := NewUploadDocumentUseCase(
		&documentDomain.MockRepository{
			CreateFunc: func(ctx context.Context, doc *documentDomain.Document) error {
				return errors.New("db error")
			},
		},
		&caseDomain.MockCaseRepository{
			FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*caseDomain.Case, error) {
				return &caseDomain.Case{ID: caseID}, nil
			},
		},
		&timelineDomain.MockRepository{},
	)

	output, err := uc.Execute(context.Background(), UploadDocumentInput{
		CaseID:     caseID.String(),
		Name:       "Contrato Válido",
		Type:       "contract",
		FileName:   "file.pdf",
		MimeType:   "application/pdf",
		FileSize:   1024,
		StorageKey: "key",
	})

	require.Error(t, err)
	assert.Nil(t, output)
}

func TestUploadDocument_ShouldReturnErrorWhenTimelineRepoFails(t *testing.T) {
	caseID := uuid.New()

	uc := NewUploadDocumentUseCase(
		&documentDomain.MockRepository{
			CreateFunc: func(ctx context.Context, doc *documentDomain.Document) error {
				doc.ID = uuid.New()
				return nil
			},
		},
		&caseDomain.MockCaseRepository{
			FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*caseDomain.Case, error) {
				return &caseDomain.Case{ID: caseID}, nil
			},
		},
		&timelineDomain.MockRepository{
			CreateFunc: func(ctx context.Context, event *timelineDomain.TimelineEvent) error {
				return errors.New("timeline error")
			},
		},
	)

	output, err := uc.Execute(context.Background(), UploadDocumentInput{
		CaseID:     caseID.String(),
		Name:       "Contrato Válido",
		Type:       "contract",
		FileName:   "file.pdf",
		MimeType:   "application/pdf",
		FileSize:   1024,
		StorageKey: "key",
	})

	require.Error(t, err)
	assert.Nil(t, output)
}

func TestUploadDocument_ShouldCreateTimelineEvent(t *testing.T) {
	caseID := uuid.New()
	var capturedEvent *timelineDomain.TimelineEvent

	uc := NewUploadDocumentUseCase(
		&documentDomain.MockRepository{
			CreateFunc: func(ctx context.Context, doc *documentDomain.Document) error {
				doc.ID = uuid.New()
				return nil
			},
		},
		&caseDomain.MockCaseRepository{
			FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*caseDomain.Case, error) {
				return &caseDomain.Case{ID: caseID}, nil
			},
		},
		&timelineDomain.MockRepository{
			CreateFunc: func(ctx context.Context, event *timelineDomain.TimelineEvent) error {
				capturedEvent = event
				return nil
			},
		},
	)

	_, err := uc.Execute(context.Background(), UploadDocumentInput{
		CaseID:     caseID.String(),
		Name:       "Contrato de Honorários",
		Type:       "contract",
		FileName:   "contrato.pdf",
		MimeType:   "application/pdf",
		FileSize:   2048,
		StorageKey: "key",
	})

	require.NoError(t, err)
	require.NotNil(t, capturedEvent)
	assert.Equal(t, timelineDomain.EventDocumentUploaded, capturedEvent.Type)
	assert.Equal(t, caseID, capturedEvent.CaseID)
	assert.Contains(t, capturedEvent.Description, "Contrato de Honorários")
	require.NotEmpty(t, capturedEvent.Metadata)

	var meta timelineDomain.DocumentMetadata
	err = json.Unmarshal(capturedEvent.Metadata, &meta)
	require.NoError(t, err)
	assert.NotEmpty(t, meta.DocumentID)
	assert.Equal(t, "contract", meta.DocumentType)
	assert.Equal(t, "contrato.pdf", meta.FileName)
}
