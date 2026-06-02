package client

import "time"

type GetClientResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CPF       string    `json:"cpf"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateClientResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	CPF   string `json:"cpf"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type ClientItem struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	CPF   string `json:"cpf"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type PaginationInfo struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type ListClientsResponse struct {
	Data       []ClientItem   `json:"data"`
	Pagination PaginationInfo `json:"pagination"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
