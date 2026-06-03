package contract

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	clientDomain "legalflow/internal/domain/client"
	contractDomain "legalflow/internal/domain/contract"
	caseDomain "legalflow/internal/domain/legalcase"
	timelineDomain "legalflow/internal/domain/timeline"
)

type ClientRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*clientDomain.Client, error)
}

type CaseRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*caseDomain.Case, error)
}

type ContractRepository interface {
	Create(ctx context.Context, contract *contractDomain.Contract) error
}

type TimelineRepository interface {
	Create(ctx context.Context, event *timelineDomain.TimelineEvent) error
}

type CreateContractUseCase struct {
	clientRepo   ClientRepository
	caseRepo     CaseRepository
	contractRepo ContractRepository
	timelineRepo TimelineRepository
}

func NewCreateContractUseCase(
	clientRepo ClientRepository,
	caseRepo CaseRepository,
	contractRepo ContractRepository,
	timelineRepo TimelineRepository,
) *CreateContractUseCase {
	return &CreateContractUseCase{
		clientRepo:   clientRepo,
		caseRepo:     caseRepo,
		contractRepo: contractRepo,
		timelineRepo: timelineRepo,
	}
}

func (uc *CreateContractUseCase) Execute(ctx context.Context, input CreateContractInput) (*CreateContractOutput, error) {
	clientID, err := uuid.Parse(input.ClientID)
	if err != nil {
		return nil, ErrInvalidClientID
	}

	caseID, err := uuid.Parse(input.CaseID)
	if err != nil {
		return nil, ErrInvalidCaseID
	}

	client, err := uc.clientRepo.FindByID(ctx, clientID)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, ErrClientNotFound
	}

	c, err := uc.caseRepo.FindByID(ctx, caseID)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, ErrCaseNotFound
	}

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

	contract, err := contractDomain.NewContract(
		clientID, caseID, input.Title, input.Description,
		contractDomain.ContractType(input.Type), amount, startDate, endDate,
	)
	if err != nil {
		return nil, ErrInvalidInput
	}

	if err := uc.contractRepo.Create(ctx, contract); err != nil {
		return nil, ErrContractPersistence
	}

	event := timelineDomain.NewContractCreatedEvent(contract.CaseID, contract)
	if err := uc.timelineRepo.Create(ctx, event); err != nil {
		return nil, ErrTimelinePersistence
	}

	var outputEndDate *string
	if contract.EndDate != nil {
		s := contract.EndDate.Format(time.RFC3339)
		outputEndDate = &s
	}

	return &CreateContractOutput{
		ID:          contract.ID.String(),
		ClientID:    contract.ClientID.String(),
		CaseID:      contract.CaseID.String(),
		Title:       contract.Title,
		Description: contract.Description,
		Type:        string(contract.Type),
		Amount:      contract.Amount.String(),
		StartDate:   contract.StartDate.Format(time.RFC3339),
		EndDate:     outputEndDate,
		Active:      contract.Active,
	}, nil
}
