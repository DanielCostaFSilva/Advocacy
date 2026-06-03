package contract

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	contractDomain "legalflow/internal/domain/contract"
	timelineDomain "legalflow/internal/domain/timeline"
)

type mockUpdateContractRepo struct {
	FindByIDFunc func(ctx context.Context, id uuid.UUID) (*contractDomain.Contract, error)
	UpdateFunc   func(ctx context.Context, contract *contractDomain.Contract) error
}

func (m *mockUpdateContractRepo) Create(ctx context.Context, contract *contractDomain.Contract) error {
	return nil
}

func (m *mockUpdateContractRepo) FindByID(ctx context.Context, id uuid.UUID) (*contractDomain.Contract, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockUpdateContractRepo) ListByCaseID(ctx context.Context, caseID uuid.UUID) ([]contractDomain.Contract, error) {
	return nil, nil
}

func (m *mockUpdateContractRepo) List(ctx context.Context, params contractDomain.ListContractsParams) ([]contractDomain.Contract, int64, error) {
	return nil, 0, nil
}

func (m *mockUpdateContractRepo) Update(ctx context.Context, contract *contractDomain.Contract) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, contract)
	}
	return nil
}

type mockUpdateTimelineRepo struct {
	CreateFunc func(ctx context.Context, event *timelineDomain.TimelineEvent) error
}

func (m *mockUpdateTimelineRepo) Create(ctx context.Context, event *timelineDomain.TimelineEvent) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, event)
	}
	return nil
}

func TestUpdateContract_ShouldUpdateSuccessfully(t *testing.T) {
	now := time.Now().UTC()
	contractID := uuid.New()
	clientID := uuid.New()
	caseID := uuid.New()

	contract := &contractDomain.Contract{
		ID:        contractID,
		ClientID:  clientID,
		CaseID:    caseID,
		Title:     "Contrato Original",
		Type:      contractDomain.ContractTypeFixedFee,
		Amount:    decimal.NewFromFloat(10000),
		StartDate: now,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	repo := &mockUpdateContractRepo{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*contractDomain.Contract, error) {
			return contract, nil
		},
		UpdateFunc: func(ctx context.Context, c *contractDomain.Contract) error {
			c.UpdatedAt = time.Now()
			return nil
		},
	}
	timelineRepo := &mockUpdateTimelineRepo{
		CreateFunc: func(ctx context.Context, event *timelineDomain.TimelineEvent) error {
			return nil
		},
	}

	uc := NewUpdateContractUseCase(repo, timelineRepo)
	output, err := uc.Execute(context.Background(), UpdateContractInput{
		ID:        contractID.String(),
		Title:     "Contrato Atualizado",
		Type:      "fixed_fee",
		Amount:    "15000.00",
		StartDate: now.Format(time.RFC3339),
		Active:    true,
	})

	require.NoError(t, err)
	require.NotNil(t, output)
	assert.Equal(t, contractID.String(), output.ID)
	assert.Equal(t, "Contrato Atualizado", output.Title)
	assert.Equal(t, "15000", output.Amount)
	assert.True(t, output.Active)
	assert.False(t, output.UpdatedAt.IsZero())
}

func TestUpdateContract_ShouldReturnErrInvalidContractID(t *testing.T) {
	repo := &mockUpdateContractRepo{}
	timelineRepo := &mockUpdateTimelineRepo{}
	uc := NewUpdateContractUseCase(repo, timelineRepo)

	output, err := uc.Execute(context.Background(), UpdateContractInput{
		ID: "invalid-uuid",
	})

	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidContractID))
	assert.Nil(t, output)
}

func TestUpdateContract_ShouldReturnErrContractNotFound(t *testing.T) {
	repo := &mockUpdateContractRepo{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*contractDomain.Contract, error) {
			return nil, nil
		},
	}
	timelineRepo := &mockUpdateTimelineRepo{}
	uc := NewUpdateContractUseCase(repo, timelineRepo)

	output, err := uc.Execute(context.Background(), UpdateContractInput{
		ID:        uuid.New().String(),
		Title:     "Teste",
		Type:      "fixed_fee",
		Amount:    "1000.00",
		StartDate: time.Now().UTC().Format(time.RFC3339),
		Active:    true,
	})

	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrContractNotFound))
	assert.Nil(t, output)
}

func TestUpdateContract_ShouldReturnValidationError(t *testing.T) {
	now := time.Now().UTC()
	contract := &contractDomain.Contract{
		ID:        uuid.New(),
		ClientID:  uuid.New(),
		CaseID:    uuid.New(),
		Title:     "Original",
		Type:      contractDomain.ContractTypeFixedFee,
		Amount:    decimal.NewFromFloat(10000),
		StartDate: now,
		Active:    true,
	}

	repo := &mockUpdateContractRepo{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*contractDomain.Contract, error) {
			return contract, nil
		},
	}
	timelineRepo := &mockUpdateTimelineRepo{}
	uc := NewUpdateContractUseCase(repo, timelineRepo)

	output, err := uc.Execute(context.Background(), UpdateContractInput{
		ID:        contract.ID.String(),
		Title:     "AB",
		Type:      "fixed_fee",
		Amount:    "1000.00",
		StartDate: now.Format(time.RFC3339),
		Active:    true,
	})

	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidInput))
	assert.Nil(t, output)
}

