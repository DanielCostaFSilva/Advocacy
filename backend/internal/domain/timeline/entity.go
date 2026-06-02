package timeline

import (
	"time"

	"github.com/google/uuid"
)

type TimelineEvent struct {
	ID          uuid.UUID
	CaseID      uuid.UUID
	Type        EventType
	Description string
	Metadata    []byte
	CreatedAt   time.Time
}

func NewEvent(caseID uuid.UUID, eventType EventType, description string) *TimelineEvent {
	return &TimelineEvent{
		ID:          uuid.New(),
		CaseID:      caseID,
		Type:        eventType,
		Description: description,
		CreatedAt:   time.Now(),
	}
}
