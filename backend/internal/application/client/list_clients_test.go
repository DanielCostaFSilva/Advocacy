package client

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domain "legalflow/internal/domain/client"
)

func TestListClients_ShouldReturnPaginatedResult(t *testing.T) {
	clients := []domain.Client{
		{ID: uuid.New(), Name: "Alice", CPF: "11111111111"},
		{ID: uuid.New(), Name: "Bob", CPF: "22222222222"},
	}
	repo := &domain.MockClientRepository{
		ListFunc: func(ctx context.Context, params domain.ListClientsParams) ([]domain.Client, int64, error) {
			return clients, 2, nil
		},
	}

	uc := NewListClientsUseCase(repo)
	output, err := uc.Execute(context.Background(), ListClientsInput{
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

func TestListClients_ShouldCapPageSizeAt100(t *testing.T) {
	repo := &domain.MockClientRepository{
		ListFunc: func(ctx context.Context, params domain.ListClientsParams) ([]domain.Client, int64, error) {
			assert.Equal(t, 100, params.Limit)
			return nil, 0, nil
		},
	}

	uc := NewListClientsUseCase(repo)
	output, err := uc.Execute(context.Background(), ListClientsInput{
		Page:     1,
		PageSize: 200,
	})

	require.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, 100, output.PageSize)
}

func TestListClients_ShouldUseDefaultPageWhenInvalid(t *testing.T) {
	repo := &domain.MockClientRepository{
		ListFunc: func(ctx context.Context, params domain.ListClientsParams) ([]domain.Client, int64, error) {
			assert.Equal(t, 0, params.Offset)
			return nil, 0, nil
		},
	}

	uc := NewListClientsUseCase(repo)
	output, err := uc.Execute(context.Background(), ListClientsInput{
		Page:     0,
		PageSize: 20,
	})

	require.NoError(t, err)
	assert.Equal(t, 1, output.Page)
}

func TestListClients_ShouldReturnErrorWhenRepoFails(t *testing.T) {
	repo := &domain.MockClientRepository{
		ListFunc: func(ctx context.Context, params domain.ListClientsParams) ([]domain.Client, int64, error) {
			return nil, 0, errors.New("db error")
		},
	}

	uc := NewListClientsUseCase(repo)
	output, err := uc.Execute(context.Background(), ListClientsInput{
		Page:     1,
		PageSize: 20,
	})

	require.Error(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "db error")
}

func TestListClients_ShouldFilterByName(t *testing.T) {
	repo := &domain.MockClientRepository{
		ListFunc: func(ctx context.Context, params domain.ListClientsParams) ([]domain.Client, int64, error) {
			assert.Equal(t, "maria", params.Name)
			return nil, 0, nil
		},
	}

	uc := NewListClientsUseCase(repo)
	_, err := uc.Execute(context.Background(), ListClientsInput{
		Page:     1,
		PageSize: 20,
		Name:     "maria",
	})

	require.NoError(t, err)
}

func TestListClients_ShouldFilterByCPF(t *testing.T) {
	repo := &domain.MockClientRepository{
		ListFunc: func(ctx context.Context, params domain.ListClientsParams) ([]domain.Client, int64, error) {
			assert.Equal(t, "12345678900", params.CPF)
			return nil, 0, nil
		},
	}

	uc := NewListClientsUseCase(repo)
	_, err := uc.Execute(context.Background(), ListClientsInput{
		Page:     1,
		PageSize: 20,
		CPF:      "12345678900",
	})

	require.NoError(t, err)
}

func TestListClients_ShouldCalculateCorrectOffset(t *testing.T) {
	var capturedParams domain.ListClientsParams
	repo := &domain.MockClientRepository{
		ListFunc: func(ctx context.Context, params domain.ListClientsParams) ([]domain.Client, int64, error) {
			capturedParams = params
			return nil, 0, nil
		},
	}

	uc := NewListClientsUseCase(repo)
	_, err := uc.Execute(context.Background(), ListClientsInput{
		Page:     3,
		PageSize: 15,
	})

	require.NoError(t, err)
	assert.Equal(t, 15, capturedParams.Limit)
	assert.Equal(t, 30, capturedParams.Offset)
}

func TestListClients_ShouldCalculateTotalPages(t *testing.T) {
	repo := &domain.MockClientRepository{
		ListFunc: func(ctx context.Context, params domain.ListClientsParams) ([]domain.Client, int64, error) {
			clients := make([]domain.Client, 10)
			for i := range clients {
				clients[i] = domain.Client{ID: uuid.New(), Name: "Client"}
			}
			return clients, 25, nil
		},
	}

	uc := NewListClientsUseCase(repo)
	output, err := uc.Execute(context.Background(), ListClientsInput{
		Page:     1,
		PageSize: 10,
	})

	require.NoError(t, err)
	assert.Equal(t, int64(25), output.Total)
	assert.Equal(t, 3, output.TotalPages)
}
