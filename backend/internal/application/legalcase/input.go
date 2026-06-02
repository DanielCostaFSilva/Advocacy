package legalcase

type GetCaseByIDInput struct {
	ID string
}

type ListCasesInput struct {
	Page     int
	PageSize int
	Number   string
	ClientID string
	Status   string
	Sort     string
	Order    string
}

type CreateCaseInput struct {
	ClientID    string
	Number      string
	Title       string
	Description string
	Court       string
}
