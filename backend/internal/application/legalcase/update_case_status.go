package legalcase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	domain "legalflow/internal/domain/legalcase"
	timelinedomain "legalflow/internal/domain/timeline"
)

type UpdateCaseStatusUseCase struct {
	repo         domain.CaseRepository
	timelineRepo TimelineEventCreator
}

func NewUpdateCaseStatusUseCase(repo domain.CaseRepository, timelineRepo TimelineEventCreator) *UpdateCaseStatusUseCase {
	return &UpdateCaseStatusUseCase{
		repo:         repo,
		timelineRepo: timelineRepo,
	}
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

	oldStatus := c.Status

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

	event := timelinedomain.NewEvent(c.ID, timelinedomain.EventStatusChanged, fmt.Sprintf("Status alterado de %s para %s", oldStatus, c.Status))
	if err := uc.timelineRepo.Create(ctx, event); err != nil {
		return nil, err
	}

	return &UpdateCaseStatusOutput{
		ID:     c.ID.String(),
		Status: string(c.Status),
	}, nil
}
