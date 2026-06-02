package client

type CreateClientRequest struct {
	Name  string `json:"name"`
	CPF   string `json:"cpf"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}
