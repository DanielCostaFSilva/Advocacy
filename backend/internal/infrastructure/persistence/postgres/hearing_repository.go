package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

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

func (r *HearingRepository) List(ctx context.Context, params domain.ListHearingsParams) ([]domain.Hearing, int64, error) {
	var conditions []string
	var args []any
	argIdx := 1

	if params.CaseID != "" {
		conditions = append(conditions, fmt.Sprintf("case_id = $%d", argIdx))
		args = append(args, params.CaseID)
		argIdx++
	}
	if params.Type != "" {
		conditions = append(conditions, fmt.Sprintf("type = $%d", argIdx))
		args = append(args, params.Type)
		argIdx++
	}
	if params.StartDate != nil {
		conditions = append(conditions, fmt.Sprintf("scheduled_at >= $%d", argIdx))
		args = append(args, *params.StartDate)
		argIdx++
	}
	if params.EndDate != nil {
		conditions = append(conditions, fmt.Sprintf("scheduled_at <= $%d", argIdx))
		args = append(args, *params.EndDate)
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	sort := "scheduled_at"
	switch params.Sort {
	case "created_at":
		sort = "created_at"
	case "title":
		sort = "title"
	case "scheduled_at":
		sort = "scheduled_at"
	}

	order := "asc"
	if params.Order == "desc" {
		order = "desc"
	}

	query := fmt.Sprintf(
		"SELECT id, case_id, title, description, type, location, scheduled_at, created_at, updated_at FROM hearings %s ORDER BY %s %s LIMIT $%d OFFSET $%d",
		whereClause, sort, order, argIdx, argIdx+1,
	)
	args = append(args, params.Limit, params.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var hearings []domain.Hearing
	for rows.Next() {
		var h Hearing
		if err := rows.Scan(&h.ID, &h.CaseID, &h.Title, &h.Description, &h.Type, &h.Location, &h.ScheduledAt, &h.CreatedAt, &h.UpdatedAt); err != nil {
			return nil, 0, err
		}
		hearings = append(hearings, *toHearingDomain(h))
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	var total int64
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM hearings %s", whereClause)
	if err := r.db.QueryRowContext(ctx, countQuery, args[:len(args)-2]...).Scan(&total); err != nil {
		return nil, 0, err
	}

	return hearings, total, nil
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
