package document

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domain "legalflow/internal/domain/document"
	caseDomain "legalflow/internal/domain/legalcase"
)

func TestListCaseDocuments_ShouldReturnPaginatedResult(t *testing.T) {
	now := time.Now().UTC()
	caseID := uuid.New()
	docs := []domain.Document{
		{ID: uuid.New(), CaseID: caseID, Name: "Doc 1", Type: domain.DocumentTypeContract, FileName: "f1.pdf", MimeType: "application/pdf", FileSize: 100, StorageKey: "k1", CreatedAt: now},
		{ID: uuid.New(), CaseID: caseID, Name: "Doc 2", Type: domain.DocumentTypeEvidence, FileName: "f2.pdf", MimeType: "application/pdf", FileSize: 200, StorageKey: "k2", CreatedAt: now},
	}

	uc := NewListCaseDocumentsUseCase(
		&domain.MockRepository{
			ListByCaseIDPaginatedFunc: func(ctx context.Context, params domain.ListDocumentsParams) ([]domain.Document, int64, error) {
				return docs, 2, nil
			},
		},
		&caseDomain.MockCaseRepository{
			FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*caseDomain.Case, error) {
				return &caseDomain.Case{ID: caseID}, nil
			},
		},
	)

	output, err := uc.Execute(context.Background(), ListCaseDocumentsInput{
		CaseID:   caseID.String(),
		Page:     1,
		PageSize: 20,
	})

	require.NoError(t, err)
	require.NotNil(t, output)
	assert.Len(t, output.Data, 2)
	assert.Equal(t, 1, output.Page)
	assert.Equal(t, 20, output.PageSize)
	assert.Equal(t, int64(2), output.Total)
	assert.Equal(t, 1, output.TotalPages)
}

func TestListCaseDocuments_ShouldCapPageSizeAt100(t *testing.T) {
	caseID := uuid.New()
	uc := NewListCaseDocumentsUseCase(
		&domain.MockRepository{
			ListByCaseIDPaginatedFunc: func(ctx context.Context, params domain.ListDocumentsParams) ([]domain.Document, int64, error) {
				assert.Equal(t, 100, params.Limit)
				return nil, 0, nil
			},
		},
		&caseDomain.MockCaseRepository{
			FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*caseDomain.Case, error) {
				return &caseDomain.Case{ID: caseID}, nil
			},
		},
	)

	output, err := uc.Execute(context.Background(), ListCaseDocumentsInput{
		CaseID:   caseID.String(),
		Page:     1,
		PageSize: 200,
	})

	require.NoError(t, err)
	assert.Equal(t, 100, output.PageSize)
}

func TestListCaseDocuments_ShouldUseDefaultPageWhenInvalid(t *testing.T) {
	caseID := uuid.New()
	uc := NewListCaseDocumentsUseCase(
		&domain.MockRepository{
			ListByCaseIDPaginatedFunc: func(ctx context.Context, params domain.ListDocumentsParams) ([]domain.Document, int64, error) {
				assert.Equal(t, 0, params.Offset)
				return nil, 0, nil
			},
		},
		&caseDomain.MockCaseRepository{
			FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*caseDomain.Case, error) {
				return &caseDomain.Case{ID: caseID}, nil
			},
		},
	)

	output, err := uc.Execute(context.Background(), ListCaseDocumentsInput{
		CaseID:   caseID.String(),
		Page:     0,
		PageSize: 20,
	})

	require.NoError(t, err)
	assert.Equal(t, 1, output.Page)
}

func TestListCaseDocuments_ShouldFilterByType(t *testing.T) {
	caseID := uuid.New()
	uc := NewListCaseDocumentsUseCase(
		&domain.MockRepository{
			ListByCaseIDPaginatedFunc: func(ctx context.Context, params domain.ListDocumentsParams) ([]domain.Document, int64, error) {
				assert.Equal(t, "contract", params.Type)
				return nil, 0, nil
			},
		},
		&caseDomain.MockCaseRepository{
			FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*caseDomain.Case, error) {
				return &caseDomain.Case{ID: caseID}, nil
			},
		},
	)

	_, err := uc.Execute(context.Background(), ListCaseDocumentsInput{
		CaseID:   caseID.String(),
		Page:     1,
		PageSize: 20,
		Type:     "contract",
	})
	require.NoError(t, err)
}

func TestListCaseDocuments_ShouldReturnErrorWhenRepoFails(t *testing.T) {
	caseID := uuid.New()
	uc := NewListCaseDocumentsUseCase(
		&domain.MockRepository{
			ListByCaseIDPaginatedFunc: func(ctx context.Context, params domain.ListDocumentsParams) ([]domain.Document, int64, error) {
				return nil, 0, errors.New("db error")
			},
		},
		&caseDomain.MockCaseRepository{
			FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*caseDomain.Case, error) {
				return &caseDomain.Case{ID: caseID}, nil
			},
		},
	)

	output, err := uc.Execute(context.Background(), ListCaseDocumentsInput{
		CaseID:   caseID.String(),
		Page:     1,
		PageSize: 20,
	})
	require.Error(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "db error")
}

func TestListCaseDocuments_ShouldReturnErrCaseNotFound(t *testing.T) {
	caseID := uuid.New()
	uc := NewListCaseDocumentsUseCase(
		&domain.MockRepository{},
		&caseDomain.MockCaseRepository{
			FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*caseDomain.Case, error) {
				return nil, nil
			},
		},
	)

	output, err := uc.Execute(context.Background(), ListCaseDocumentsInput{
		CaseID:   caseID.String(),
		Page:     1,
		PageSize: 20,
	})
	require.Error(t, err)
	assert.Nil(t, output)
	assert.ErrorIs(t, err, ErrCaseNotFound)
}

func TestListCaseDocuments_ShouldReturnErrInvalidCaseID(t *testing.T) {
	uc := NewListCaseDocumentsUseCase(
		&domain.MockRepository{},
		&caseDomain.MockCaseRepository{},
	)

	output, err := uc.Execute(context.Background(), ListCaseDocumentsInput{
		CaseID: "invalid-uuid",
	})
	require.Error(t, err)
	assert.Nil(t, output)
	assert.ErrorIs(t, err, ErrInvalidCaseID)
}
