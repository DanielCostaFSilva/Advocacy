package contract

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, contract *Contract) error
	FindByID(ctx context.Context, id uuid.UUID) (*Contract, error)
	ListByCaseID(ctx context.Context, caseID uuid.UUID) ([]Contract, error)
}
