package client

import (
	"context"

	"github.com/google/uuid"
	domain "legalflow/internal/domain/client"
)

type DeleteClientUseCase struct {
	repo domain.ClientRepository
}

func NewDeleteClientUseCase(repo domain.ClientRepository) *DeleteClientUseCase {
	return &DeleteClientUseCase{repo: repo}
}

func (uc *DeleteClientUseCase) Execute(ctx context.Context, input DeleteClientInput) error {
	id, err := uuid.Parse(input.ID)
	if err != nil {
		return ErrInvalidClientID
	}

	client, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if client == nil {
		return ErrClientNotFound
	}

	return uc.repo.SoftDelete(ctx, id)
}
