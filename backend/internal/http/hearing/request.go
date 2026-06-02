package hearing

type CreateHearingRequest struct {
	CaseID      string `json:"case_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Location    string `json:"location"`
	ScheduledAt string `json:"scheduled_at"`
}
