package timeline

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domain "legalflow/internal/domain/timeline"
	caseDomain "legalflow/internal/domain/legalcase"
)

func TestGetCaseTimeline_ShouldReturnOrderedEvents(t *testing.T) {
	caseID := uuid.New()
	now := time.Now()

	uc := NewGetCaseTimelineUseCase(
		&domain.MockRepository{
			FindByCaseIDFunc: func(ctx context.Context, id uuid.UUID) ([]domain.TimelineEvent, error) {
				return []domain.TimelineEvent{
					{ID: uuid.New(), CaseID: caseID, Type: domain.EventCaseCreated, Description: "Processo criado", CreatedAt: now},
					{ID: uuid.New(), CaseID: caseID, Type: domain.EventStatusChanged, Description: "Status alterado para active", CreatedAt: now.Add(time.Hour)},
				}, nil
			},
		},
		&caseDomain.MockCaseRepository{
			FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*caseDomain.Case, error) {
				return &caseDomain.Case{ID: caseID}, nil
			},
		},
	)

	output, err := uc.Execute(context.Background(), GetCaseTimelineInput{CaseID: caseID.String()})

	require.NoError(t, err)
	require.NotNil(t, output)
	assert.Len(t, output.Events, 2)
	assert.Equal(t, "CASE_CREATED", output.Events[0].Type)
	assert.Equal(t, "STATUS_CHANGED", output.Events[1].Type)
}

func TestGetCaseTimeline_ShouldReturnErrInvalidCaseID(t *testing.T) {
	uc := NewGetCaseTimelineUseCase(&domain.MockRepository{}, &caseDomain.MockCaseRepository{})
	output, err := uc.Execute(context.Background(), GetCaseTimelineInput{CaseID: "not-a-uuid"})

	require.Error(t, err)
	assert.Nil(t, output)
	assert.ErrorIs(t, err, ErrInvalidCaseID)
}

func TestGetCaseTimeline_ShouldReturnErrCaseNotFound(t *testing.T) {
	repo := &caseDomain.MockCaseRepository{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*caseDomain.Case, error) {
			return nil, nil
		},
	}

	uc := NewGetCaseTimelineUseCase(&domain.MockRepository{}, repo)
	output, err := uc.Execute(context.Background(), GetCaseTimelineInput{CaseID: uuid.New().String()})

	require.Error(t, err)
	assert.Nil(t, output)
	assert.ErrorIs(t, err, ErrCaseNotFound)
}

func TestGetCaseTimeline_ShouldReturnErrorWhenRepoFails(t *testing.T) {
	caseID := uuid.New()
	repo := &caseDomain.MockCaseRepository{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*caseDomain.Case, error) {
			return &caseDomain.Case{ID: caseID}, nil
		},
	}

	uc := NewGetCaseTimelineUseCase(
		&domain.MockRepository{
			FindByCaseIDFunc: func(ctx context.Context, id uuid.UUID) ([]domain.TimelineEvent, error) {
				return nil, errors.New("db error")
			},
		},
		repo,
	)

	output, err := uc.Execute(context.Background(), GetCaseTimelineInput{CaseID: caseID.String()})

	require.Error(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "db error")
}
