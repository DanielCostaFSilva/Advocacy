package contract

import (
	"context"
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

type mockCloseContractRepo struct {
	FindByIDFunc func(ctx context.Context, id uuid.UUID) (*contractDomain.Contract, error)
	CloseFunc    func(ctx context.Context, id uuid.UUID) error
}

func (m *mockCloseContractRepo) Create(ctx context.Context, contract *contractDomain.Contract) error {
	return nil
}

func (m *mockCloseContractRepo) FindByID(ctx context.Context, id uuid.UUID) (*contractDomain.Contract, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockCloseContractRepo) ListByCaseID(ctx context.Context, caseID uuid.UUID) ([]contractDomain.Contract, error) {
	return nil, nil
}

func (m *mockCloseContractRepo) List(ctx context.Context, params contractDomain.ListContractsParams) ([]contractDomain.Contract, int64, error) {
	return nil, 0, nil
}

func (m *mockCloseContractRepo) Update(ctx context.Context, contract *contractDomain.Contract) error {
	return nil
}

func (m *mockCloseContractRepo) Close(ctx context.Context, id uuid.UUID) error {
	if m.CloseFunc != nil {
		return m.CloseFunc(ctx, id)
	}
	return nil
}

type mockCloseTimelineRepo struct {
	CreateFunc func(ctx context.Context, event *timelineDomain.TimelineEvent) error
}

func (m *mockCloseTimelineRepo) Create(ctx context.Context, event *timelineDomain.TimelineEvent) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, event)
	}
	return nil
}

func TestCloseContract_ShouldCloseSuccessfully(t *testing.T) {
	contractID := uuid.New()
	contract := &contractDomain.Contract{
		ID:        contractID,
		ClientID:  uuid.New(),
		CaseID:    uuid.New(),
		Title:     "Contrato Ativo",
		Type:      contractDomain.ContractTypeFixedFee,
		Amount:    decimal.NewFromFloat(15000),
		StartDate: time.Now().UTC(),
		Active:    true,
	}

	repo := &mockCloseContractRepo{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*contractDomain.Contract, error) {
			return contract, nil
		},
		CloseFunc: func(ctx context.Context, id uuid.UUID) error {
			return nil
		},
	}
	timelineRepo := &mockCloseTimelineRepo{
		CreateFunc: func(ctx context.Context, event *timelineDomain.TimelineEvent) error {
			return nil
		},
	}

	uc := NewCloseContractUseCase(repo, timelineRepo)
	err := uc.Execute(context.Background(), CloseContractInput{ContractID: contractID.String()})

	require.NoError(t, err)
}

func TestCloseContract_ShouldSetActiveFalse(t *testing.T) {
	contractID := uuid.New()
	contract := &contractDomain.Contract{
		ID:        contractID,
		ClientID:  uuid.New(),
		CaseID:    uuid.New(),
		Title:     "Contrato Ativo",
		Type:      contractDomain.ContractTypeFixedFee,
		Amount:    decimal.NewFromFloat(15000),
		StartDate: time.Now().UTC(),
		Active:    true,
	}

	repo := &mockCloseContractRepo{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*contractDomain.Contract, error) {
			return contract, nil
		},
		CloseFunc: func(ctx context.Context, id uuid.UUID) error {
			contract.Active = false
			return nil
		},
	}
	timelineRepo := &mockCloseTimelineRepo{
		CreateFunc: func(ctx context.Context, event *timelineDomain.TimelineEvent) error {
			return nil
		},
	}

	uc := NewCloseContractUseCase(repo, timelineRepo)
	err := uc.Execute(context.Background(), CloseContractInput{ContractID: contractID.String()})

	require.NoError(t, err)
	assert.False(t, contract.Active)
}

func TestCloseContract_ShouldSetEndDate(t *testing.T) {
	contractID := uuid.New()
	contract := &contractDomain.Contract{
		ID:        contractID,
		ClientID:  uuid.New(),
		CaseID:    uuid.New(),
		Title:     "Contrato Ativo",
		Type:      contractDomain.ContractTypeFixedFee,
		Amount:    decimal.NewFromFloat(15000),
		StartDate: time.Now().UTC(),
		Active:    true,
	}

	repo := &mockCloseContractRepo{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*contractDomain.Contract, error) {
			return contract, nil
		},
		CloseFunc: func(ctx context.Context, id uuid.UUID) error {
			contract.Active = false
			now := time.Now().UTC()
			contract.EndDate = &now
			return nil
		},
	}
	timelineRepo := &mockCloseTimelineRepo{
		CreateFunc: func(ctx context.Context, event *timelineDomain.TimelineEvent) error {
			return nil
		},
	}

	uc := NewCloseContractUseCase(repo, timelineRepo)
	err := uc.Execute(context.Background(), CloseContractInput{ContractID: contractID.String()})

	require.NoError(t, err)
	assert.NotNil(t, contract.EndDate)
}

