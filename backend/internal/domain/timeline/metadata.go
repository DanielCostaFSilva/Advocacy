package timeline

import "time"

type HearingMetadata struct {
	HearingID   string    `json:"hearing_id"`
	Type        string    `json:"hearing_type"`
	ScheduledAt time.Time `json:"scheduled_at"`
}

type DocumentMetadata struct {
	DocumentID   string `json:"document_id"`
	DocumentType string `json:"document_type"`
	FileName     string `json:"file_name"`
}
