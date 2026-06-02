package postgres

import (
	"context"
	"database/sql"
	"errors"

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
