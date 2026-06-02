package legalcase

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domain "legalflow/internal/domain/legalcase"
)

func TestListCases_ShouldReturnPaginatedResult(t *testing.T) {
	cases := []domain.Case{
		{ID: uuid.New(), ClientID: uuid.New(), Number: "ABC-123", Title: "Case 1", Court: "TJSP", Status: domain.CaseStatusDraft},
		{ID: uuid.New(), ClientID: uuid.New(), Number: "DEF-456", Title: "Case 2", Court: "TJPE", Status: domain.CaseStatusActive},
	}
	repo := &domain.MockCaseRepository{
		ListFunc: func(ctx context.Context, params domain.ListCasesParams) ([]domain.Case, int64, error) {
			return cases, 2, nil
		},
	}

	uc := NewListCasesUseCase(repo)
	output, err := uc.Execute(context.Background(), ListCasesInput{
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

func TestListCases_ShouldCapPageSizeAt100(t *testing.T) {
	repo := &domain.MockCaseRepository{
		ListFunc: func(ctx context.Context, params domain.ListCasesParams) ([]domain.Case, int64, error) {
			assert.Equal(t, 100, params.Limit)
			return nil, 0, nil
		},
	}

	uc := NewListCasesUseCase(repo)
	output, err := uc.Execute(context.Background(), ListCasesInput{
		Page:     1,
		PageSize: 200,
	})

	require.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, 100, output.PageSize)
}

func TestListCases_ShouldUseDefaultPageWhenInvalid(t *testing.T) {
	repo := &domain.MockCaseRepository{
		ListFunc: func(ctx context.Context, params domain.ListCasesParams) ([]domain.Case, int64, error) {
			assert.Equal(t, 0, params.Offset)
			return nil, 0, nil
		},
	}

	uc := NewListCasesUseCase(repo)
	output, err := uc.Execute(context.Background(), ListCasesInput{
		Page:     0,
		PageSize: 20,
	})

	require.NoError(t, err)
	assert.Equal(t, 1, output.Page)
}

func TestListCases_ShouldFilterByStatus(t *testing.T) {
	repo := &domain.MockCaseRepository{
		ListFunc: func(ctx context.Context, params domain.ListCasesParams) ([]domain.Case, int64, error) {
			assert.Equal(t, "active", params.Status)
			return nil, 0, nil
		},
	}

	uc := NewListCasesUseCase(repo)
	_, err := uc.Execute(context.Background(), ListCasesInput{
		Page:     1,
		PageSize: 20,
		Status:   "active",
	})

	require.NoError(t, err)
}

func TestListCases_ShouldFilterByClientID(t *testing.T) {
	clientID := uuid.New().String()
	repo := &domain.MockCaseRepository{
		ListFunc: func(ctx context.Context, params domain.ListCasesParams) ([]domain.Case, int64, error) {
			assert.Equal(t, clientID, params.ClientID)
			return nil, 0, nil
		},
	}

	uc := NewListCasesUseCase(repo)
	_, err := uc.Execute(context.Background(), ListCasesInput{
		Page:     1,
		PageSize: 20,
		ClientID: clientID,
	})

	require.NoError(t, err)
}

func TestListCases_ShouldReturnErrorWhenRepoFails(t *testing.T) {
	repo := &domain.MockCaseRepository{
		ListFunc: func(ctx context.Context, params domain.ListCasesParams) ([]domain.Case, int64, error) {
			return nil, 0, errors.New("db error")
		},
	}

	uc := NewListCasesUseCase(repo)
	output, err := uc.Execute(context.Background(), ListCasesInput{
		Page:     1,
		PageSize: 20,
	})

	require.Error(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "db error")
}
