package timeline

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	FindByCaseID(ctx context.Context, caseID uuid.UUID) ([]TimelineEvent, error)
	Create(ctx context.Context, event *TimelineEvent) error
}
