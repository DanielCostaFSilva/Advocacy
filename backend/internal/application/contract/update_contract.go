package contract

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	contractDomain "legalflow/internal/domain/contract"
	timelineDomain "legalflow/internal/domain/timeline"
)

type UpdateContractInput struct {
	ID          string
	Title       string
	Description string
	Type        string
	Amount      string
	StartDate   string
	EndDate     *string
	Active      bool
}

type UpdateContractOutput struct {
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
	UpdatedAt   time.Time
}

type UpdateContractUseCase struct {
	contractRepo  contractDomain.Repository
	timelineRepo  TimelineRepository
}

func NewUpdateContractUseCase(contractRepo contractDomain.Repository, timelineRepo TimelineRepository) *UpdateContractUseCase {
	return &UpdateContractUseCase{
		contractRepo: contractRepo,
		timelineRepo: timelineRepo,
	}
}

func (uc *UpdateContractUseCase) Execute(ctx context.Context, input UpdateContractInput) (*UpdateContractOutput, error) {
	id, err := uuid.Parse(input.ID)
	if err != nil {
		return nil, ErrInvalidContractID
	}

	contract, err := uc.contractRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if contract == nil {
		return nil, ErrContractNotFound
	}

	oldAmount := contract.Amount.String()

	amount, err := decimal.NewFromString(input.Amount)
	if err != nil {
		return nil, ErrInvalidInput
	}

	startDate, err := time.Parse(time.RFC3339, input.StartDate)
	if err != nil {
		return nil, ErrInvalidInput
	}

	var endDate *time.Time
	if input.EndDate != nil {
		parsed, err := time.Parse(time.RFC3339, *input.EndDate)
		if err != nil {
			return nil, ErrInvalidInput
		}
		endDate = &parsed
	}

	updated, err := contractDomain.UpdateContract(
		contract,
		input.Title,
		input.Description,
		contractDomain.ContractType(input.Type),
		amount,
		startDate,
		endDate,
		input.Active,
	)
	if err != nil {
		return nil, ErrInvalidInput
	}

	if err := uc.contractRepo.Update(ctx, updated); err != nil {
		return nil, ErrContractPersistence
	}

	event := timelineDomain.NewContractUpdatedEvent(updated, oldAmount)
	if err := uc.timelineRepo.Create(ctx, event); err != nil {
		return nil, ErrTimelinePersistence
	}

	var outputEndDate *string
	if updated.EndDate != nil {
		s := updated.EndDate.Format(time.RFC3339)
		outputEndDate = &s
	}

	return &UpdateContractOutput{
		ID:          updated.ID.String(),
		ClientID:    updated.ClientID.String(),
		CaseID:      updated.CaseID.String(),
		Title:       updated.Title,
		Description: updated.Description,
		Type:        string(updated.Type),
		Amount:      updated.Amount.String(),
		StartDate:   updated.StartDate.Format(time.RFC3339),
		EndDate:     outputEndDate,
		Active:      updated.Active,
		UpdatedAt:   updated.UpdatedAt,
	}, nil
}
