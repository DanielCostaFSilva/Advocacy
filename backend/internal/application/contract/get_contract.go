package contract

import (
	"context"
	"time"

	"github.com/google/uuid"
	domain "legalflow/internal/domain/contract"
)

type GetContractInput struct {
	ID string
}

type GetContractOutput struct {
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
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type GetContractUseCase struct {
	repo domain.Repository
}

func NewGetContractUseCase(repo domain.Repository) *GetContractUseCase {
	return &GetContractUseCase{repo: repo}
}

func (uc *GetContractUseCase) Execute(ctx context.Context, input GetContractInput) (*GetContractOutput, error) {
	id, err := uuid.Parse(input.ID)
	if err != nil {
		return nil, ErrInvalidContractID
	}

	contract, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if contract == nil {
		return nil, ErrContractNotFound
	}

	var endDate *string
	if contract.EndDate != nil {
		s := contract.EndDate.Format(time.RFC3339)
		endDate = &s
	}

	return &GetContractOutput{
		ID:          contract.ID.String(),
		ClientID:    contract.ClientID.String(),
		CaseID:      contract.CaseID.String(),
		Title:       contract.Title,
		Description: contract.Description,
		Type:        string(contract.Type),
		Amount:      contract.Amount.String(),
		StartDate:   contract.StartDate.Format(time.RFC3339),
		EndDate:     endDate,
		Active:      contract.Active,
		CreatedAt:   contract.CreatedAt,
		UpdatedAt:   contract.UpdatedAt,
	}, nil
}
