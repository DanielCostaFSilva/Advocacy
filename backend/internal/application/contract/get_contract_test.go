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
	domain "legalflow/internal/domain/contract"
)

type mockGetContractRepo struct {
	FindByIDFunc func(ctx context.Context, id uuid.UUID) (*domain.Contract, error)
}

func (m *mockGetContractRepo) Create(ctx context.Context, contract *domain.Contract) error {
	return nil
}

func (m *mockGetContractRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Contract, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockGetContractRepo) ListByCaseID(ctx context.Context, caseID uuid.UUID) ([]domain.Contract, error) {
	return nil, nil
}

func (m *mockGetContractRepo) List(ctx context.Context, params domain.ListContractsParams) ([]domain.Contract, int64, error) {
	return nil, 0, nil
}

func TestGetContract_ShouldReturnContract(t *testing.T) {
	now := time.Now().UTC()
	contractID := uuid.New()
	clientID := uuid.New()
	caseID := uuid.New()

	repo := &mockGetContractRepo{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Contract, error) {
			return &domain.Contract{
				ID:        contractID,
				ClientID:  clientID,
				CaseID:    caseID,
				Title:     "Contrato de Honorários",
				Type:      domain.ContractTypeFixedFee,
				Amount:    decimal.NewFromFloat(15000),
				StartDate: now,
				Active:    true,
				CreatedAt: now,
				UpdatedAt: now,
			}, nil
		},
	}

	uc := NewGetContractUseCase(repo)
	output, err := uc.Execute(context.Background(), GetContractInput{ID: contractID.String()})

	require.NoError(t, err)
	require.NotNil(t, output)
	assert.Equal(t, contractID.String(), output.ID)
	assert.Equal(t, clientID.String(), output.ClientID)
	assert.Equal(t, caseID.String(), output.CaseID)
	assert.Equal(t, "Contrato de Honorários", output.Title)
	assert.Equal(t, "fixed_fee", output.Type)
	assert.Equal(t, "15000", output.Amount)
	assert.True(t, output.Active)
}

func TestGetContract_ShouldReturnErrInvalidContractID(t *testing.T) {
	repo := &mockGetContractRepo{}
	uc := NewGetContractUseCase(repo)
	output, err := uc.Execute(context.Background(), GetContractInput{ID: "invalid-uuid"})

	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidContractID))
	assert.Nil(t, output)
}

func TestGetContract_ShouldReturnErrContractNotFound(t *testing.T) {
	repo := &mockGetContractRepo{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Contract, error) {
			return nil, nil
		},
	}

	uc := NewGetContractUseCase(repo)
	output, err := uc.Execute(context.Background(), GetContractInput{ID: uuid.New().String()})

	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrContractNotFound))
	assert.Nil(t, output)
}

func TestGetContract_ShouldReturnErrorWhenRepoFails(t *testing.T) {
	repo := &mockGetContractRepo{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Contract, error) {
			return nil, errors.New("db error")
		},
	}

	uc := NewGetContractUseCase(repo)
	output, err := uc.Execute(context.Background(), GetContractInput{ID: uuid.New().String()})

	require.Error(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "db error")
}

func TestGetContract_ShouldMapOutputCorrectly(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	endDate := now.AddDate(1, 0, 0)
	contractID := uuid.New()

	repo := &mockGetContractRepo{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Contract, error) {
			return &domain.Contract{
				ID:          contractID,
				ClientID:    uuid.New(),
				CaseID:      uuid.New(),
				Title:       "Contrato Teste",
				Description: "Descrição detalhada",
				Type:        domain.ContractTypeMonthly,
				Amount:      decimal.NewFromFloat(5000.50),
				StartDate:   now,
				EndDate:     &endDate,
				Active:      true,
				CreatedAt:   now,
				UpdatedAt:   now,
			}, nil
		},
	}

	uc := NewGetContractUseCase(repo)
	output, err := uc.Execute(context.Background(), GetContractInput{ID: contractID.String()})

	require.NoError(t, err)
	assert.Equal(t, "Descrição detalhada", output.Description)
	assert.Equal(t, "monthly", output.Type)
	assert.Equal(t, "5000.5", output.Amount)
	require.NotNil(t, output.EndDate)
	assert.Equal(t, endDate.Format(time.RFC3339), *output.EndDate)
	assert.Equal(t, now.Format(time.RFC3339), output.StartDate)
	assert.Equal(t, now, output.CreatedAt)
	assert.Equal(t, now, output.UpdatedAt)
}
