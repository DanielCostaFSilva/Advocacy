package legalcase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	clientdomain "legalflow/internal/domain/client"
	casedomain "legalflow/internal/domain/legalcase"
	timelinedomain "legalflow/internal/domain/timeline"
)

type TimelineEventCreator interface {
	Create(ctx context.Context, event *timelinedomain.TimelineEvent) error
}

type CreateCaseUseCase struct {
	caseRepo      casedomain.CaseRepository
	clientRepo    clientdomain.ClientRepository
	timelineRepo  TimelineEventCreator
}

func NewCreateCaseUseCase(caseRepo casedomain.CaseRepository, clientRepo clientdomain.ClientRepository, timelineRepo TimelineEventCreator) *CreateCaseUseCase {
	return &CreateCaseUseCase{
		caseRepo:      caseRepo,
		clientRepo:    clientRepo,
		timelineRepo:  timelineRepo,
	}
}

func (uc *CreateCaseUseCase) Execute(ctx context.Context, input CreateCaseInput) (*CreateCaseOutput, error) {
	clientID, err := uuid.Parse(input.ClientID)
	if err != nil {
		return nil, ErrInvalidClientID
	}

	client, err := uc.clientRepo.FindByID(ctx, clientID)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, ErrClientNotFound
	}

	existing, err := uc.caseRepo.FindByNumber(ctx, input.Number)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrCaseAlreadyExists
	}

	c, err := casedomain.NewCase(clientID, input.Number, input.Title, input.Description, input.Court)
	if err != nil {
		return nil, err
	}

	if err := uc.caseRepo.Create(ctx, c); err != nil {
		return nil, err
	}

	event := timelinedomain.NewEvent(c.ID, timelinedomain.EventCaseCreated, fmt.Sprintf("Processo %s criado", c.Number))
	if err := uc.timelineRepo.Create(ctx, event); err != nil {
		return nil, err
	}

	return &CreateCaseOutput{
		ID:          c.ID.String(),
		ClientID:    c.ClientID.String(),
		Number:      c.Number,
		Title:       c.Title,
		Description: c.Description,
		Court:       c.Court,
		Status:      string(c.Status),
	}, nil
}
