package legalcase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domain "legalflow/internal/domain/legalcase"
)

func TestGetCaseByID_ShouldReturnCase(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	caseID := uuid.New()
	clientID := uuid.New()

	repo := &domain.MockCaseRepository{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Case, error) {
			return &domain.Case{
				ID:          caseID,
				ClientID:    clientID,
				Number:      "ABC-123",
				Title:       "Test Case",
				Description: "Description",
				Court:       "TJSP",
				Status:      domain.CaseStatusActive,
				CreatedAt:   now,
				UpdatedAt:   now,
			}, nil
		},
	}

	uc := NewGetCaseByIDUseCase(repo)
	output, err := uc.Execute(context.Background(), GetCaseByIDInput{ID: caseID.String()})

	require.NoError(t, err)
	require.NotNil(t, output)
	assert.Equal(t, caseID.String(), output.ID)
	assert.Equal(t, clientID.String(), output.ClientID)
	assert.Equal(t, "ABC-123", output.Number)
	assert.Equal(t, "Test Case", output.Title)
	assert.Equal(t, "Description", output.Description)
	assert.Equal(t, "TJSP", output.Court)
	assert.Equal(t, "active", output.Status)
	assert.Equal(t, now, output.CreatedAt)
	assert.Equal(t, now, output.UpdatedAt)
}

func TestGetCaseByID_ShouldReturnErrInvalidCaseID(t *testing.T) {
	uc := NewGetCaseByIDUseCase(&domain.MockCaseRepository{})
	output, err := uc.Execute(context.Background(), GetCaseByIDInput{ID: "not-a-uuid"})

	require.Error(t, err)
	assert.Nil(t, output)
	assert.ErrorIs(t, err, ErrInvalidCaseID)
}

func TestGetCaseByID_ShouldReturnErrCaseNotFound(t *testing.T) {
	repo := &domain.MockCaseRepository{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Case, error) {
			return nil, nil
		},
	}

	uc := NewGetCaseByIDUseCase(repo)
	output, err := uc.Execute(context.Background(), GetCaseByIDInput{ID: uuid.New().String()})

	require.Error(t, err)
	assert.Nil(t, output)
	assert.ErrorIs(t, err, ErrCaseNotFound)
}

func TestGetCaseByID_ShouldReturnErrorWhenRepoFails(t *testing.T) {
	repo := &domain.MockCaseRepository{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Case, error) {
			return nil, errors.New("db error")
		},
	}

	uc := NewGetCaseByIDUseCase(repo)
	output, err := uc.Execute(context.Background(), GetCaseByIDInput{ID: uuid.New().String()})

	require.Error(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "db error")
}
