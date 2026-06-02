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

func TestUpdateClient_ShouldUpdateSuccessfully(t *testing.T) {
	id := uuid.New()
	repo := &domain.MockClientRepository{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Client, error) {
			return &domain.Client{
				ID:    id,
				Name:  "Maria Silva",
				CPF:   "12345678900",
				Email: "maria@example.com",
				Phone: "81999999999",
			}, nil
		},
		UpdateFunc: func(ctx context.Context, client *domain.Client) error {
			return nil
		},
	}

	uc := NewUpdateClientUseCase(repo)
	output, err := uc.Execute(context.Background(), UpdateClientInput{
		ID:    id.String(),
		Name:  "Maria da Silva",
		Email: "maria.silva@example.com",
		Phone: "81988887777",
	})

	require.NoError(t, err)
	require.NotNil(t, output)
	assert.Equal(t, id.String(), output.ID)
	assert.Equal(t, "Maria da Silva", output.Name)
	assert.Equal(t, "12345678900", output.CPF)
	assert.Equal(t, "maria.silva@example.com", output.Email)
	assert.Equal(t, "81988887777", output.Phone)
}

func TestUpdateClient_ShouldReturnErrInvalidClientIDWhenInvalidUUID(t *testing.T) {
	repo := &domain.MockClientRepository{}
	uc := NewUpdateClientUseCase(repo)
	output, err := uc.Execute(context.Background(), UpdateClientInput{
		ID:    "not-a-uuid",
		Name:  "Maria",
		Email: "maria@example.com",
		Phone: "81999999999",
	})

	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidClientID))
	assert.Nil(t, output)
}

func TestUpdateClient_ShouldReturnErrClientNotFoundWhenNotExists(t *testing.T) {
	repo := &domain.MockClientRepository{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Client, error) {
			return nil, nil
		},
	}

	uc := NewUpdateClientUseCase(repo)
	output, err := uc.Execute(context.Background(), UpdateClientInput{
		ID:    uuid.New().String(),
		Name:  "Maria",
		Email: "maria@example.com",
		Phone: "81999999999",
	})

	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrClientNotFound))
	assert.Nil(t, output)
}

func TestUpdateClient_ShouldReturnValidationErrorWhenInvalidEmail(t *testing.T) {
	repo := &domain.MockClientRepository{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Client, error) {
			return &domain.Client{
				ID:   id,
				Name: "Maria",
				CPF:  "12345678900",
			}, nil
		},
	}

	uc := NewUpdateClientUseCase(repo)
	output, err := uc.Execute(context.Background(), UpdateClientInput{
		ID:    uuid.New().String(),
		Name:  "Maria",
		Email: "invalid-email",
		Phone: "81999999999",
	})

	require.Error(t, err)
	assert.Nil(t, output)
}

func TestUpdateClient_ShouldReturnErrorWhenRepoUpdateFails(t *testing.T) {
	repo := &domain.MockClientRepository{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Client, error) {
			return &domain.Client{
				ID:   id,
				Name: "Maria",
				CPF:  "12345678900",
			}, nil
		},
		UpdateFunc: func(ctx context.Context, client *domain.Client) error {
			return errors.New("db error")
		},
	}

	uc := NewUpdateClientUseCase(repo)
	output, err := uc.Execute(context.Background(), UpdateClientInput{
		ID:    uuid.New().String(),
		Name:  "Maria",
		Email: "maria@example.com",
		Phone: "81999999999",
	})

	require.Error(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "db error")
}

func TestUpdateClient_ShouldPreserveCPF(t *testing.T) {
	id := uuid.New()
	repo := &domain.MockClientRepository{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Client, error) {
			return &domain.Client{
				ID:    id,
				Name:  "Maria Silva",
				CPF:   "12345678900",
				Email: "maria@example.com",
				Phone: "81999999999",
			}, nil
		},
		UpdateFunc: func(ctx context.Context, client *domain.Client) error {
			assert.Equal(t, "12345678900", client.CPF, "CPF must not change")
			return nil
		},
	}

	uc := NewUpdateClientUseCase(repo)
	output, err := uc.Execute(context.Background(), UpdateClientInput{
		ID:    id.String(),
		Name:  "Maria da Silva",
		Email: "maria.silva@example.com",
		Phone: "81988887777",
	})

	require.NoError(t, err)
	assert.Equal(t, "12345678900", output.CPF)
}
