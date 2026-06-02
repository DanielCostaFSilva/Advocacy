package timeline

import "time"

type EventItem struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type TimelineResponse struct {
	Events []EventItem `json:"events"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