func TestUpdateContract_ShouldReturnErrorWhenRepoFails(t *testing.T) {
	now := time.Now().UTC()
	contract := &contractDomain.Contract{
		ID:        uuid.New(),
		ClientID:  uuid.New(),
		CaseID:    uuid.New(),
		Title:     "Original",
		Type:      contractDomain.ContractTypeFixedFee,
		Amount:    decimal.NewFromFloat(10000),
		StartDate: now,
		Active:    true,
	}

	repo := &mockUpdateContractRepo{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*contractDomain.Contract, error) {
			return contract, nil
		},
		UpdateFunc: func(ctx context.Context, c *contractDomain.Contract) error {
			return errors.New("db error")
		},
	}
	timelineRepo := &mockUpdateTimelineRepo{}
	uc := NewUpdateContractUseCase(repo, timelineRepo)

	output, err := uc.Execute(context.Background(), UpdateContractInput{
		ID:        contract.ID.String(),
		Title:     "Updated Title",
		Type:      "fixed_fee",
		Amount:    "1000.00",
		StartDate: now.Format(time.RFC3339),
		Active:    true,
	})

	require.Error(t, err)
	assert.Nil(t, output)
}

func TestUpdateContract_ShouldReturnErrorWhenTimelineFails(t *testing.T) {
	now := time.Now().UTC()
	contract := &contractDomain.Contract{
		ID:        uuid.New(),
		ClientID:  uuid.New(),
		CaseID:    uuid.New(),
		Title:     "Original",
		Type:      contractDomain.ContractTypeFixedFee,
		Amount:    decimal.NewFromFloat(10000),
		StartDate: now,
		Active:    true,
	}

	repo := &mockUpdateContractRepo{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*contractDomain.Contract, error) {
			return contract, nil
		},
		UpdateFunc: func(ctx context.Context, c *contractDomain.Contract) error {
			c.UpdatedAt = time.Now()
			return nil
		},
	}
	timelineRepo := &mockUpdateTimelineRepo{
		CreateFunc: func(ctx context.Context, event *timelineDomain.TimelineEvent) error {
			return errors.New("timeline error")
		},
	}
	uc := NewUpdateContractUseCase(repo, timelineRepo)

	output, err := uc.Execute(context.Background(), UpdateContractInput{
		ID:        contract.ID.String(),
		Title:     "Updated Title",
		Type:      "fixed_fee",
		Amount:    "1000.00",
		StartDate: now.Format(time.RFC3339),
		Active:    true,
	})

	require.Error(t, err)
	assert.Nil(t, output)
}

func TestUpdateContract_ShouldGenerateCONTRACTUPDATEDEvent(t *testing.T) {
	now := time.Now().UTC()
	contractID := uuid.New()

	contract := &contractDomain.Contract{
		ID:        contractID,
		ClientID:  uuid.New(),
		CaseID:    uuid.New(),
		Title:     "Original",
		Type:      contractDomain.ContractTypeFixedFee,
		Amount:    decimal.NewFromFloat(10000),
		StartDate: now,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	var capturedEvent *timelineDomain.TimelineEvent
	repo := &mockUpdateContractRepo{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*contractDomain.Contract, error) {
			return contract, nil
		},
		UpdateFunc: func(ctx context.Context, c *contractDomain.Contract) error {
			c.UpdatedAt = time.Now()
			return nil
		},
	}
	timelineRepo := &mockUpdateTimelineRepo{
		CreateFunc: func(ctx context.Context, event *timelineDomain.TimelineEvent) error {
			capturedEvent = event
			return nil
		},
	}

	uc := NewUpdateContractUseCase(repo, timelineRepo)
	_, err := uc.Execute(context.Background(), UpdateContractInput{
		ID:        contractID.String(),
		Title:     "Contrato Atualizado",
		Type:      "fixed_fee",
		Amount:    "20000.00",
		StartDate: now.Format(time.RFC3339),
		Active:    true,
	})

	require.NoError(t, err)
	require.NotNil(t, capturedEvent)
	assert.Equal(t, timelineDomain.EventContractUpdated, capturedEvent.Type)
	assert.Contains(t, capturedEvent.Description, "Contrato Atualizado")
}

func TestUpdateContract_ShouldPopulateMetadata(t *testing.T) {
	now := time.Now().UTC()
	contractID := uuid.New()

	contract := &contractDomain.Contract{
		ID:        contractID,
		ClientID:  uuid.New(),
		CaseID:    uuid.New(),
		Title:     "Original",
		Type:      contractDomain.ContractTypeFixedFee,
		Amount:    decimal.NewFromFloat(10000),
		StartDate: now,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	var capturedEvent *timelineDomain.TimelineEvent
	repo := &mockUpdateContractRepo{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*contractDomain.Contract, error) {
			return contract, nil
		},
		UpdateFunc: func(ctx context.Context, c *contractDomain.Contract) error {
			c.UpdatedAt = time.Now()
			return nil
		},
	}
	timelineRepo := &mockUpdateTimelineRepo{
		CreateFunc: func(ctx context.Context, event *timelineDomain.TimelineEvent) error {
			capturedEvent = event
			return nil
		},
	}

	uc := NewUpdateContractUseCase(repo, timelineRepo)
	_, err := uc.Execute(context.Background(), UpdateContractInput{
		ID:        contractID.String(),
		Title:     "Contrato Atualizado",
		Type:      "fixed_fee",
		Amount:    "20000.00",
		StartDate: now.Format(time.RFC3339),
		Active:    true,
	})

	require.NoError(t, err)
	require.NotNil(t, capturedEvent)
	require.NotNil(t, capturedEvent.Metadata)

	var meta timelineDomain.ContractUpdatedMetadata
	err = json.Unmarshal(capturedEvent.Metadata, &meta)
	require.NoError(t, err)
	assert.Equal(t, contractID.String(), meta.ContractID)
	assert.Equal(t, "10000", meta.OldAmount)
	assert.Equal(t, "20000", meta.NewAmount)
	assert.Equal(t, "fixed_fee", meta.ContractType)
}
