package legalcase

import (
	"context"

	"github.com/google/uuid"
	domain "legalflow/internal/domain/legalcase"
)

type UpdateCaseStatusUseCase struct {
	repo domain.CaseRepository
}

func NewUpdateCaseStatusUseCase(repo domain.CaseRepository) *UpdateCaseStatusUseCase {
	return &UpdateCaseStatusUseCase{repo: repo}
}

func (uc *UpdateCaseStatusUseCase) Execute(ctx context.Context, input UpdateCaseStatusInput) (*UpdateCaseStatusOutput, error) {
	id, err := uuid.Parse(input.CaseID)
	if err != nil {
		return nil, ErrInvalidCaseID
	}

	newStatus := domain.CaseStatus(input.Status)
	switch newStatus {
	case domain.CaseStatusDraft, domain.CaseStatusActive, domain.CaseStatusSuspended, domain.CaseStatusClosed:
	default:
		return nil, ErrInvalidStatus
	}

	c, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, ErrCaseNotFound
	}

	if err := c.ChangeStatus(newStatus); err != nil {
		if err == domain.ErrInvalidStatus {
			return nil, ErrInvalidStatus
		}
		if err == domain.ErrInvalidStatusTransition {
			return nil, ErrInvalidStatusTransition
		}
		return nil, err
	}

	if err := uc.repo.UpdateStatus(ctx, c.ID, c.Status); err != nil {
		return nil, err
	}

	return &UpdateCaseStatusOutput{
		ID:     c.ID.String(),
		Status: string(c.Status),
	}, nil
}
