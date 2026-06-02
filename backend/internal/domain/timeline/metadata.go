package timeline

import "time"

type HearingMetadata struct {
	HearingID   string    `json:"hearing_id"`
	Type        string    `json:"hearing_type"`
	ScheduledAt time.Time `json:"scheduled_at"`
}
