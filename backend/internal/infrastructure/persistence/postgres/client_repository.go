package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	domain "legalflow/internal/domain/client"
)

type ClientRepository struct {
	q *Queries
}

func NewClientRepository(db *sql.DB) *ClientRepository {
	return &ClientRepository{
		q: New(db),
	}
}

func (r *ClientRepository) Create(ctx context.Context, client *domain.Client) error {
	params := CreateClientParams{
		Name:  client.Name,
		Cpf:   client.CPF,
		Email: client.Email,
		Phone: client.Phone,
	}

	created, err := r.q.CreateClient(ctx, params)
	if err != nil {
		return err
	}

	client.ID = created.ID
	client.CreatedAt = created.CreatedAt
	client.UpdatedAt = created.UpdatedAt

	return nil
}

func (r *ClientRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Client, error) {
	result, err := r.q.GetClientByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return toClientDomain(result), nil
}

func (r *ClientRepository) FindByCPF(ctx context.Context, cpf string) (*domain.Client, error) {
	result, err := r.q.GetClientByCPF(ctx, cpf)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return toClientDomain(result), nil
}

func toClientDomain(c Client) *domain.Client {
	return &domain.Client{
		ID:        c.ID,
		Name:      c.Name,
		CPF:       c.Cpf,
		Email:     c.Email,
		Phone:     c.Phone,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
