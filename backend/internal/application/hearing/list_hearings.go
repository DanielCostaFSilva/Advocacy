package hearing

import "time"

type ListHearingsInput struct {
	Page      int
	PageSize  int
	CaseID    string
	Type      string
	StartDate string
	EndDate   string
	Sort      string
	Order     string
}

type HearingDTO struct {
	ID          string    `json:"id"`
	CaseID      string    `json:"case_id"`
	Title       string    `json:"title"`
	Type        string    `json:"type"`
	Location    string    `json:"location"`
	ScheduledAt time.Time `json:"scheduled_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type ListHearingsOutput struct {
	Data       []HearingDTO
	Page       int
	PageSize   int
	Total      int64
	TotalPages int
}
