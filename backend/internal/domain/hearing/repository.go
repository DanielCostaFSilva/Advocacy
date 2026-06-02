package hearing

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, hearing *Hearing) error
	FindByID(ctx context.Context, id uuid.UUID) (*Hearing, error)
	ListByCaseID(ctx context.Context, caseID uuid.UUID) ([]Hearing, error)
}
