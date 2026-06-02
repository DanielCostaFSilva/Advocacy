package client

import (
	"context"

	"github.com/google/uuid"
)

type ClientRepository interface {
	Create(ctx context.Context, client *Client) error
	FindByID(ctx context.Context, id uuid.UUID) (*Client, error)
	FindByCPF(ctx context.Context, cpf string) (*Client, error)
}
