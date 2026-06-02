package legalcase

import (
	"context"

	domain "legalflow/internal/domain/legalcase"
)

type ListCasesUseCase struct {
	repo domain.CaseRepository
}

func NewListCasesUseCase(repo domain.CaseRepository) *ListCasesUseCase {
	return &ListCasesUseCase{repo: repo}
}

func (uc *ListCasesUseCase) Execute(ctx context.Context, input ListCasesInput) (*ListCasesOutput, error) {
	page := input.Page
	if page < 1 {
		page = 1
	}

	pageSize := input.PageSize
	if pageSize < 1 {
		pageSize = 20
	} else if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	cases, total, err := uc.repo.List(ctx, domain.ListCasesParams{
		Offset:   offset,
		Limit:    pageSize,
		Number:   input.Number,
		ClientID: input.ClientID,
		Status:   input.Status,
		Sort:     input.Sort,
		Order:    input.Order,
	})
	if err != nil {
		return nil, err
	}

	totalPages := int(total / int64(pageSize))
	if total%int64(pageSize) > 0 {
		totalPages++
	}

	data := make([]CaseDTO, len(cases))
	for i, c := range cases {
		data[i] = CaseDTO{
			ID:        c.ID.String(),
			ClientID:  c.ClientID.String(),
			Number:    c.Number,
			Title:     c.Title,
			Court:     c.Court,
			Status:    string(c.Status),
			CreatedAt: c.CreatedAt,
		}
	}

	return &ListCasesOutput{
		Data:       data,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}
