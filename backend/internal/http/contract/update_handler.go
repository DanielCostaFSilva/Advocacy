package contract

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	app "legalflow/internal/application/contract"
)

type UpdateContractUseCase interface {
	Execute(ctx context.Context, input app.UpdateContractInput) (*app.UpdateContractOutput, error)
}

type UpdateHandler struct {
	useCase UpdateContractUseCase
}

func NewUpdateHandler(useCase UpdateContractUseCase) *UpdateHandler {
	return &UpdateHandler{useCase: useCase}
}

func (h *UpdateHandler) Register(r chi.Router) {
	r.Put("/contracts/{id}", h.Update)
}

func (h *UpdateHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req UpdateContractRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	output, err := h.useCase.Execute(r.Context(), app.UpdateContractInput{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		Amount:      req.Amount,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Active:      req.Active,
	})
	if err != nil {
		if errors.Is(err, app.ErrInvalidContractID) {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}
		if errors.Is(err, app.ErrContractNotFound) {
			writeJSON(w, http.StatusNotFound, ErrorResponse{Error: err.Error()})
			return
		}
		if errors.Is(err, app.ErrInvalidInput) {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, UpdateContractResponse{
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
		UpdatedAt:   output.UpdatedAt,
	})
}

type UpdateContractResponse struct {
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
	UpdatedAt   time.Time `json:"updated_at"`
}
