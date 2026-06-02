package legalcase

import (
	"context"

	"github.com/google/uuid"
)

type MockCaseRepository struct {
	CreateFunc       func(ctx context.Context, legalCase *Case) error
	FindByIDFunc     func(ctx context.Context, id uuid.UUID) (*Case, error)
	FindByNumberFunc func(ctx context.Context, number string) (*Case, error)
}

func (m *MockCaseRepository) Create(ctx context.Context, legalCase *Case) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, legalCase)
	}
	return nil
}

func (m *MockCaseRepository) FindByID(ctx context.Context, id uuid.UUID) (*Case, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockCaseRepository) FindByNumber(ctx context.Context, number string) (*Case, error) {
	if m.FindByNumberFunc != nil {
		return m.FindByNumberFunc(ctx, number)
	}
	return nil, nil
}
