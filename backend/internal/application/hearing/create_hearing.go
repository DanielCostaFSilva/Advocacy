package hearing

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	caseDomain "legalflow/internal/domain/legalcase"
	hearingDomain "legalflow/internal/domain/hearing"
	timelineDomain "legalflow/internal/domain/timeline"
)

type CreateHearingUseCase struct {
	hearingRepo  hearingDomain.Repository
	caseRepo     caseDomain.CaseRepository
	timelineRepo timelineDomain.Repository
}

func NewCreateHearingUseCase(
	hearingRepo hearingDomain.Repository,
	caseRepo caseDomain.CaseRepository,
	timelineRepo timelineDomain.Repository,
) *CreateHearingUseCase {
	return &CreateHearingUseCase{
		hearingRepo:  hearingRepo,
		caseRepo:     caseRepo,
		timelineRepo: timelineRepo,
	}
}

func (uc *CreateHearingUseCase) Execute(ctx context.Context, input CreateHearingInput) (*CreateHearingOutput, error) {
	caseID, err := uuid.Parse(input.CaseID)
	if err != nil {
		return nil, ErrInvalidCaseID
	}

	c, err := uc.caseRepo.FindByID(ctx, caseID)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, ErrCaseNotFound
	}

	h, err := hearingDomain.NewHearing(caseID, input.Title, input.Description, hearingDomain.HearingType(input.Type), input.Location, input.ScheduledAt)
	if err != nil {
		return nil, ErrInvalidInput
	}

	if err := uc.hearingRepo.Create(ctx, h); err != nil {
		return nil, err
	}

	event := timelineDomain.NewEvent(caseID, timelineDomain.EventHearingCreated, fmt.Sprintf("Audiência agendada: %s", input.Title))
	if err := uc.timelineRepo.Create(ctx, event); err != nil {
		return nil, err
	}

	return &CreateHearingOutput{
		ID:          h.ID.String(),
		CaseID:      h.CaseID.String(),
		Title:       h.Title,
		Description: h.Description,
		Type:        string(h.Type),
		Location:    h.Location,
		ScheduledAt: h.ScheduledAt,
	}, nil
}
