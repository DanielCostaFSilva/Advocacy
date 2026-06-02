package hearing

import "time"

type CreateHearingOutput struct {
	ID          string
	CaseID      string
	Title       string
	Description string
	Type        string
	Location    string
	ScheduledAt time.Time
}
