package legalcase

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	app "legalflow/internal/application/legalcase"
)

type ListCasesUseCase interface {
	Execute(ctx context.Context, input app.ListCasesInput) (*app.ListCasesOutput, error)
}

type ListHandler struct {
	useCase ListCasesUseCase
}

func NewListHandler(useCase ListCasesUseCase) *ListHandler {
	return &ListHandler{useCase: useCase}
}

func (h *ListHandler) Register(r chi.Router) {
	r.Get("/cases", h.List)
}

func (h *ListHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	input := app.ListCasesInput{
		Page:     page,
		PageSize: pageSize,
		Number:   r.URL.Query().Get("number"),
		ClientID: r.URL.Query().Get("client_id"),
		Status:   r.URL.Query().Get("status"),
		Sort:     r.URL.Query().Get("sort"),
		Order:    r.URL.Query().Get("order"),
	}

	output, err := h.useCase.Execute(r.Context(), input)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}

	resp := ListCasesResponse{
		Data: make([]CaseItem, len(output.Data)),
		Pagination: PaginationInfo{
			Page:       output.Page,
			PageSize:   output.PageSize,
			Total:      output.Total,
			TotalPages: output.TotalPages,
		},
	}
	for i, c := range output.Data {
		resp.Data[i] = CaseItem{
			ID:        c.ID,
			ClientID:  c.ClientID,
			Number:    c.Number,
			Title:     c.Title,
			Court:     c.Court,
			Status:    c.Status,
			CreatedAt: c.CreatedAt,
		}
	}

	writeJSON(w, http.StatusOK, resp)
}
