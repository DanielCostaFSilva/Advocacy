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
	clientDomain "legalflow/internal/domain/client"
	contractDomain "legalflow/internal/domain/contract"
	caseDomain "legalflow/internal/domain/legalcase"
	timelineDomain "legalflow/internal/domain/timeline"
)

type mockClientRepository struct {
	FindByIDFunc func(ctx context.Context, id uuid.UUID) (*clientDomain.Client, error)
}

func (m *mockClientRepository) FindByID(ctx context.Context, id uuid.UUID) (*clientDomain.Client, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}

type mockCaseRepository struct {
	FindByIDFunc func(ctx context.Context, id uuid.UUID) (*caseDomain.Case, error)
}

func (m *mockCaseRepository) FindByID(ctx context.Context, id uuid.UUID) (*caseDomain.Case, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}

type mockContractRepository struct {
	CreateFunc func(ctx context.Context, contract *contractDomain.Contract) error
}

func (m *mockContractRepository) Create(ctx context.Context, contract *contractDomain.Contract) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, contract)
	}
	return nil
}

type mockTimelineRepository struct {
	CreateFunc func(ctx context.Context, event *timelineDomain.TimelineEvent) error
}

func (m *mockTimelineRepository) Create(ctx context.Context, event *timelineDomain.TimelineEvent) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, event)
	}
	return nil
}

func validInput() CreateContractInput {
	return CreateContractInput{
		ClientID:    uuid.New().String(),
		CaseID:      uuid.New().String(),
		Title:       "Contrato de Honorários",
		Description: "Descrição do contrato",
		Type:        "fixed_fee",
		Amount:      "15000.00",
		StartDate:   time.Now().AddDate(0, 0, 1).Format(time.RFC3339),
	}
}

func TestCreateContract_ShouldSucceed(t *testing.T) {
	clientID := uuid.New()
	caseID := uuid.New()

	clientRepo := &mockClientRepository{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*clientDomain.Client, error) {
			return &clientDomain.Client{ID: clientID}, nil
		},
	}
	caseRepo := &mockCaseRepository{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*caseDomain.Case, error) {
			return &caseDomain.Case{ID: caseID}, nil
		},
	}
	contractRepo := &mockContractRepository{
		CreateFunc: func(ctx context.Context, contract *contractDomain.Contract) error {
			contract.ID = uuid.New()
			return nil
		},
	}
	timelineRepo := &mockTimelineRepository{
		CreateFunc: func(ctx context.Context, event *timelineDomain.TimelineEvent) error {
			return nil
		},
	}

	uc := NewCreateContractUseCase(clientRepo, caseRepo, contractRepo, timelineRepo)
	input := validInput()
	input.ClientID = clientID.String()
	input.CaseID = caseID.String()

	output, err := uc.Execute(context.Background(), input)

	require.NoError(t, err)
	require.NotNil(t, output)
	assert.NotEmpty(t, output.ID)
	assert.Equal(t, clientID.String(), output.ClientID)
	assert.Equal(t, caseID.String(), output.CaseID)
	assert.Equal(t, "Contrato de Honorários", output.Title)
	assert.Equal(t, "fixed_fee", output.Type)
	assert.Equal(t, "15000", output.Amount)
	assert.True(t, output.Active)
}

func TestCreateContract_ShouldReturnErrClientNotFound(t *testing.T) {
	uc := NewCreateContractUseCase(
		&mockClientRepository{
			FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*clientDomain.Client, error) {
				return nil, nil
			},
		},
		&mockCaseRepository{},
		&mockContractRepository{},
		&mockTimelineRepository{},
	)

	output, err := uc.Execute(context.Background(), validInput())

	require.Error(t, err)
	assert.Nil(t, output)
	assert.ErrorIs(t, err, ErrClientNotFound)
}

func TestCreateContract_ShouldReturnErrCaseNotFound(t *testing.T) {
	clientID := uuid.New()

	uc := NewCreateContractUseCase(
		&mockClientRepository{
			FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*clientDomain.Client, error) {
				return &clientDomain.Client{ID: clientID}, nil
			},
		},
		&mockCaseRepository{
			FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*caseDomain.Case, error) {
				return nil, nil
			},
		},
		&mockContractRepository{},
		&mockTimelineRepository{},
	)

	input := validInput()
	input.ClientID = clientID.String()
	output, err := uc.Execute(context.Background(), input)

	require.Error(t, err)
	assert.Nil(t, output)
	assert.ErrorIs(t, err, ErrCaseNotFound)
}

func TestCreateContract_ShouldReturnErrInvalidClientID(t *testing.T) {
	uc := NewCreateContractUseCase(
		&mockClientRepository{}, &mockCaseRepository{},
		&mockContractRepository{}, &mockTimelineRepository{},
	)

	input := validInput()
	input.ClientID = "not-a-uuid"
	output, err := uc.Execute(context.Background(), input)

	require.Error(t, err)
	assert.Nil(t, output)
	assert.ErrorIs(t, err, ErrInvalidClientID)
}

func TestCreateContract_ShouldReturnErrInvalidCaseID(t *testing.T) {
	uc := NewCreateContractUseCase(
		&mockClientRepository{}, &mockCaseRepository{},
		&mockContractRepository{}, &mockTimelineRepository{},
	)

	input := validInput()
	input.CaseID = "not-a-uuid"
	output, err := uc.Execute(context.Background(), input)

	require.Error(t, err)
	assert.Nil(t, output)
	assert.ErrorIs(t, err, ErrInvalidCaseID)
}

