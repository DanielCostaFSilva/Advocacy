package legalcase

import "time"

type CaseDTO struct {
	ID        string    `json:"id"`
	ClientID  string    `json:"client_id"`
	Number    string    `json:"number"`
	Title     string    `json:"title"`
	Court     string    `json:"court"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type ListCasesOutput struct {
	Data       []CaseDTO
	Page       int
	PageSize   int
	Total      int64
	TotalPages int
}

type UpdateCaseStatusOutput struct {
	ID     string
	Status string
}

type GetCaseByIDOutput struct {
	ID          string
	ClientID    string
	Number      string
	Title       string
	Description string
	Court       string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateCaseOutput struct {
	ID          string
	ClientID    string
	Number      string
	Title       string
	Description string
	Court       string
	Status      string
}
