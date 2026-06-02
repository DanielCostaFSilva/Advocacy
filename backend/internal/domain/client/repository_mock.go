package client

import (
	"context"

	"github.com/google/uuid"
)

type MockClientRepository struct {
	CreateFunc    func(ctx context.Context, client *Client) error
	FindByIDFunc  func(ctx context.Context, id uuid.UUID) (*Client, error)
	FindByCPFFunc func(ctx context.Context, cpf string) (*Client, error)
	UpdateFunc    func(ctx context.Context, client *Client) error
	ListFunc      func(ctx context.Context, params ListClientsParams) ([]Client, int64, error)
}

func (m *MockClientRepository) Create(ctx context.Context, client *Client) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, client)
	}
	return nil
}

func (m *MockClientRepository) FindByID(ctx context.Context, id uuid.UUID) (*Client, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockClientRepository) FindByCPF(ctx context.Context, cpf string) (*Client, error) {
	if m.FindByCPFFunc != nil {
		return m.FindByCPFFunc(ctx, cpf)
	}
	return nil, nil
}

func (m *MockClientRepository) Update(ctx context.Context, client *Client) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, client)
	}
	return nil
}

func (m *MockClientRepository) List(ctx context.Context, params ListClientsParams) ([]Client, int64, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, params)
	}
	return nil, 0, nil
}
