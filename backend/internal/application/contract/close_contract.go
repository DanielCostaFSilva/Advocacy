package contract

import (
	"context"

	"github.com/google/uuid"
	contractDomain "legalflow/internal/domain/contract"
	timelineDomain "legalflow/internal/domain/timeline"
)

type CloseContractInput struct {
	ContractID string
}

type CloseContractUseCase struct {
	contractRepo contractDomain.Repository
	timelineRepo TimelineRepository
}

func NewCloseContractUseCase(contractRepo contractDomain.Repository, timelineRepo TimelineRepository) *CloseContractUseCase {
	return &CloseContractUseCase{
		contractRepo: contractRepo,
		timelineRepo: timelineRepo,
	}
}

func (uc *CloseContractUseCase) Execute(ctx context.Context, input CloseContractInput) error {
	id, err := uuid.Parse(input.ContractID)
	if err != nil {
		return ErrInvalidContractID
	}

	contract, err := uc.contractRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if contract == nil {
		return ErrContractNotFound
	}

	if !contract.Active {
		return ErrContractAlreadyClosed
	}

	if err := uc.contractRepo.Close(ctx, id); err != nil {
		return ErrContractPersistence
	}

	event := timelineDomain.NewContractClosedEvent(contract)
	if err := uc.timelineRepo.Create(ctx, event); err != nil {
		return ErrTimelinePersistence
	}

	return nil
}
