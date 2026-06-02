package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domain "legalflow/internal/user/domain/user"
)

func TestGetCurrentUser_ShouldReturnUserWhenExists(t *testing.T) {
	userID := uuid.New()
	repo := &domain.MockRepository{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.User, error) {
			return &domain.User{
				ID:    userID,
				Name:  "John Doe",
				Email: "john@example.com",
			}, nil
		},
	}

	uc := NewGetCurrentUserUseCase(repo)
	output, err := uc.Execute(context.Background(), GetCurrentUserInput{
		UserID: userID.String(),
	})

	require.NoError(t, err)
	require.NotNil(t, output)
	assert.Equal(t, userID.String(), output.ID)
	assert.Equal(t, "John Doe", output.Name)
	assert.Equal(t, "john@example.com", output.Email)
}

func TestGetCurrentUser_ShouldReturnErrUserNotFoundWhenNotExists(t *testing.T) {
	repo := &domain.MockRepository{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.User, error) {
			return nil, nil
		},
	}

	uc := NewGetCurrentUserUseCase(repo)
	output, err := uc.Execute(context.Background(), GetCurrentUserInput{
		UserID: uuid.New().String(),
	})

	assert.ErrorIs(t, err, ErrUserNotFound)
	assert.Nil(t, output)
}

func TestGetCurrentUser_ShouldReturnErrorWhenRepositoryFails(t *testing.T) {
	repo := &domain.MockRepository{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.User, error) {
			return nil, errors.New("db connection failed")
		},
	}

	uc := NewGetCurrentUserUseCase(repo)
	output, err := uc.Execute(context.Background(), GetCurrentUserInput{
		UserID: uuid.New().String(),
	})

	require.Error(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "db connection failed")
}

func TestGetCurrentUser_ShouldReturnErrInvalidInputWhenUserIDEmpty(t *testing.T) {
	repo := &domain.MockRepository{}

	uc := NewGetCurrentUserUseCase(repo)
	output, err := uc.Execute(context.Background(), GetCurrentUserInput{
		UserID: "",
	})

	assert.ErrorIs(t, err, ErrInvalidInput)
	assert.Nil(t, output)
}

func TestGetCurrentUser_ShouldReturnErrInvalidInputWhenUserIDNotUUID(t *testing.T) {
	repo := &domain.MockRepository{}

	uc := NewGetCurrentUserUseCase(repo)
	output, err := uc.Execute(context.Background(), GetCurrentUserInput{
		UserID: "not-a-uuid",
	})

	assert.ErrorIs(t, err, ErrInvalidInput)
	assert.Nil(t, output)
}
