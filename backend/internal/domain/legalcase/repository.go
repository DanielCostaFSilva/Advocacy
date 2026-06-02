package legalcase

import (
	"context"

	"github.com/google/uuid"
)

type CaseRepository interface {
	Create(ctx context.Context, legalCase *Case) error
	FindByID(ctx context.Context, id uuid.UUID) (*Case, error)
	FindByNumber(ctx context.Context, number string) (*Case, error)
}
