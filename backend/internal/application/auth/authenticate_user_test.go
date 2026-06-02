package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domain "legalflow/internal/user/domain/user"
)

func TestAuthenticateUser_ShouldSucceedWhenValidCredentials(t *testing.T) {
	user, _ := domain.NewUser("John Doe", "john@example.com", "hashed-pass")

	repo := &domain.MockRepository{
		FindByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
			return user, nil
		},
	}
	verifier := &MockPasswordVerifier{
		CompareFunc: func(plain, hash string) error {
			return nil
		},
	}

	uc := NewAuthenticateUserUseCase(repo, verifier)
	input := AuthenticateUserInput{
		Email:    "john@example.com",
		Password: "correctpassword",
	}

	output, err := uc.Execute(context.Background(), input)
	require.NoError(t, err)
	require.NotNil(t, output)
	assert.Equal(t, user.ID.String(), output.UserID)
	assert.Equal(t, "John Doe", output.Name)
	assert.Equal(t, "john@example.com", output.Email)
}

func TestAuthenticateUser_ShouldReturnErrInvalidCredentialsWhenEmailNotFound(t *testing.T) {
	repo := &domain.MockRepository{
		FindByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
			return nil, nil
		},
	}
	verifier := &MockPasswordVerifier{}

	uc := NewAuthenticateUserUseCase(repo, verifier)
	input := AuthenticateUserInput{
		Email:    "unknown@example.com",
		Password: "anypassword",
	}

	output, err := uc.Execute(context.Background(), input)
	assert.ErrorIs(t, err, ErrInvalidCredentials)
	assert.Nil(t, output)
}

func TestAuthenticateUser_ShouldReturnErrInvalidCredentialsWhenWrongPassword(t *testing.T) {
	user, _ := domain.NewUser("John Doe", "john@example.com", "hashed-pass")

	repo := &domain.MockRepository{
		FindByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
			return user, nil
		},
	}
	verifier := &MockPasswordVerifier{
		CompareFunc: func(plain, hash string) error {
			return errors.New("wrong password")
		},
	}

	uc := NewAuthenticateUserUseCase(repo, verifier)
	input := AuthenticateUserInput{
		Email:    "john@example.com",
		Password: "wrongpassword",
	}

	output, err := uc.Execute(context.Background(), input)
	assert.ErrorIs(t, err, ErrInvalidCredentials)
	assert.Nil(t, output)
}

func TestAuthenticateUser_ShouldReturnErrInvalidInputWhenEmailEmpty(t *testing.T) {
	repo := &domain.MockRepository{}
	verifier := &MockPasswordVerifier{}

	uc := NewAuthenticateUserUseCase(repo, verifier)
	input := AuthenticateUserInput{
		Email:    "",
		Password: "somepassword",
	}

	output, err := uc.Execute(context.Background(), input)
	assert.ErrorIs(t, err, ErrInvalidInput)
	assert.Nil(t, output)
}

func TestAuthenticateUser_ShouldReturnErrInvalidInputWhenPasswordEmpty(t *testing.T) {
	repo := &domain.MockRepository{}
	verifier := &MockPasswordVerifier{}

	uc := NewAuthenticateUserUseCase(repo, verifier)
	input := AuthenticateUserInput{
		Email:    "john@example.com",
		Password: "",
	}

	output, err := uc.Execute(context.Background(), input)
	assert.ErrorIs(t, err, ErrInvalidInput)
	assert.Nil(t, output)
}

func TestAuthenticateUser_ShouldReturnErrorWhenRepositoryFails(t *testing.T) {
	repo := &domain.MockRepository{
		FindByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
			return nil, errors.New("db connection error")
		},
	}
	verifier := &MockPasswordVerifier{}

	uc := NewAuthenticateUserUseCase(repo, verifier)
	input := AuthenticateUserInput{
		Email:    "john@example.com",
		Password: "somepassword",
	}

	output, err := uc.Execute(context.Background(), input)
	require.Error(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "db connection error")
}

func TestAuthenticateUser_ShouldReturnErrInvalidCredentialsWhenVerifierFails(t *testing.T) {
	user, _ := domain.NewUser("John Doe", "john@example.com", "hashed-pass")

	repo := &domain.MockRepository{
		FindByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
			return user, nil
		},
	}
	verifier := &MockPasswordVerifier{
		CompareFunc: func(plain, hash string) error {
			return errors.New("crypto error")
		},
	}

	uc := NewAuthenticateUserUseCase(repo, verifier)
	input := AuthenticateUserInput{
		Email:    "john@example.com",
		Password: "somepassword",
	}

	output, err := uc.Execute(context.Background(), input)
	assert.ErrorIs(t, err, ErrInvalidCredentials)
	assert.Nil(t, output)
}
