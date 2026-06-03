package contract

type CreateContractOutput struct {
	ID          string
	ClientID    string
	CaseID      string
	Title       string
	Description string
	Type        string
	Amount      string
	StartDate   string
	EndDate     *string
	Active      bool
}
