package timeline

import "time"

type TimelineEventDTO struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type GetCaseTimelineOutput struct {
	Events []TimelineEventDTO
}
