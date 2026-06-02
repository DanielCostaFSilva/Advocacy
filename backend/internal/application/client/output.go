package client

type CreateClientOutput struct {
	ID    string
	Name  string
	CPF   string
	Email string
	Phone string
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
