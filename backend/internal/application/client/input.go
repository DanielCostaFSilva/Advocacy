package client

type CreateClientInput struct {
	Name  string
	CPF   string
	Email string
	Phone string
}

type GetClientByIDInput struct {
	ID string
}

type ListClientsInput struct {
	Page     int
	PageSize int
	Name     string
	CPF      string
	Sort     string
	Order    string
}
