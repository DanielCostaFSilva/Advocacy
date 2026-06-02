package client

import (
	"context"
	"net/http"
	"strconv"

	app "legalflow/internal/application/client"
)

type ListClientsUseCase interface {
	Execute(ctx context.Context, input app.ListClientsInput) (*app.ListClientsOutput, error)
}

type ListHandler struct {
	useCase ListClientsUseCase
}

func NewListHandler(useCase ListClientsUseCase) *ListHandler {
	return &ListHandler{useCase: useCase}
}

func (h *ListHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	page, _ := strconv.Atoi(q.Get("page"))
	pageSize, _ := strconv.Atoi(q.Get("page_size"))

	input := app.ListClientsInput{
		Page:     page,
		PageSize: pageSize,
		Name:     q.Get("name"),
		CPF:      q.Get("cpf"),
		Sort:     q.Get("sort"),
		Order:    q.Get("order"),
	}

	output, err := h.useCase.Execute(r.Context(), input)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}

	resp := ListClientsResponse{
		Data: make([]ClientItem, len(output.Data)),
		Pagination: PaginationInfo{
			Page:       output.Page,
			PageSize:   output.PageSize,
			Total:      output.Total,
			TotalPages: output.TotalPages,
		},
	}
	for i, d := range output.Data {
		resp.Data[i] = ClientItem{
			ID:    d.ID,
			Name:  d.Name,
			CPF:   d.CPF,
			Email: d.Email,
			Phone: d.Phone,
		}
	}

	writeJSON(w, http.StatusOK, resp)
}


