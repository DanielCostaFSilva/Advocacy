package client

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domain "legalflow/internal/domain/client"
)

func TestCreateClient_ShouldSucceedWhenValidInput(t *testing.T) {
	repo := &domain.MockClientRepository{
		FindByCPFFunc: func(ctx context.Context, cpf string) (*domain.Client, error) {
			return nil, nil
		},
	}

	uc := NewCreateClientUseCase(repo)
	input := CreateClientInput{
		Name:  "John Doe",
		CPF:   "12345678901",
		Email: "john@example.com",
		Phone: "11999999999",
	}

	output, err := uc.Execute(context.Background(), input)
	require.NoError(t, err)
	require.NotNil(t, output)
	assert.NotEmpty(t, output.ID)
	assert.Equal(t, "John Doe", output.Name)
	assert.Equal(t, "12345678901", output.CPF)
	assert.Equal(t, "john@example.com", output.Email)
	assert.Equal(t, "11999999999", output.Phone)
}

func TestCreateClient_ShouldReturnErrClientAlreadyExists(t *testing.T) {
	existing, _ := domain.NewClient("Existing", "99988877766", "existing@example.com", "")
	repo := &domain.MockClientRepository{
		FindByCPFFunc: func(ctx context.Context, cpf string) (*domain.Client, error) {
			return existing, nil
		},
	}

	uc := NewCreateClientUseCase(repo)
	input := CreateClientInput{
		Name:  "John Doe",
		CPF:   "99988877766",
		Email: "john@example.com",
		Phone: "11999999999",
	}

	output, err := uc.Execute(context.Background(), input)
	assert.ErrorIs(t, err, ErrClientAlreadyExists)
	assert.Nil(t, output)
}

func TestCreateClient_ShouldReturnDomainErrorWhenInvalidCPF(t *testing.T) {
	repo := &domain.MockClientRepository{
		FindByCPFFunc: func(ctx context.Context, cpf string) (*domain.Client, error) {
			return nil, nil
		},
	}

	uc := NewCreateClientUseCase(repo)
	input := CreateClientInput{
		Name:  "John Doe",
		CPF:   "123",
		Email: "john@example.com",
		Phone: "11999999999",
	}

	output, err := uc.Execute(context.Background(), input)
	require.Error(t, err)
	assert.Nil(t, output)
}

func TestCreateClient_ShouldReturnErrorWhenRepoCreateFails(t *testing.T) {
	repo := &domain.MockClientRepository{
		FindByCPFFunc: func(ctx context.Context, cpf string) (*domain.Client, error) {
			return nil, nil
		},
		CreateFunc: func(ctx context.Context, client *domain.Client) error {
			return errors.New("db error")
		},
	}

	uc := NewCreateClientUseCase(repo)
	input := CreateClientInput{
		Name:  "John Doe",
		CPF:   "12345678901",
		Email: "john@example.com",
		Phone: "11999999999",
	}

	output, err := uc.Execute(context.Background(), input)
	require.Error(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "db error")
}

func TestCreateClient_ShouldReturnErrorWhenRepoSearchFails(t *testing.T) {
	repo := &domain.MockClientRepository{
		FindByCPFFunc: func(ctx context.Context, cpf string) (*domain.Client, error) {
			return nil, errors.New("search error")
		},
	}

	uc := NewCreateClientUseCase(repo)
	input := CreateClientInput{
		Name:  "John Doe",
		CPF:   "12345678901",
		Email: "john@example.com",
		Phone: "11999999999",
	}

	output, err := uc.Execute(context.Background(), input)
	require.Error(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "search error")
}

func TestCreateClient_ShouldInvokeRepoCreateWithCorrectData(t *testing.T) {
	var capturedClient *domain.Client
	repo := &domain.MockClientRepository{
		FindByCPFFunc: func(ctx context.Context, cpf string) (*domain.Client, error) {
			return nil, nil
		},
		CreateFunc: func(ctx context.Context, client *domain.Client) error {
			capturedClient = client
			return nil
		},
	}

	uc := NewCreateClientUseCase(repo)
	input := CreateClientInput{
		Name:  "Jane Doe",
		CPF:   "11122233344",
		Email: "jane@example.com",
		Phone: "11988888888",
	}

	output, err := uc.Execute(context.Background(), input)
	require.NoError(t, err)
	require.NotNil(t, output)

	require.NotNil(t, capturedClient)
	assert.Equal(t, "Jane Doe", capturedClient.Name)
	assert.Equal(t, "11122233344", capturedClient.CPF)
	assert.Equal(t, "jane@example.com", capturedClient.Email)
	assert.Equal(t, "11988888888", capturedClient.Phone)
}
