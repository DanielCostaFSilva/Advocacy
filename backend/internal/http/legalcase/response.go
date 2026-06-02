package legalcase

import "time"

type CreateCaseResponse struct {
	ID          string `json:"id"`
	ClientID    string `json:"client_id"`
	Number      string `json:"number"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Court       string `json:"court"`
	Status      string `json:"status"`
}

type CaseItem struct {
	ID        string    `json:"id"`
	ClientID  string    `json:"client_id"`
	Number    string    `json:"number"`
	Title     string    `json:"title"`
	Court     string    `json:"court"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type PaginationInfo struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type GetCaseResponse struct {
	ID          string    `json:"id"`
	ClientID    string    `json:"client_id"`
	Number      string    `json:"number"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Court       string    `json:"court"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ListCasesResponse struct {
	Data       []CaseItem     `json:"data"`
	Pagination PaginationInfo `json:"pagination"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
