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

func TestUpdateCaseStatus_ShouldUpdateSuccessfully(t *testing.T) {
	caseID := uuid.New()
	repo := &domain.MockCaseRepository{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Case, error) {
			return &domain.Case{
				ID:     caseID,
				Status: domain.CaseStatusDraft,
			}, nil
		},
		UpdateStatusFunc: func(ctx context.Context, id uuid.UUID, status domain.CaseStatus) error {
			assert.Equal(t, caseID, id)
			assert.Equal(t, domain.CaseStatusActive, status)
			return nil
		},
	}

	uc := NewUpdateCaseStatusUseCase(repo, &mockTimelineRepo{})
	output, err := uc.Execute(context.Background(), UpdateCaseStatusInput{
		CaseID: caseID.String(),
		Status: "active",
	})

	require.NoError(t, err)
	require.NotNil(t, output)
	assert.Equal(t, caseID.String(), output.ID)
	assert.Equal(t, "active", output.Status)
}

func TestUpdateCaseStatus_ShouldReturnErrInvalidCaseID(t *testing.T) {
	uc := NewUpdateCaseStatusUseCase(&domain.MockCaseRepository{}, &mockTimelineRepo{})
	output, err := uc.Execute(context.Background(), UpdateCaseStatusInput{
		CaseID: "not-a-uuid",
		Status: "active",
	})

	require.Error(t, err)
	assert.Nil(t, output)
	assert.ErrorIs(t, err, ErrInvalidCaseID)
}

func TestUpdateCaseStatus_ShouldReturnErrInvalidStatus(t *testing.T) {
	uc := NewUpdateCaseStatusUseCase(&domain.MockCaseRepository{}, &mockTimelineRepo{})
	output, err := uc.Execute(context.Background(), UpdateCaseStatusInput{
		CaseID: uuid.New().String(),
		Status: "invalid",
	})

	require.Error(t, err)
	assert.Nil(t, output)
	assert.ErrorIs(t, err, ErrInvalidStatus)
}

func TestUpdateCaseStatus_ShouldReturnErrCaseNotFound(t *testing.T) {
	repo := &domain.MockCaseRepository{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Case, error) {
			return nil, nil
		},
	}

	uc := NewUpdateCaseStatusUseCase(repo, &mockTimelineRepo{})
	output, err := uc.Execute(context.Background(), UpdateCaseStatusInput{
		CaseID: uuid.New().String(),
		Status: "active",
	})

	require.Error(t, err)
	assert.Nil(t, output)
	assert.ErrorIs(t, err, ErrCaseNotFound)
}

func TestUpdateCaseStatus_ShouldReturnErrInvalidStatusTransition(t *testing.T) {
	caseID := uuid.New()
	repo := &domain.MockCaseRepository{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Case, error) {
			return &domain.Case{
				ID:     caseID,
				Status: domain.CaseStatusDraft,
			}, nil
		},
	}

	uc := NewUpdateCaseStatusUseCase(repo, &mockTimelineRepo{})
	output, err := uc.Execute(context.Background(), UpdateCaseStatusInput{
		CaseID: caseID.String(),
	Status: "suspended",
	})

	require.Error(t, err)
	assert.Nil(t, output)
	assert.ErrorIs(t, err, ErrInvalidStatusTransition)
}

func TestUpdateCaseStatus_ShouldReturnErrorWhenRepoFails(t *testing.T) {
	caseID := uuid.New()
	repo := &domain.MockCaseRepository{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Case, error) {
			return &domain.Case{
				ID:     caseID,
				Status: domain.CaseStatusDraft,
			}, nil
		},
		UpdateStatusFunc: func(ctx context.Context, id uuid.UUID, status domain.CaseStatus) error {
			return errors.New("db error")
		},
	}

	uc := NewUpdateCaseStatusUseCase(repo, &mockTimelineRepo{})
	output, err := uc.Execute(context.Background(), UpdateCaseStatusInput{
		CaseID: caseID.String(),
		Status: "active",
	})

	require.Error(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "db error")
}
