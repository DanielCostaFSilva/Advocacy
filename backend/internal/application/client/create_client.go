package client

import (
	"context"

	domain "legalflow/internal/domain/client"
)

type CreateClientUseCase struct {
	repo domain.ClientRepository
}

func NewCreateClientUseCase(repo domain.ClientRepository) *CreateClientUseCase {
	return &CreateClientUseCase{repo: repo}
}

func (uc *CreateClientUseCase) Execute(ctx context.Context, input CreateClientInput) (*CreateClientOutput, error) {
	if input.Name == "" || input.CPF == "" {
		return nil, ErrInvalidInput
	}

	existing, err := uc.repo.FindByCPF(ctx, input.CPF)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrClientAlreadyExists
	}

	client, err := domain.NewClient(input.Name, input.CPF, input.Email, input.Phone)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, client); err != nil {
		return nil, err
	}

	return &CreateClientOutput{
		ID:    client.ID.String(),
		Name:  client.Name,
		CPF:   client.CPF,
		Email: client.Email,
		Phone: client.Phone,
	}, nil
}
