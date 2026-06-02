package auth

import (
	"context"

	domain "legalflow/internal/user/domain/user"
)

type PasswordVerifier interface {
	Compare(plainPassword, hash string) error
}

type AuthenticateUserUseCase struct {
	repo     domain.Repository
	verifier PasswordVerifier
}

func NewAuthenticateUserUseCase(repo domain.Repository, verifier PasswordVerifier) *AuthenticateUserUseCase {
	return &AuthenticateUserUseCase{repo: repo, verifier: verifier}
}

func (uc *AuthenticateUserUseCase) Execute(ctx context.Context, input AuthenticateUserInput) (*AuthenticateUserOutput, error) {
	if input.Email == "" || input.Password == "" {
		return nil, ErrInvalidInput
	}

	user, err := uc.repo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if err := uc.verifier.Compare(input.Password, user.PasswordHash); err != nil {
		return nil, ErrInvalidCredentials
	}

	return &AuthenticateUserOutput{
		UserID: user.ID.String(),
		Name:   user.Name,
		Email:  user.Email,
	}, nil
}
