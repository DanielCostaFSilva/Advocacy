package hearing

import "time"

type CreateHearingInput struct {
	CaseID      string
	Title       string
	Description string
	Type        string
	Location    string
	ScheduledAt time.Time
}
