package hearing

import (
	"context"

	"github.com/google/uuid"
)

type MockRepository struct {
	CreateFunc      func(ctx context.Context, hearing *Hearing) error
	FindByIDFunc    func(ctx context.Context, id uuid.UUID) (*Hearing, error)
	ListByCaseIDFunc func(ctx context.Context, caseID uuid.UUID) ([]Hearing, error)
	ListFunc        func(ctx context.Context, params ListHearingsParams) ([]Hearing, int64, error)
}

func (m *MockRepository) Create(ctx context.Context, hearing *Hearing) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, hearing)
	}
	return nil
}

func (m *MockRepository) FindByID(ctx context.Context, id uuid.UUID) (*Hearing, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockRepository) ListByCaseID(ctx context.Context, caseID uuid.UUID) ([]Hearing, error) {
	if m.ListByCaseIDFunc != nil {
		return m.ListByCaseIDFunc(ctx, caseID)
	}
	return nil, nil
}

func (m *MockRepository) List(ctx context.Context, params ListHearingsParams) ([]Hearing, int64, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, params)
	}
	return nil, 0, nil
}
