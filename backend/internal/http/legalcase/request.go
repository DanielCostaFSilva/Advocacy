package legalcase

type CreateCaseRequest struct {
	ClientID    string `json:"client_id"`
	Number      string `json:"number"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Court       string `json:"court"`
}
