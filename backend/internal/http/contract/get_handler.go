package contract

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	app "legalflow/internal/application/contract"
)

type GetContractUseCase interface {
	Execute(ctx context.Context, input app.GetContractInput) (*app.GetContractOutput, error)
}

type GetHandler struct {
	useCase GetContractUseCase
}

func NewGetHandler(useCase GetContractUseCase) *GetHandler {
	return &GetHandler{useCase: useCase}
}

func (h *GetHandler) Register(r chi.Router) {
	r.Get("/contracts/{id}", h.Get)
}

func (h *GetHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	output, err := h.useCase.Execute(r.Context(), app.GetContractInput{ID: id})
	if err != nil {
		if errors.Is(err, app.ErrInvalidContractID) {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}
		if errors.Is(err, app.ErrContractNotFound) {
			writeJSON(w, http.StatusNotFound, ErrorResponse{Error: err.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, GetContractResponse{
		ID:          output.ID,
		ClientID:    output.ClientID,
		CaseID:      output.CaseID,
		Title:       output.Title,
		Description: output.Description,
		Type:        output.Type,
		Amount:      output.Amount,
		StartDate:   output.StartDate,
		EndDate:     output.EndDate,
		Active:      output.Active,
		CreatedAt:   output.CreatedAt,
		UpdatedAt:   output.UpdatedAt,
	})
}

type GetContractResponse struct {
	ID          string    `json:"id"`
	ClientID    string    `json:"client_id"`
	CaseID      string    `json:"case_id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Type        string    `json:"type"`
	Amount      string    `json:"amount"`
	StartDate   string    `json:"start_date"`
	EndDate     *string   `json:"end_date,omitempty"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
