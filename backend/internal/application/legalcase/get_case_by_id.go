package legalcase

import (
	"context"

	"github.com/google/uuid"
	domain "legalflow/internal/domain/legalcase"
)

type GetCaseByIDUseCase struct {
	repo domain.CaseRepository
}

func NewGetCaseByIDUseCase(repo domain.CaseRepository) *GetCaseByIDUseCase {
	return &GetCaseByIDUseCase{repo: repo}
}

func (uc *GetCaseByIDUseCase) Execute(ctx context.Context, input GetCaseByIDInput) (*GetCaseByIDOutput, error) {
	id, err := uuid.Parse(input.ID)
	if err != nil {
		return nil, ErrInvalidCaseID
	}

	c, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, ErrCaseNotFound
	}

	return &GetCaseByIDOutput{
		ID:          c.ID.String(),
		ClientID:    c.ClientID.String(),
		Number:      c.Number,
		Title:       c.Title,
		Description: c.Description,
		Court:       c.Court,
		Status:      string(c.Status),
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}, nil
}
