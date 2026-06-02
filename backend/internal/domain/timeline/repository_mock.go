package timeline

import (
	"context"

	"github.com/google/uuid"
)

type MockRepository struct {
	FindByCaseIDFunc func(ctx context.Context, caseID uuid.UUID) ([]TimelineEvent, error)
	CreateFunc       func(ctx context.Context, event *TimelineEvent) error
}

func (m *MockRepository) FindByCaseID(ctx context.Context, caseID uuid.UUID) ([]TimelineEvent, error) {
	if m.FindByCaseIDFunc != nil {
		return m.FindByCaseIDFunc(ctx, caseID)
	}
	return nil, nil
}

func (m *MockRepository) Create(ctx context.Context, event *TimelineEvent) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, event)
	}
	return nil
}
