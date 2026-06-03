package document

import (
	"context"

	"github.com/google/uuid"
)

type ListDocumentsParams struct {
	CaseID string
	Offset int
	Limit  int
	Type   string
	Sort   string
	Order  string
}

type Repository interface {
	Create(ctx context.Context, document *Document) error
	FindByID(ctx context.Context, id uuid.UUID) (*Document, error)
	ListByCaseID(ctx context.Context, caseID uuid.UUID) ([]Document, error)
	ListByCaseIDPaginated(ctx context.Context, params ListDocumentsParams) ([]Document, int64, error)
}
