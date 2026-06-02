package document

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, document *Document) error
	FindByID(ctx context.Context, id uuid.UUID) (*Document, error)
	ListByCaseID(ctx context.Context, caseID uuid.UUID) ([]Document, error)
}
