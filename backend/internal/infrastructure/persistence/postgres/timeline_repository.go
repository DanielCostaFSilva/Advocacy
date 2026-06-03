package postgres

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	domain "legalflow/internal/domain/timeline"
)

type TimelineRepository struct {
	db *sql.DB
	q  *Queries
}

func NewTimelineRepository(db *sql.DB) *TimelineRepository {
	return &TimelineRepository{
		db: db,
		q:  New(db),
	}
}

func (r *TimelineRepository) FindByCaseID(ctx context.Context, caseID uuid.UUID) ([]domain.TimelineEvent, error) {
	rows, err := r.q.GetTimelineByCaseID(ctx, caseID)
	if err != nil {
		return nil, err
	}

	events := make([]domain.TimelineEvent, len(rows))
	for i, row := range rows {
			var metadata json.RawMessage
		if row.Metadata.Valid {
			metadata = make(json.RawMessage, len(row.Metadata.RawMessage))
			copy(metadata, row.Metadata.RawMessage)
		}
		events[i] = domain.TimelineEvent{
			ID:          row.ID,
			CaseID:      row.CaseID,
			Type:        domain.EventType(row.Type),
			Description: row.Description,
			Metadata:    metadata,
			CreatedAt:   row.CreatedAt,
		}
	}
	return events, nil
}

func (r *TimelineRepository) Create(ctx context.Context, event *domain.TimelineEvent) error {
	var metadata pqtype.NullRawMessage
	if len(event.Metadata) > 0 {
		metadata = pqtype.NullRawMessage{RawMessage: []byte(event.Metadata), Valid: true}
	}

	created, err := r.q.CreateTimelineEvent(ctx, CreateTimelineEventParams{
		CaseID:      event.CaseID,
		Type:        string(event.Type),
		Description: event.Description,
		Metadata:    metadata,
	})
	if err != nil {
		return err
	}

	event.ID = created.ID
	event.CreatedAt = created.CreatedAt
	return nil
}
