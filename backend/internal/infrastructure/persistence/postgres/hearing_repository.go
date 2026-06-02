package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	domain "legalflow/internal/domain/hearing"
)

type HearingRepository struct {
	db *sql.DB
	q  *Queries
}

func NewHearingRepository(db *sql.DB) *HearingRepository {
	return &HearingRepository{
		db: db,
		q:  New(db),
	}
}

func (r *HearingRepository) Create(ctx context.Context, h *domain.Hearing) error {
	params := CreateHearingParams{
		CaseID:      h.CaseID,
		Title:       h.Title,
		Description: sql.NullString{String: h.Description, Valid: h.Description != ""},
		Type:        string(h.Type),
		Location:    h.Location,
		ScheduledAt: h.ScheduledAt,
	}

	created, err := r.q.CreateHearing(ctx, params)
	if err != nil {
		return err
	}

	h.ID = created.ID
	h.CreatedAt = created.CreatedAt
	h.UpdatedAt = created.UpdatedAt
	return nil
}

func (r *HearingRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Hearing, error) {
	result, err := r.q.GetHearingByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return toHearingDomain(result), nil
}

func (r *HearingRepository) ListByCaseID(ctx context.Context, caseID uuid.UUID) ([]domain.Hearing, error) {
	rows, err := r.q.ListHearingsByCaseID(ctx, caseID)
	if err != nil {
		return nil, err
	}

	hearings := make([]domain.Hearing, len(rows))
	for i, row := range rows {
		hearings[i] = *toHearingDomain(row)
	}
	return hearings, nil
}

func toHearingDomain(h Hearing) *domain.Hearing {
	return &domain.Hearing{
		ID:          h.ID,
		CaseID:      h.CaseID,
		Title:       h.Title,
		Description: h.Description.String,
		Type:        domain.HearingType(h.Type),
		Location:    h.Location,
		ScheduledAt: h.ScheduledAt,
		CreatedAt:   h.CreatedAt,
		UpdatedAt:   h.UpdatedAt,
	}
}
