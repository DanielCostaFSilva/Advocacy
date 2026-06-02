package auth

import (
	"context"

	"github.com/google/uuid"
	domain "legalflow/internal/user/domain/user"
)

type UserRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
}

type GetCurrentUserUseCase struct {
	repo UserRepository
}

func NewGetCurrentUserUseCase(repo UserRepository) *GetCurrentUserUseCase {
	return &GetCurrentUserUseCase{repo: repo}
}

func (uc *GetCurrentUserUseCase) Execute(ctx context.Context, input GetCurrentUserInput) (*GetCurrentUserOutput, error) {
	if input.UserID == "" {
		return nil, ErrInvalidInput
	}

	id, err := uuid.Parse(input.UserID)
	if err != nil {
		return nil, ErrInvalidInput
	}

	user, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return &GetCurrentUserOutput{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
