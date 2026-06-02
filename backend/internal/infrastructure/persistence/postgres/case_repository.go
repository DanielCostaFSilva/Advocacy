package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	domain "legalflow/internal/domain/legalcase"
)

type CaseRepository struct {
	q *Queries
}

func NewCaseRepository(db *sql.DB) *CaseRepository {
	return &CaseRepository{
		q: New(db),
	}
}

func (r *CaseRepository) Create(ctx context.Context, c *domain.Case) error {
	params := CreateCaseParams{
		ClientID:    c.ClientID,
		Number:      c.Number,
		Title:       c.Title,
		Description: c.Description,
		Court:       c.Court,
		Status:      string(c.Status),
	}

	created, err := r.q.CreateCase(ctx, params)
	if err != nil {
		return err
	}

	c.ID = created.ID
	c.CreatedAt = created.CreatedAt
	c.UpdatedAt = created.UpdatedAt
	return nil
}

func (r *CaseRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Case, error) {
	result, err := r.q.GetCaseByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return toCaseDomain(result), nil
}

func (r *CaseRepository) List(ctx context.Context, params domain.ListCasesParams) ([]domain.Case, int64, error) {
	var conditions []string
	var args []any
	argIdx := 1

	if params.Number != "" {
		conditions = append(conditions, fmt.Sprintf("number ILIKE '%%' || $%d || '%%'", argIdx))
		args = append(args, params.Number)
		argIdx++
	}
	if params.ClientID != "" {
		conditions = append(conditions, fmt.Sprintf("client_id = $%d", argIdx))
		args = append(args, params.ClientID)
		argIdx++
	}
	if params.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, params.Status)
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	sort := "created_at"
	switch params.Sort {
	case "title":
		sort = "title"
	case "number":
		sort = "number"
	}

	order := "asc"
	if params.Order == "desc" {
		order = "desc"
	}

	query := fmt.Sprintf("SELECT id, client_id, number, title, description, court, status, created_at, updated_at FROM cases %s ORDER BY %s %s LIMIT $%d OFFSET $%d",
		whereClause, sort, order, argIdx, argIdx+1)
	args = append(args, params.Limit, params.Offset)

	rows, err := r.q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var cases []domain.Case
	for rows.Next() {
		var c Case
		if err := rows.Scan(&c.ID, &c.ClientID, &c.Number, &c.Title, &c.Description, &c.Court, &c.Status, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, 0, err
		}
		cases = append(cases, *toCaseDomain(c))
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	var total int64
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM cases %s", whereClause)
	if err := r.q.db.QueryRowContext(ctx, countQuery, args[:len(args)-2]...).Scan(&total); err != nil {
		return nil, 0, err
	}

	return cases, total, nil
}

func (r *CaseRepository) FindByNumber(ctx context.Context, number string) (*domain.Case, error) {
	result, err := r.q.GetCaseByNumber(ctx, number)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return toCaseDomain(result), nil
}

func (r *CaseRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.CaseStatus) error {
	return r.q.UpdateCaseStatus(ctx, UpdateCaseStatusParams{
		ID:     id,
		Status: string(status),
	})
}

func toCaseDomain(c Case) *domain.Case {
	return &domain.Case{
		ID:          c.ID,
		ClientID:    c.ClientID,
		Number:      c.Number,
		Title:       c.Title,
		Description: c.Description,
		Court:       c.Court,
		Status:      domain.CaseStatus(c.Status),
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}
