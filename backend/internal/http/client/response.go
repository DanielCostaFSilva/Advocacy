package client

type CreateClientResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	CPF   string `json:"cpf"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
