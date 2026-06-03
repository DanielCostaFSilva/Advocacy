package hearing

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ListHearingsParams struct {
	Offset    int
	Limit     int
	CaseID    string
	Type      string
	StartDate *time.Time
	EndDate   *time.Time
	Sort      string
	Order     string
}

type Repository interface {
	Create(ctx context.Context, hearing *Hearing) error
	FindByID(ctx context.Context, id uuid.UUID) (*Hearing, error)
	ListByCaseID(ctx context.Context, caseID uuid.UUID) ([]Hearing, error)
	List(ctx context.Context, params ListHearingsParams) ([]Hearing, int64, error)
}
