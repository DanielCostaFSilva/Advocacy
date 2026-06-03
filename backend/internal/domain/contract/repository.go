package contract

import (
	"context"

	"github.com/google/uuid"
)

type ListContractsParams struct {
	Offset   int
	Limit    int
	ClientID string
	CaseID   string
	Type     string
	Active   *bool
	Sort     string
	Order    string
}

type Repository interface {
	Create(ctx context.Context, contract *Contract) error
	FindByID(ctx context.Context, id uuid.UUID) (*Contract, error)
	ListByCaseID(ctx context.Context, caseID uuid.UUID) ([]Contract, error)
	List(ctx context.Context, params ListContractsParams) ([]Contract, int64, error)
	Update(ctx context.Context, contract *Contract) error
}
