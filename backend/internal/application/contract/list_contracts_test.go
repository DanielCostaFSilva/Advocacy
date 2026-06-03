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

type mockContractRepo struct {
	ListFunc func(ctx context.Context, params domain.ListContractsParams) ([]domain.Contract, int64, error)
}

func (m *mockContractRepo) Create(ctx context.Context, contract *domain.Contract) error {
	return nil
}

func (m *mockContractRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Contract, error) {
	return nil, nil
}

func (m *mockContractRepo) ListByCaseID(ctx context.Context, caseID uuid.UUID) ([]domain.Contract, error) {
	return nil, nil
}

func (m *mockContractRepo) List(ctx context.Context, params domain.ListContractsParams) ([]domain.Contract, int64, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, params)
	}
	return nil, 0, nil
}

func (m *mockContractRepo) Update(ctx context.Context, contract *domain.Contract) error {
	return nil
}

func TestListContracts_ShouldReturnPaginatedResult(t *testing.T) {
	now := time.Now().UTC()
	contracts := []domain.Contract{
		{ID: uuid.New(), ClientID: uuid.New(), CaseID: uuid.New(), Title: "Contrato 1", Type: domain.ContractTypeFixedFee, Amount: decimal.NewFromFloat(10000), StartDate: now, Active: true, CreatedAt: now},
		{ID: uuid.New(), ClientID: uuid.New(), CaseID: uuid.New(), Title: "Contrato 2", Type: domain.ContractTypeHourly, Amount: decimal.NewFromFloat(5000), StartDate: now, Active: true, CreatedAt: now},
	}
	repo := &mockContractRepo{
		ListFunc: func(ctx context.Context, params domain.ListContractsParams) ([]domain.Contract, int64, error) {
			return contracts, 2, nil
		},
	}

	uc := NewListContractsUseCase(repo)
	output, err := uc.Execute(context.Background(), ListContractsInput{
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

func TestListContracts_ShouldCapPageSizeAt100(t *testing.T) {
	repo := &mockContractRepo{
		ListFunc: func(ctx context.Context, params domain.ListContractsParams) ([]domain.Contract, int64, error) {
			assert.Equal(t, 100, params.Limit)
			return nil, 0, nil
		},
	}

	uc := NewListContractsUseCase(repo)
	output, err := uc.Execute(context.Background(), ListContractsInput{
		Page:     1,
		PageSize: 200,
	})

	require.NoError(t, err)
	assert.Equal(t, 100, output.PageSize)
}

func TestListContracts_ShouldUseDefaultPageWhenInvalid(t *testing.T) {
	repo := &mockContractRepo{
		ListFunc: func(ctx context.Context, params domain.ListContractsParams) ([]domain.Contract, int64, error) {
			assert.Equal(t, 0, params.Offset)
			return nil, 0, nil
		},
	}

	uc := NewListContractsUseCase(repo)
	output, err := uc.Execute(context.Background(), ListContractsInput{
		Page:     0,
		PageSize: 20,
	})

	require.NoError(t, err)
	assert.Equal(t, 1, output.Page)
}

func TestListContracts_ShouldFilterByClientID(t *testing.T) {
	clientID := uuid.New().String()
	repo := &mockContractRepo{
		ListFunc: func(ctx context.Context, params domain.ListContractsParams) ([]domain.Contract, int64, error) {
			assert.Equal(t, clientID, params.ClientID)
			return nil, 0, nil
		},
	}

	uc := NewListContractsUseCase(repo)
	_, err := uc.Execute(context.Background(), ListContractsInput{
		Page: 1, PageSize: 20,
		ClientID: clientID,
	})
	require.NoError(t, err)
}

func TestListContracts_ShouldFilterByCaseID(t *testing.T) {
	caseID := uuid.New().String()
	repo := &mockContractRepo{
		ListFunc: func(ctx context.Context, params domain.ListContractsParams) ([]domain.Contract, int64, error) {
			assert.Equal(t, caseID, params.CaseID)
			return nil, 0, nil
		},
	}

	uc := NewListContractsUseCase(repo)
	_, err := uc.Execute(context.Background(), ListContractsInput{
		Page: 1, PageSize: 20,
		CaseID: caseID,
	})
	require.NoError(t, err)
}

func TestListContracts_ShouldFilterByActive(t *testing.T) {
	repo := &mockContractRepo{
		ListFunc: func(ctx context.Context, params domain.ListContractsParams) ([]domain.Contract, int64, error) {
			require.NotNil(t, params.Active)
			assert.True(t, *params.Active)
			return nil, 0, nil
		},
	}

	uc := NewListContractsUseCase(repo)
	_, err := uc.Execute(context.Background(), ListContractsInput{
		Page: 1, PageSize: 20,
		Active: "true",
	})
	require.NoError(t, err)
}

func TestListContracts_ShouldReturnErrorWhenRepoFails(t *testing.T) {
	repo := &mockContractRepo{
		ListFunc: func(ctx context.Context, params domain.ListContractsParams) ([]domain.Contract, int64, error) {
			return nil, 0, errors.New("db error")
		},
	}

	uc := NewListContractsUseCase(repo)
	output, err := uc.Execute(context.Background(), ListContractsInput{
		Page: 1, PageSize: 20,
	})
	require.Error(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "db error")
}
