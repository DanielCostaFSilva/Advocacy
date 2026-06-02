package user

import (
	"context"
	"errors"

	domain "legalflow/internal/user/domain/user"
)

type CreateUserUseCase struct {
	repo   domain.Repository
	hasher PasswordHasher
}

func NewCreateUserUseCase(repo domain.Repository, hasher PasswordHasher) *CreateUserUseCase {
	return &CreateUserUseCase{repo: repo, hasher: hasher}
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {
	if input.Name == "" || input.Email == "" || len(input.Password) < 8 {
		return nil, ErrInvalidInput
	}

	existing, err := uc.repo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrEmailAlreadyExists
	}

	passwordHash, err := uc.hasher.Hash(input.Password)
	if err != nil {
		return nil, errors.Join(ErrHashPassword, err)
	}

	user, err := domain.NewUser(input.Name, input.Email, passwordHash)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return &CreateUserOutput{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
