package contract

import (
	"context"
	"time"

	domain "legalflow/internal/domain/contract"
)

type ListContractsInput struct {
	Page     int
	PageSize int
	ClientID string
	CaseID   string
	Type     string
	Active   string
	Sort     string
	Order    string
}

type ContractDTO struct {
	ID          string `json:"id"`
	ClientID    string `json:"client_id"`
	CaseID      string `json:"case_id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type"`
	Amount      string `json:"amount"`
	StartDate   string `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
	Active      bool   `json:"active"`
	CreatedAt   string `json:"created_at"`
}

type ListContractsOutput struct {
	Data       []ContractDTO
	Page       int
	PageSize   int
	Total      int64
	TotalPages int
}

type ListContractsUseCase struct {
	repo domain.Repository
}

func NewListContractsUseCase(repo domain.Repository) *ListContractsUseCase {
	return &ListContractsUseCase{repo: repo}
}

func (uc *ListContractsUseCase) Execute(ctx context.Context, input ListContractsInput) (*ListContractsOutput, error) {
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

	var active *bool
	if input.Active != "" {
		v := input.Active == "true"
		active = &v
	}

	params := domain.ListContractsParams{
		Offset:   offset,
		Limit:    pageSize,
		ClientID: input.ClientID,
		CaseID:   input.CaseID,
		Type:     input.Type,
		Active:   active,
		Sort:     input.Sort,
		Order:    input.Order,
	}

	contracts, total, err := uc.repo.List(ctx, params)
	if err != nil {
		return nil, err
	}

	totalPages := int(total / int64(pageSize))
	if total%int64(pageSize) > 0 {
		totalPages++
	}

	data := make([]ContractDTO, len(contracts))
	for i, c := range contracts {
		var endDate *string
		if c.EndDate != nil {
			s := c.EndDate.Format(time.RFC3339)
			endDate = &s
		}

		data[i] = ContractDTO{
			ID:          c.ID.String(),
			ClientID:    c.ClientID.String(),
			CaseID:      c.CaseID.String(),
			Title:       c.Title,
			Description: c.Description,
			Type:        string(c.Type),
			Amount:      c.Amount.String(),
			StartDate:   c.StartDate.Format(time.RFC3339),
			EndDate:     endDate,
			Active:      c.Active,
			CreatedAt:   c.CreatedAt.Format(time.RFC3339),
		}
	}

	return &ListContractsOutput{
		Data:       data,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}