func TestCloseContract_ShouldReturnErrInvalidContractID(t *testing.T) {
	repo := &mockCloseContractRepo{}
	timelineRepo := &mockCloseTimelineRepo{}
	uc := NewCloseContractUseCase(repo, timelineRepo)

	err := uc.Execute(context.Background(), CloseContractInput{ContractID: "invalid-uuid"})

	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidContractID))
}

func TestCloseContract_ShouldReturnErrContractNotFound(t *testing.T) {
	repo := &mockCloseContractRepo{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*contractDomain.Contract, error) {
			return nil, nil
		},
	}
	timelineRepo := &mockCloseTimelineRepo{}
	uc := NewCloseContractUseCase(repo, timelineRepo)

	err := uc.Execute(context.Background(), CloseContractInput{ContractID: uuid.New().String()})

	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrContractNotFound))
}

func TestCloseContract_ShouldReturnErrAlreadyClosed(t *testing.T) {
	contractID := uuid.New()
	contract := &contractDomain.Contract{
		ID:        contractID,
		ClientID:  uuid.New(),
		CaseID:    uuid.New(),
		Title:     "Contrato Encerrado",
		Type:      contractDomain.ContractTypeFixedFee,
		Amount:    decimal.NewFromFloat(15000),
		StartDate: time.Now().UTC(),
		Active:    false,
	}

	repo := &mockCloseContractRepo{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*contractDomain.Contract, error) {
			return contract, nil
		},
	}
	timelineRepo := &mockCloseTimelineRepo{}
	uc := NewCloseContractUseCase(repo, timelineRepo)

	err := uc.Execute(context.Background(), CloseContractInput{ContractID: contractID.String()})

	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrContractAlreadyClosed))
}

func TestCloseContract_ShouldReturnErrorWhenRepoFails(t *testing.T) {
	contractID := uuid.New()
	contract := &contractDomain.Contract{
		ID:        contractID,
		ClientID:  uuid.New(),
		CaseID:    uuid.New(),
		Title:     "Contrato Ativo",
		Type:      contractDomain.ContractTypeFixedFee,
		Amount:    decimal.NewFromFloat(15000),
		StartDate: time.Now().UTC(),
		Active:    true,
	}

	repo := &mockCloseContractRepo{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*contractDomain.Contract, error) {
			return contract, nil
		},
		CloseFunc: func(ctx context.Context, id uuid.UUID) error {
			return errors.New("db error")
		},
	}
	timelineRepo := &mockCloseTimelineRepo{}
	uc := NewCloseContractUseCase(repo, timelineRepo)

	err := uc.Execute(context.Background(), CloseContractInput{ContractID: contractID.String()})

	require.Error(t, err)
}

func TestCloseContract_ShouldReturnErrorWhenTimelineFails(t *testing.T) {
	contractID := uuid.New()
	contract := &contractDomain.Contract{
		ID:        contractID,
		ClientID:  uuid.New(),
		CaseID:    uuid.New(),
		Title:     "Contrato Ativo",
		Type:      contractDomain.ContractTypeFixedFee,
		Amount:    decimal.NewFromFloat(15000),
		StartDate: time.Now().UTC(),
		Active:    true,
	}

	repo := &mockCloseContractRepo{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*contractDomain.Contract, error) {
			return contract, nil
		},
		CloseFunc: func(ctx context.Context, id uuid.UUID) error {
			return nil
		},
	}
	timelineRepo := &mockCloseTimelineRepo{
		CreateFunc: func(ctx context.Context, event *timelineDomain.TimelineEvent) error {
			return errors.New("timeline error")
		},
	}

	uc := NewCloseContractUseCase(repo, timelineRepo)
	err := uc.Execute(context.Background(), CloseContractInput{ContractID: contractID.String()})

	require.Error(t, err)
}

func TestCloseContract_ShouldGenerateCONTRACTCLOSEDEvent(t *testing.T) {
	contractID := uuid.New()
	contract := &contractDomain.Contract{
		ID:        contractID,
		ClientID:  uuid.New(),
		CaseID:    uuid.New(),
		Title:     "Contrato Ativo",
		Type:      contractDomain.ContractTypeFixedFee,
		Amount:    decimal.NewFromFloat(15000),
		StartDate: time.Now().UTC(),
		Active:    true,
	}

	var capturedEvent *timelineDomain.TimelineEvent
	repo := &mockCloseContractRepo{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*contractDomain.Contract, error) {
			return contract, nil
		},
		CloseFunc: func(ctx context.Context, id uuid.UUID) error {
			contract.Active = false
			return nil
		},
	}
	timelineRepo := &mockCloseTimelineRepo{
		CreateFunc: func(ctx context.Context, event *timelineDomain.TimelineEvent) error {
			capturedEvent = event
			return nil
		},
	}

	uc := NewCloseContractUseCase(repo, timelineRepo)
	err := uc.Execute(context.Background(), CloseContractInput{ContractID: contractID.String()})

	require.NoError(t, err)
	require.NotNil(t, capturedEvent)
	assert.Equal(t, timelineDomain.EventContractClosed, capturedEvent.Type)
	assert.Contains(t, capturedEvent.Description, "Contrato Ativo")
}
