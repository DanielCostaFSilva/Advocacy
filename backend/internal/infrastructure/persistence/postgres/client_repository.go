package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

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

func (r *ClientRepository) Update(ctx context.Context, client *domain.Client) error {
	params := UpdateClientParams{
		ID:    client.ID,
		Name:  client.Name,
		Email: client.Email,
		Phone: client.Phone,
	}

	updated, err := r.q.UpdateClient(ctx, params)
	if err != nil {
		return err
	}

	client.UpdatedAt = updated.UpdatedAt
	return nil
}

func (r *ClientRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.q.SoftDeleteClient(ctx, id)
}

func (r *ClientRepository) List(ctx context.Context, params domain.ListClientsParams) ([]domain.Client, int64, error) {
	var conditions []string
	var args []any
	argIdx := 1

	if params.Name != "" {
		conditions = append(conditions, fmt.Sprintf("name ILIKE '%%' || $%d || '%%'", argIdx))
		args = append(args, params.Name)
		argIdx++
	}
	if params.CPF != "" {
		conditions = append(conditions, fmt.Sprintf("cpf = $%d", argIdx))
		args = append(args, params.CPF)
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ") + " AND deleted_at IS NULL"
	} else {
		whereClause = "WHERE deleted_at IS NULL"
	}

	sort := "name"
	if params.Sort == "created_at" {
		sort = "created_at"
	}

	order := "asc"
	if params.Order == "desc" {
		order = "desc"
	}

	query := fmt.Sprintf("SELECT id, name, cpf, email, phone, created_at, updated_at, deleted_at FROM clients %s ORDER BY %s %s LIMIT $%d OFFSET $%d",
		whereClause, sort, order, argIdx, argIdx+1)
	args = append(args, params.Limit, params.Offset)

	rows, err := r.q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var clients []domain.Client
	for rows.Next() {
		var c Client
		if err := rows.Scan(&c.ID, &c.Name, &c.Cpf, &c.Email, &c.Phone, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt); err != nil {
			return nil, 0, err
		}
		clients = append(clients, *toClientDomain(c))
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	var total int64
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM clients %s", whereClause)
	if err := r.q.db.QueryRowContext(ctx, countQuery, args[:len(args)-2]...).Scan(&total); err != nil {
		return nil, 0, err
	}

	return clients, total, nil
}

func toClientDomain(c Client) *domain.Client {
	var deletedAt *time.Time
	if c.DeletedAt.Valid {
		deletedAt = &c.DeletedAt.Time
	}
	return &domain.Client{
		ID:        c.ID,
		Name:      c.Name,
		CPF:       c.Cpf,
		Email:     c.Email,
		Phone:     c.Phone,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		DeletedAt: deletedAt,
	}
}
