package client

import (
	"context"

	"github.com/google/uuid"
	domain "legalflow/internal/domain/client"
)

type GetClientByIDUseCase struct {
	repo domain.ClientRepository
}

func NewGetClientByIDUseCase(repo domain.ClientRepository) *GetClientByIDUseCase {
	return &GetClientByIDUseCase{repo: repo}
}

func (uc *GetClientByIDUseCase) Execute(ctx context.Context, input GetClientByIDInput) (*GetClientByIDOutput, error) {
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

	return &GetClientByIDOutput{
		ID:        client.ID.String(),
		Name:      client.Name,
		CPF:       client.CPF,
		Email:     client.Email,
		Phone:     client.Phone,
		CreatedAt: client.CreatedAt,
		UpdatedAt: client.UpdatedAt,
	}, nil
}
