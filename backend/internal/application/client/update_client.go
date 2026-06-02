package client

import (
	"context"

	"github.com/google/uuid"
	domain "legalflow/internal/domain/client"
)

type UpdateClientUseCase struct {
	repo domain.ClientRepository
}

func NewUpdateClientUseCase(repo domain.ClientRepository) *UpdateClientUseCase {
	return &UpdateClientUseCase{repo: repo}
}

func (uc *UpdateClientUseCase) Execute(ctx context.Context, input UpdateClientInput) (*UpdateClientOutput, error) {
	id, err := uuid.Parse(input.ID)
	if err != nil {
		return nil, ErrInvalidClientID
	}

	client, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, ErrClientNotFound
	}

	if err := client.Update(input.Name, input.Email, input.Phone); err != nil {
		return nil, err
	}

	if err := uc.repo.Update(ctx, client); err != nil {
		return nil, err
	}

	return &UpdateClientOutput{
		ID:    client.ID.String(),
		Name:  client.Name,
		CPF:   client.CPF,
		Email: client.Email,
		Phone: client.Phone,
	}, nil
}
