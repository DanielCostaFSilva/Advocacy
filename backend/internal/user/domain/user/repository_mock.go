package user

import (
	"context"

	"github.com/google/uuid"
)

type MockRepository struct {
	CreateFunc    func(ctx context.Context, user *User) error
	FindByIDFunc  func(ctx context.Context, id uuid.UUID) (*User, error)
	FindByEmailFunc func(ctx context.Context, email string) (*User, error)
}

func (m *MockRepository) Create(ctx context.Context, user *User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, user)
	}
	return nil
}

func (m *MockRepository) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	if m.FindByEmailFunc != nil {
		return m.FindByEmailFunc(ctx, email)
	}
	return nil, nil
}
