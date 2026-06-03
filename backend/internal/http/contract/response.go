package contract

type CreateContractResponse struct {
	ID          string  `json:"id"`
	ClientID    string  `json:"client_id"`
	CaseID      string  `json:"case_id"`
	Title       string  `json:"title"`
	Description string  `json:"description,omitempty"`
	Type        string  `json:"type"`
	Amount      string  `json:"amount"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
	Active      bool    `json:"active"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
