package legalcase

type CreateCaseResponse struct {
	ID          string `json:"id"`
	ClientID    string `json:"client_id"`
	Number      string `json:"number"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Court       string `json:"court"`
	Status      string `json:"status"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
