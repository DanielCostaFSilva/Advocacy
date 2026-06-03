package timeline

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type TimelineEvent struct {
	ID          uuid.UUID
	CaseID      uuid.UUID
	Type        EventType
	Description string
	Metadata    json.RawMessage
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

func NewEventWithMetadata(caseID uuid.UUID, eventType EventType, description string, metadata json.RawMessage) *TimelineEvent {
	event := NewEvent(caseID, eventType, description)
	event.Metadata = metadata
	return event
}
