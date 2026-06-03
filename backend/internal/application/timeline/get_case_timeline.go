package timeline

import (
	"context"

	"github.com/google/uuid"
	timelineDomain "legalflow/internal/domain/timeline"
	legalcaseDomain "legalflow/internal/domain/legalcase"
)

type GetCaseTimelineUseCase struct {
	timelineRepo timelineDomain.Repository
	caseRepo     legalcaseDomain.CaseRepository
}

func NewGetCaseTimelineUseCase(
	timelineRepo timelineDomain.Repository,
	caseRepo legalcaseDomain.CaseRepository,
) *GetCaseTimelineUseCase {
	return &GetCaseTimelineUseCase{
		timelineRepo: timelineRepo,
		caseRepo:     caseRepo,
	}
}

func (uc *GetCaseTimelineUseCase) Execute(ctx context.Context, input GetCaseTimelineInput) (*GetCaseTimelineOutput, error) {
	id, err := uuid.Parse(input.CaseID)
	if err != nil {
		return nil, ErrInvalidCaseID
	}

	c, err := uc.caseRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, ErrCaseNotFound
	}

	events, err := uc.timelineRepo.FindByCaseID(ctx, id)
	if err != nil {
		return nil, err
	}

	dtos := make([]TimelineEventDTO, len(events))
	for i, e := range events {
		dto := TimelineEventDTO{
			ID:          e.ID.String(),
			Type:        string(e.Type),
			Description: e.Description,
			CreatedAt:   e.CreatedAt,
		}
		if len(e.Metadata) > 0 {
			dto.Metadata = &e.Metadata
		}
		dtos[i] = dto
	}

	return &GetCaseTimelineOutput{
		Events: dtos,
	}, nil
}
