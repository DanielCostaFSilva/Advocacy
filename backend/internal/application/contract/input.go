package contract

type CreateContractInput struct {
	ClientID    string
	CaseID      string
	Title       string
	Description string
	Type        string
	Amount      string
	StartDate   string
	EndDate     *string
}
