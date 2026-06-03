package contract

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	app "legalflow/internal/application/contract"
)

type ListContractsUseCase interface {
	Execute(ctx context.Context, input app.ListContractsInput) (*app.ListContractsOutput, error)
}

type ListHandler struct {
	useCase ListContractsUseCase
}

func NewListHandler(useCase ListContractsUseCase) *ListHandler {
	return &ListHandler{useCase: useCase}
}

func (h *ListHandler) Register(r chi.Router) {
	r.Get("/contracts", h.List)
}

func (h *ListHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	input := app.ListContractsInput{
		Page:     page,
		PageSize: pageSize,
		ClientID: r.URL.Query().Get("client_id"),
		CaseID:   r.URL.Query().Get("case_id"),
		Type:     r.URL.Query().Get("type"),
		Active:   r.URL.Query().Get("active"),
		Sort:     r.URL.Query().Get("sort"),
		Order:    r.URL.Query().Get("order"),
	}

	output, err := h.useCase.Execute(r.Context(), input)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}

	resp := ListContractsResponse{
		Data: make([]ContractItem, len(output.Data)),
		Pagination: PaginationInfo{
			Page:       output.Page,
			PageSize:   output.PageSize,
			Total:      output.Total,
			TotalPages: output.TotalPages,
		},
	}
	for i, d := range output.Data {
		resp.Data[i] = ContractItem{
			ID:          d.ID,
			ClientID:    d.ClientID,
			CaseID:      d.CaseID,
			Title:       d.Title,
			Description: d.Description,
			Type:        d.Type,
			Amount:      d.Amount,
			StartDate:   d.StartDate,
			EndDate:     d.EndDate,
			Active:      d.Active,
			CreatedAt:   d.CreatedAt,
		}
	}

	writeJSON(w, http.StatusOK, resp)
}

type ListContractsResponse struct {
	Data       []ContractItem `json:"data"`
	Pagination PaginationInfo `json:"pagination"`
}

type ContractItem struct {
	ID          string  `json:"id"`
	ClientID    string  `json:"client_id"`
	CaseID      string  `json:"case_id"`
	Title       string  `json:"title"`
	Description string  `json:"description,omitempty"`
	Type        string  `json:"type"`
	Amount      string  `json:"amount"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
	Active      bool    `json:"active"`
	CreatedAt   string  `json:"created_at"`
}

type PaginationInfo struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}
