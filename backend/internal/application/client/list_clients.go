package client

import (
	"context"

	domain "legalflow/internal/domain/client"
)

type ListClientsUseCase struct {
	repo domain.ClientRepository
}

func NewListClientsUseCase(repo domain.ClientRepository) *ListClientsUseCase {
	return &ListClientsUseCase{repo: repo}
}

func (uc *ListClientsUseCase) Execute(ctx context.Context, input ListClientsInput) (*ListClientsOutput, error) {
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

	clients, total, err := uc.repo.List(ctx, domain.ListClientsParams{
		Offset: offset,
		Limit:  pageSize,
		Name:   input.Name,
		CPF:    input.CPF,
		Sort:   input.Sort,
		Order:  input.Order,
	})
	if err != nil {
		return nil, err
	}

	totalPages := int(total / int64(pageSize))
	if total%int64(pageSize) > 0 {
		totalPages++
	}

	data := make([]ClientDTO, len(clients))
	for i, c := range clients {
		data[i] = ClientDTO{
			ID:    c.ID.String(),
			Name:  c.Name,
			CPF:   c.CPF,
			Email: c.Email,
			Phone: c.Phone,
		}
	}

	return &ListClientsOutput{
		Data:       data,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}
