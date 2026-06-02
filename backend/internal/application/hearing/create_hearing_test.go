package hearing

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	caseDomain "legalflow/internal/domain/legalcase"
	hearingDomain "legalflow/internal/domain/hearing"
	timelineDomain "legalflow/internal/domain/timeline"
)

func TestCreateHearing_ShouldSucceed(t *testing.T) {
	caseID := uuid.New()
	future := time.Now().Add(48 * time.Hour)

	uc := NewCreateHearingUseCase(
		&hearingDomain.MockRepository{
			CreateFunc: func(ctx context.Context, h *hearingDomain.Hearing) error {
				h.ID = uuid.New()
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

	output, err := uc.Execute(context.Background(), CreateHearingInput{
		CaseID:      caseID.String(),
		Title:       "Audiência de Conciliação",
		Description: "Descrição",
		Type:        "conciliation",
		Location:    "Fórum Central",
		ScheduledAt: future,
	})

	require.NoError(t, err)
	require.NotNil(t, output)
	assert.NotEmpty(t, output.ID)
	assert.Equal(t, caseID.String(), output.CaseID)
	assert.Equal(t, "Audiência de Conciliação", output.Title)
	assert.Equal(t, "conciliation", output.Type)
	assert.Equal(t, "Fórum Central", output.Location)
	assert.Equal(t, future, output.ScheduledAt)
}

func TestCreateHearing_ShouldReturnErrInvalidCaseID(t *testing.T) {
	uc := NewCreateHearingUseCase(
		&hearingDomain.MockRepository{},
		&caseDomain.MockCaseRepository{},
		&timelineDomain.MockRepository{},
	)

	output, err := uc.Execute(context.Background(), CreateHearingInput{
		CaseID: "not-a-uuid",
	})

	require.Error(t, err)
	assert.Nil(t, output)
	assert.ErrorIs(t, err, ErrInvalidCaseID)
}

func TestCreateHearing_ShouldReturnErrCaseNotFound(t *testing.T) {
	caseID := uuid.New().String()

	uc := NewCreateHearingUseCase(
		&hearingDomain.MockRepository{},
		&caseDomain.MockCaseRepository{
			FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*caseDomain.Case, error) {
				return nil, nil
			},
		},
		&timelineDomain.MockRepository{},
	)

	output, err := uc.Execute(context.Background(), CreateHearingInput{
		CaseID: caseID,
	})

	require.Error(t, err)
	assert.Nil(t, output)
	assert.ErrorIs(t, err, ErrCaseNotFound)
}

func TestCreateHearing_ShouldReturnErrorOnInvalidHearingData(t *testing.T) {
	caseID := uuid.New()

	uc := NewCreateHearingUseCase(
		&hearingDomain.MockRepository{},
		&caseDomain.MockCaseRepository{
			FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*caseDomain.Case, error) {
				return &caseDomain.Case{ID: caseID}, nil
			},
		},
		&timelineDomain.MockRepository{},
	)

	output, err := uc.Execute(context.Background(), CreateHearingInput{
		CaseID:      caseID.String(),
		Title:       "AB",
		Type:        "conciliation",
		Location:    "Local",
		ScheduledAt: time.Now().Add(48 * time.Hour),
	})

	require.Error(t, err)
	assert.Nil(t, output)
}

func TestCreateHearing_ShouldReturnErrorWhenHearingRepoFails(t *testing.T) {
	caseID := uuid.New()

	uc := NewCreateHearingUseCase(
		&hearingDomain.MockRepository{
			CreateFunc: func(ctx context.Context, h *hearingDomain.Hearing) error {
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

	output, err := uc.Execute(context.Background(), CreateHearingInput{
		CaseID:      caseID.String(),
		Title:       "Audiência Válida",
		Type:        "conciliation",
		Location:    "Local",
		ScheduledAt: time.Now().Add(48 * time.Hour),
	})

	require.Error(t, err)
	assert.Nil(t, output)
}

func TestCreateHearing_ShouldReturnErrorWhenTimelineRepoFails(t *testing.T) {
	caseID := uuid.New()

	uc := NewCreateHearingUseCase(
		&hearingDomain.MockRepository{
			CreateFunc: func(ctx context.Context, h *hearingDomain.Hearing) error {
				h.ID = uuid.New()
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

	output, err := uc.Execute(context.Background(), CreateHearingInput{
		CaseID:      caseID.String(),
		Title:       "Audiência Válida",
		Type:        "conciliation",
		Location:    "Local",
		ScheduledAt: time.Now().Add(48 * time.Hour),
	})

	require.Error(t, err)
	assert.Nil(t, output)
}

func TestCreateHearing_ShouldCreateTimelineEvent(t *testing.T) {
	caseID := uuid.New()
	future := time.Now().UTC().Add(48 * time.Hour)
	var capturedEvent *timelineDomain.TimelineEvent

	uc := NewCreateHearingUseCase(
		&hearingDomain.MockRepository{
			CreateFunc: func(ctx context.Context, h *hearingDomain.Hearing) error {
				h.ID = uuid.New()
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

	_, err := uc.Execute(context.Background(), CreateHearingInput{
		CaseID:      caseID.String(),
		Title:       "Audiência de Conciliação",
		Description: "Desc",
		Type:        "conciliation",
		Location:    "Fórum",
		ScheduledAt: future,
	})

	require.NoError(t, err)
	require.NotNil(t, capturedEvent)
	assert.Equal(t, timelineDomain.EventHearingCreated, capturedEvent.Type)
	assert.Equal(t, caseID, capturedEvent.CaseID)
	assert.Contains(t, capturedEvent.Description, "Audiência de Conciliação")
	require.NotEmpty(t, capturedEvent.Metadata)

	var meta timelineDomain.HearingMetadata
	err = json.Unmarshal(capturedEvent.Metadata, &meta)
	require.NoError(t, err)
	assert.NotEmpty(t, meta.HearingID)
	assert.Equal(t, "conciliation", meta.Type)
	assert.True(t, future.Equal(meta.ScheduledAt))
}
