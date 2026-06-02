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

func TestDeleteClient_ShouldSoftDeleteSuccessfully(t *testing.T) {
	id := uuid.New()
	repo := &domain.MockClientRepository{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Client, error) {
			return &domain.Client{ID: id, Name: "Maria", CPF: "12345678900"}, nil
		},
		SoftDeleteFunc: func(ctx context.Context, id uuid.UUID) error {
			return nil
		},
	}

	uc := NewDeleteClientUseCase(repo)
	err := uc.Execute(context.Background(), DeleteClientInput{ID: id.String()})

	require.NoError(t, err)
}

func TestDeleteClient_ShouldReturnErrInvalidClientIDWhenInvalidUUID(t *testing.T) {
	repo := &domain.MockClientRepository{}
	uc := NewDeleteClientUseCase(repo)
	err := uc.Execute(context.Background(), DeleteClientInput{ID: "not-a-uuid"})

	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidClientID))
}

func TestDeleteClient_ShouldReturnErrClientNotFoundWhenNotExists(t *testing.T) {
	repo := &domain.MockClientRepository{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Client, error) {
			return nil, nil
		},
	}

	uc := NewDeleteClientUseCase(repo)
	err := uc.Execute(context.Background(), DeleteClientInput{ID: uuid.New().String()})

	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrClientNotFound))
}

func TestDeleteClient_ShouldReturnErrorWhenRepositoryFails(t *testing.T) {
	repo := &domain.MockClientRepository{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Client, error) {
			return &domain.Client{ID: id, Name: "Maria", CPF: "12345678900"}, nil
		},
		SoftDeleteFunc: func(ctx context.Context, id uuid.UUID) error {
			return errors.New("db error")
		},
	}

	uc := NewDeleteClientUseCase(repo)
	err := uc.Execute(context.Background(), DeleteClientInput{ID: uuid.New().String()})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
}
