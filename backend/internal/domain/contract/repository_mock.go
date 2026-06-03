package contract

import (
	"context"

	"github.com/google/uuid"
)

type MockRepository struct {
	CreateFunc    func(ctx context.Context, contract *Contract) error
	FindByIDFunc  func(ctx context.Context, id uuid.UUID) (*Contract, error)
	ListByCaseIDFunc func(ctx context.Context, caseID uuid.UUID) ([]Contract, error)
}

func (m *MockRepository) Create(ctx context.Context, contract *Contract) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, contract)
	}
	return nil
}

func (m *MockRepository) FindByID(ctx context.Context, id uuid.UUID) (*Contract, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockRepository) ListByCaseID(ctx context.Context, caseID uuid.UUID) ([]Contract, error) {
	if m.ListByCaseIDFunc != nil {
		return m.ListByCaseIDFunc(ctx, caseID)
	}
	return nil, nil
}
