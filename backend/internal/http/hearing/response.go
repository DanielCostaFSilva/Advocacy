package hearing

import "time"

type CreateHearingResponse struct {
	ID          string    `json:"id"`
	CaseID      string    `json:"case_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	Location    string    `json:"location"`
	ScheduledAt time.Time `json:"scheduled_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
