package timeline

import (
	"encoding/json"
	"time"
)

type TimelineEventDTO struct {
	ID          string           `json:"id"`
	Type        string           `json:"type"`
	Description string           `json:"description"`
	Metadata    *json.RawMessage `json:"metadata,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
}

type GetCaseTimelineOutput struct {
	Events []TimelineEventDTO
}
