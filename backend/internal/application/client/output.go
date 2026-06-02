package client

import "time"

type CreateClientOutput struct {
	ID    string
	Name  string
	CPF   string
	Email string
	Phone string
}

type UpdateClientOutput struct {
	ID    string
	Name  string
	CPF   string
	Email string
	Phone string
}

type GetClientByIDOutput struct {
	ID        string
	Name      string
	CPF       string
	Email     string
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ClientDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	CPF   string `json:"cpf"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type ListClientsOutput struct {
	Data       []ClientDTO
	Page       int
	PageSize   int
	Total      int64
	TotalPages int
}
