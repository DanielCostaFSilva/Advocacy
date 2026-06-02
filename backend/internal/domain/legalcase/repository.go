package legalcase

import (
	"context"

	"github.com/google/uuid"
)

type ListCasesParams struct {
	Offset   int
	Limit    int
	Number   string
	ClientID string
	Status   string
	Sort     string
	Order    string
}

type CaseRepository interface {
	Create(ctx context.Context, legalCase *Case) error
	FindByID(ctx context.Context, id uuid.UUID) (*Case, error)
	FindByNumber(ctx context.Context, number string) (*Case, error)
	List(ctx context.Context, params ListCasesParams) ([]Case, int64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status CaseStatus) error
}
