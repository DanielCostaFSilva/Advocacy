package hearing

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domain "legalflow/internal/domain/hearing"
)

type mockHearingRepo struct {
	ListFunc func(ctx context.Context, params domain.ListHearingsParams) ([]domain.Hearing, int64, error)
}

func (m *mockHearingRepo) Create(ctx context.Context, hearing *domain.Hearing) error {
	return nil
}

func (m *mockHearingRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Hearing, error) {
	return nil, nil
}

func (m *mockHearingRepo) ListByCaseID(ctx context.Context, caseID uuid.UUID) ([]domain.Hearing, error) {
	return nil, nil
}

func (m *mockHearingRepo) List(ctx context.Context, params domain.ListHearingsParams) ([]domain.Hearing, int64, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, params)
	}
	return nil, 0, nil
}

func TestListHearings_ShouldReturnPaginatedResult(t *testing.T) {
	now := time.Now().UTC()
	hearings := []domain.Hearing{
		{ID: uuid.New(), CaseID: uuid.New(), Title: "Audiência 1", Type: domain.HearingTypeConciliation, Location: "Local", ScheduledAt: now, CreatedAt: now},
		{ID: uuid.New(), CaseID: uuid.New(), Title: "Audiência 2", Type: domain.HearingTypeInstruction, Location: "Local", ScheduledAt: now, CreatedAt: now},
	}
	repo := &mockHearingRepo{
		ListFunc: func(ctx context.Context, params domain.ListHearingsParams) ([]domain.Hearing, int64, error) {
			return hearings, 2, nil
		},
	}

	uc := NewListHearingsUseCase(repo)
	output, err := uc.Execute(context.Background(), ListHearingsInput{
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

func TestListHearings_ShouldCapPageSizeAt100(t *testing.T) {
	repo := &mockHearingRepo{
		ListFunc: func(ctx context.Context, params domain.ListHearingsParams) ([]domain.Hearing, int64, error) {
			assert.Equal(t, 100, params.Limit)
			return nil, 0, nil
		},
	}

	uc := NewListHearingsUseCase(repo)
	output, err := uc.Execute(context.Background(), ListHearingsInput{
		Page:     1,
		PageSize: 200,
	})

	require.NoError(t, err)
	assert.Equal(t, 100, output.PageSize)
}

func TestListHearings_ShouldUseDefaultPageWhenInvalid(t *testing.T) {
	repo := &mockHearingRepo{
		ListFunc: func(ctx context.Context, params domain.ListHearingsParams) ([]domain.Hearing, int64, error) {
			assert.Equal(t, 0, params.Offset)
			return nil, 0, nil
		},
	}

	uc := NewListHearingsUseCase(repo)
	output, err := uc.Execute(context.Background(), ListHearingsInput{
		Page:     0,
		PageSize: 20,
	})

	require.NoError(t, err)
	assert.Equal(t, 1, output.Page)
}

func TestListHearings_ShouldFilterByType(t *testing.T) {
	repo := &mockHearingRepo{
		ListFunc: func(ctx context.Context, params domain.ListHearingsParams) ([]domain.Hearing, int64, error) {
			assert.Equal(t, "conciliation", params.Type)
			return nil, 0, nil
		},
	}

	uc := NewListHearingsUseCase(repo)
	_, err := uc.Execute(context.Background(), ListHearingsInput{
		Page: 1, PageSize: 20,
		Type: "conciliation",
	})
	require.NoError(t, err)
}

func TestListHearings_ShouldFilterByCaseID(t *testing.T) {
	caseID := uuid.New().String()
	repo := &mockHearingRepo{
		ListFunc: func(ctx context.Context, params domain.ListHearingsParams) ([]domain.Hearing, int64, error) {
			assert.Equal(t, caseID, params.CaseID)
			return nil, 0, nil
		},
	}

	uc := NewListHearingsUseCase(repo)
	_, err := uc.Execute(context.Background(), ListHearingsInput{
		Page: 1, PageSize: 20,
		CaseID: caseID,
	})
	require.NoError(t, err)
}

func TestListHearings_ShouldFilterByDateRange(t *testing.T) {
	repo := &mockHearingRepo{
		ListFunc: func(ctx context.Context, params domain.ListHearingsParams) ([]domain.Hearing, int64, error) {
			assert.NotNil(t, params.StartDate)
			assert.NotNil(t, params.EndDate)
			return nil, 0, nil
		},
	}

	uc := NewListHearingsUseCase(repo)
	_, err := uc.Execute(context.Background(), ListHearingsInput{
		Page:      1,
		PageSize:  20,
		StartDate: "2026-08-01",
		EndDate:   "2026-08-31",
	})
	require.NoError(t, err)
}

func TestListHearings_ShouldReturnErrorWhenRepoFails(t *testing.T) {
	repo := &mockHearingRepo{
		ListFunc: func(ctx context.Context, params domain.ListHearingsParams) ([]domain.Hearing, int64, error) {
			return nil, 0, errors.New("db error")
		},
	}

	uc := NewListHearingsUseCase(repo)
	output, err := uc.Execute(context.Background(), ListHearingsInput{
		Page: 1, PageSize: 20,
	})
	require.Error(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "db error")
}
