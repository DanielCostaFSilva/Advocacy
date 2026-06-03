package contract

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Contract struct {
	ID          uuid.UUID
	ClientID    uuid.UUID
	CaseID      uuid.UUID
	Title       string
	Description string
	Type        ContractType
	Amount      decimal.Decimal
	StartDate   time.Time
	EndDate     *time.Time
	Active      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewContract(
	clientID uuid.UUID,
	caseID uuid.UUID,
	title string,
	description string,
	contractType ContractType,
	amount decimal.Decimal,
	startDate time.Time,
	endDate *time.Time,
) (*Contract, error) {
	if clientID == uuid.Nil {
		return nil, ErrInvalidClientID
	}
	if caseID == uuid.Nil {
		return nil, ErrInvalidCaseID
	}
	if len(title) < 3 {
		return nil, ErrInvalidTitle
	}
	switch contractType {
	case ContractTypeFixedFee, ContractTypeHourly, ContractTypeSuccessFee, ContractTypeMonthly:
	default:
		return nil, ErrInvalidContractType
	}
	if !amount.IsPositive() {
		return nil, ErrInvalidAmount
	}
	if startDate.IsZero() {
		return nil, ErrInvalidStartDate
	}
	if endDate != nil && !endDate.After(startDate) {
		return nil, ErrInvalidEndDate
	}

	now := time.Now()
	return &Contract{
		ID:          uuid.New(),
		ClientID:    clientID,
		CaseID:      caseID,
		Title:       title,
		Description: description,
		Type:        contractType,
		Amount:      amount,
		StartDate:   startDate,
		EndDate:     endDate,
		Active:      true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}
