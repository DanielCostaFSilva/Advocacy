package hearing

import (
	"context"
	"time"

	domain "legalflow/internal/domain/hearing"
)

type ListHearingsUseCase struct {
	repo domain.Repository
}

func NewListHearingsUseCase(repo domain.Repository) *ListHearingsUseCase {
	return &ListHearingsUseCase{repo: repo}
}

func (uc *ListHearingsUseCase) Execute(ctx context.Context, input ListHearingsInput) (*ListHearingsOutput, error) {
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

	params := domain.ListHearingsParams{
		Offset: offset,
		Limit:  pageSize,
		CaseID: input.CaseID,
		Type:   input.Type,
		Sort:   input.Sort,
		Order:  input.Order,
	}

	if input.StartDate != "" {
		t, err := time.Parse("2006-01-02", input.StartDate)
		if err == nil {
			params.StartDate = &t
		}
	}
	if input.EndDate != "" {
		t, err := time.Parse("2006-01-02", input.EndDate)
		if err == nil {
			params.EndDate = &t
		}
	}

	hearings, total, err := uc.repo.List(ctx, params)
	if err != nil {
		return nil, err
	}

	totalPages := int(total / int64(pageSize))
	if total%int64(pageSize) > 0 {
		totalPages++
	}

	data := make([]HearingDTO, len(hearings))
	for i, h := range hearings {
		data[i] = HearingDTO{
			ID:          h.ID.String(),
			CaseID:      h.CaseID.String(),
			Title:       h.Title,
			Type:        string(h.Type),
			Location:    h.Location,
			ScheduledAt: h.ScheduledAt,
			CreatedAt:   h.CreatedAt,
		}
	}

	return &ListHearingsOutput{
		Data:       data,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}