func TestCreateContract_ShouldReturnErrInvalidInputWhenInvalidData(t *testing.T) {
	clientID := uuid.New()

	uc := NewCreateContractUseCase(
		&mockClientRepository{
			FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*clientDomain.Client, error) {
				return &clientDomain.Client{ID: clientID}, nil
			},
		},
		&mockCaseRepository{
			FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*caseDomain.Case, error) {
				return &caseDomain.Case{ID: uuid.New()}, nil
			},
		},
		&mockContractRepository{},
		&mockTimelineRepository{},
	)

	input := validInput()
	input.ClientID = clientID.String()
	input.Amount = "not-a-number"
	output, err := uc.Execute(context.Background(), input)

	require.Error(t, err)
	assert.Nil(t, output)
	assert.ErrorIs(t, err, ErrInvalidInput)
}

func TestCreateContract_ShouldReturnErrorWhenContractRepoFails(t *testing.T) {
	clientID := uuid.New()
	caseID := uuid.New()

	uc := NewCreateContractUseCase(
		&mockClientRepository{
			FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*clientDomain.Client, error) {
				return &clientDomain.Client{ID: clientID}, nil
			},
		},
		&mockCaseRepository{
			FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*caseDomain.Case, error) {
				return &caseDomain.Case{ID: caseID}, nil
			},
		},
		&mockContractRepository{
			CreateFunc: func(ctx context.Context, contract *contractDomain.Contract) error {
				return errors.New("db error")
			},
		},
		&mockTimelineRepository{},
	)

	input := validInput()
	input.ClientID = clientID.String()
	input.CaseID = caseID.String()
	output, err := uc.Execute(context.Background(), input)

	require.Error(t, err)
	assert.Nil(t, output)
	assert.ErrorIs(t, err, ErrContractPersistence)
}

func TestCreateContract_ShouldReturnErrorWhenTimelineRepoFails(t *testing.T) {
	clientID := uuid.New()
	caseID := uuid.New()

	uc := NewCreateContractUseCase(
		&mockClientRepository{
			FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*clientDomain.Client, error) {
				return &clientDomain.Client{ID: clientID}, nil
			},
		},
		&mockCaseRepository{
			FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*caseDomain.Case, error) {
				return &caseDomain.Case{ID: caseID}, nil
			},
		},
		&mockContractRepository{
			CreateFunc: func(ctx context.Context, contract *contractDomain.Contract) error {
				contract.ID = uuid.New()
				return nil
			},
		},
		&mockTimelineRepository{
			CreateFunc: func(ctx context.Context, event *timelineDomain.TimelineEvent) error {
				return errors.New("timeline error")
			},
		},
	)

	input := validInput()
	input.ClientID = clientID.String()
	input.CaseID = caseID.String()
	output, err := uc.Execute(context.Background(), input)

	require.Error(t, err)
	assert.Nil(t, output)
	assert.ErrorIs(t, err, ErrTimelinePersistence)
}

func TestCreateContract_ShouldCreateTimelineEvent(t *testing.T) {
	clientID := uuid.New()
	caseID := uuid.New()
	var capturedEvent *timelineDomain.TimelineEvent

	uc := NewCreateContractUseCase(
		&mockClientRepository{
			FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*clientDomain.Client, error) {
				return &clientDomain.Client{ID: clientID}, nil
			},
		},
		&mockCaseRepository{
			FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*caseDomain.Case, error) {
				return &caseDomain.Case{ID: caseID}, nil
			},
		},
		&mockContractRepository{
			CreateFunc: func(ctx context.Context, contract *contractDomain.Contract) error {
				contract.ID = uuid.New()
				return nil
			},
		},
		&mockTimelineRepository{
			CreateFunc: func(ctx context.Context, event *timelineDomain.TimelineEvent) error {
				capturedEvent = event
				return nil
			},
		},
	)

	input := validInput()
	input.ClientID = clientID.String()
	input.CaseID = caseID.String()
	_, err := uc.Execute(context.Background(), input)

	require.NoError(t, err)
	require.NotNil(t, capturedEvent)
	assert.Equal(t, timelineDomain.EventContractCreated, capturedEvent.Type)
	assert.Equal(t, caseID, capturedEvent.CaseID)
	assert.Contains(t, capturedEvent.Description, "Contrato de Honorários")
}

func TestCreateContract_ShouldPopulateMetadata(t *testing.T) {
	clientID := uuid.New()
	caseID := uuid.New()
	var capturedEvent *timelineDomain.TimelineEvent

	uc := NewCreateContractUseCase(
		&mockClientRepository{
			FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*clientDomain.Client, error) {
				return &clientDomain.Client{ID: clientID}, nil
			},
		},
		&mockCaseRepository{
			FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*caseDomain.Case, error) {
				return &caseDomain.Case{ID: caseID}, nil
			},
		},
		&mockContractRepository{
			CreateFunc: func(ctx context.Context, contract *contractDomain.Contract) error {
				contract.ID = uuid.New()
				contract.Amount = decimal.NewFromFloat(15000)
				return nil
			},
		},
		&mockTimelineRepository{
			CreateFunc: func(ctx context.Context, event *timelineDomain.TimelineEvent) error {
				capturedEvent = event
				return nil
			},
		},
	)

	input := validInput()
	input.ClientID = clientID.String()
	input.CaseID = caseID.String()
	_, err := uc.Execute(context.Background(), input)

	require.NoError(t, err)
	require.NotNil(t, capturedEvent)
	require.NotEmpty(t, capturedEvent.Metadata)

	var meta timelineDomain.ContractMetadata
	err = json.Unmarshal(capturedEvent.Metadata, &meta)
	require.NoError(t, err)
	assert.NotEmpty(t, meta.ContractID)
	assert.Equal(t, "fixed_fee", meta.ContractType)
	assert.NotEmpty(t, meta.Amount)
}
