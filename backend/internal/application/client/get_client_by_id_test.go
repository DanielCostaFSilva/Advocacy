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

func TestGetClientByID_ShouldReturnClientWhenExists(t *testing.T) {
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
	}

	uc := NewGetClientByIDUseCase(repo)
	output, err := uc.Execute(context.Background(), GetClientByIDInput{ID: id.String()})

	require.NoError(t, err)
	require.NotNil(t, output)
	assert.Equal(t, id.String(), output.ID)
	assert.Equal(t, "Maria Silva", output.Name)
	assert.Equal(t, "12345678900", output.CPF)
	assert.Equal(t, "maria@example.com", output.Email)
	assert.Equal(t, "81999999999", output.Phone)
}

func TestGetClientByID_ShouldReturnErrInvalidClientIDWhenInvalidUUID(t *testing.T) {
	repo := &domain.MockClientRepository{}
	uc := NewGetClientByIDUseCase(repo)
	output, err := uc.Execute(context.Background(), GetClientByIDInput{ID: "not-a-uuid"})

	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidClientID))
	assert.Nil(t, output)
}

func TestGetClientByID_ShouldReturnErrClientNotFoundWhenNotExists(t *testing.T) {
	repo := &domain.MockClientRepository{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Client, error) {
			return nil, nil
		},
	}

	uc := NewGetClientByIDUseCase(repo)
	output, err := uc.Execute(context.Background(), GetClientByIDInput{ID: uuid.New().String()})

	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrClientNotFound))
	assert.Nil(t, output)
}

func TestGetClientByID_ShouldReturnErrorWhenRepositoryFails(t *testing.T) {
	repo := &domain.MockClientRepository{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Client, error) {
			return nil, errors.New("db error")
		},
	}

	uc := NewGetClientByIDUseCase(repo)
	output, err := uc.Execute(context.Background(), GetClientByIDInput{ID: uuid.New().String()})

	require.Error(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "db error")
}
