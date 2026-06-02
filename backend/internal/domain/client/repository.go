package client

import (
	"context"

	"github.com/google/uuid"
)

type ListClientsParams struct {
	Offset int
	Limit  int
	Name   string
	CPF    string
	Sort   string
	Order  string
}

type ClientRepository interface {
	Create(ctx context.Context, client *Client) error
	FindByID(ctx context.Context, id uuid.UUID) (*Client, error)
	FindByCPF(ctx context.Context, cpf string) (*Client, error)
	Update(ctx context.Context, client *Client) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, params ListClientsParams) ([]Client, int64, error)
}
