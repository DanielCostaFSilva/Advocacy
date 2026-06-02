package user

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domain "legalflow/internal/user/domain/user"
)

func TestCreateUser_ShouldSucceedWhenValidInput(t *testing.T) {
	repo := &domain.MockRepository{
		FindByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
			return nil, nil
		},
	}
	hasher := &MockPasswordHasher{
		HashFunc: func(password string) (string, error) {
			return "hashed-" + password, nil
		},
	}

	uc := NewCreateUserUseCase(repo, hasher)
	input := CreateUserInput{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "securepass123",
	}

	output, err := uc.Execute(context.Background(), input)
	require.NoError(t, err)
	require.NotNil(t, output)
	assert.NotEmpty(t, output.ID)
	assert.Equal(t, "John Doe", output.Name)
	assert.Equal(t, "john@example.com", output.Email)
}

func TestCreateUser_ShouldReturnErrEmailAlreadyExists(t *testing.T) {
	existingUser, _ := domain.NewUser("Existing", "existing@example.com", "hash")
	repo := &domain.MockRepository{
		FindByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
			return existingUser, nil
		},
	}
	hasher := &MockPasswordHasher{
		HashFunc: func(password string) (string, error) {
			return "hashed-" + password, nil
		},
	}

	uc := NewCreateUserUseCase(repo, hasher)
	input := CreateUserInput{
		Name:     "John Doe",
		Email:    "existing@example.com",
		Password: "securepass123",
	}

	output, err := uc.Execute(context.Background(), input)
	assert.ErrorIs(t, err, ErrEmailAlreadyExists)
	assert.Nil(t, output)
}

func TestCreateUser_ShouldReturnErrHashPasswordWhenHashFails(t *testing.T) {
	repo := &domain.MockRepository{
		FindByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
			return nil, nil
		},
	}
	hasher := &MockPasswordHasher{
		HashFunc: func(password string) (string, error) {
			return "", errors.New("crypto failure")
		},
	}

	uc := NewCreateUserUseCase(repo, hasher)
	input := CreateUserInput{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "securepass123",
	}

	output, err := uc.Execute(context.Background(), input)
	assert.ErrorIs(t, err, ErrHashPassword)
	assert.Nil(t, output)
}

func TestCreateUser_ShouldReturnValidationErrorWhenPasswordTooShort(t *testing.T) {
	repo := &domain.MockRepository{}
	hasher := &MockPasswordHasher{}

	uc := NewCreateUserUseCase(repo, hasher)
	input := CreateUserInput{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "1234567",
	}

	output, err := uc.Execute(context.Background(), input)
	assert.ErrorIs(t, err, ErrInvalidInput)
	assert.Nil(t, output)
}

func TestCreateUser_ShouldReturnValidationErrorWhenNameEmpty(t *testing.T) {
	repo := &domain.MockRepository{}
	hasher := &MockPasswordHasher{}

	uc := NewCreateUserUseCase(repo, hasher)
	input := CreateUserInput{
		Name:     "",
		Email:    "john@example.com",
		Password: "securepass123",
	}

	output, err := uc.Execute(context.Background(), input)
	assert.ErrorIs(t, err, ErrInvalidInput)
	assert.Nil(t, output)
}

func TestCreateUser_ShouldReturnValidationErrorWhenEmailEmpty(t *testing.T) {
	repo := &domain.MockRepository{}
	hasher := &MockPasswordHasher{}

	uc := NewCreateUserUseCase(repo, hasher)
	input := CreateUserInput{
		Name:     "John Doe",
		Email:    "",
		Password: "securepass123",
	}

	output, err := uc.Execute(context.Background(), input)
	assert.ErrorIs(t, err, ErrInvalidInput)
	assert.Nil(t, output)
}

func TestCreateUser_ShouldReturnErrorWhenRepositoryCreateFails(t *testing.T) {
	repo := &domain.MockRepository{
		FindByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
			return nil, nil
		},
		CreateFunc: func(ctx context.Context, user *domain.User) error {
			return errors.New("db error")
		},
	}
	hasher := &MockPasswordHasher{
		HashFunc: func(password string) (string, error) {
			return "hashed-" + password, nil
		},
	}

	uc := NewCreateUserUseCase(repo, hasher)
	input := CreateUserInput{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "securepass123",
	}

	output, err := uc.Execute(context.Background(), input)
	require.Error(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "db error")
}

func TestCreateUser_ShouldReturnDomainErrorWhenDomainValidationFails(t *testing.T) {
	repo := &domain.MockRepository{
		FindByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
			return nil, nil
		},
	}
	hasher := &MockPasswordHasher{
		HashFunc: func(password string) (string, error) {
			return "", nil
		},
	}

	uc := NewCreateUserUseCase(repo, hasher)
	input := CreateUserInput{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "securepass123",
	}

	output, err := uc.Execute(context.Background(), input)

	assert.Error(t, err)
	assert.Nil(t, output)
}
