package legalcase

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	app "legalflow/internal/application/legalcase"
)

type UpdateCaseStatusUseCase interface {
	Execute(ctx context.Context, input app.UpdateCaseStatusInput) (*app.UpdateCaseStatusOutput, error)
}

type StatusHandler struct {
	useCase UpdateCaseStatusUseCase
}

func NewStatusHandler(useCase UpdateCaseStatusUseCase) *StatusHandler {
	return &StatusHandler{useCase: useCase}
}

func (h *StatusHandler) Register(r chi.Router) {
	r.Patch("/cases/{id}/status", h.UpdateStatus)
}

func (h *StatusHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	output, err := h.useCase.Execute(r.Context(), app.UpdateCaseStatusInput{
		CaseID: id,
		Status: req.Status,
	})
	if err != nil {
		if errors.Is(err, app.ErrInvalidCaseID) || errors.Is(err, app.ErrInvalidStatus) {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}
		if errors.Is(err, app.ErrInvalidStatusTransition) {
			writeJSON(w, http.StatusUnprocessableEntity, ErrorResponse{Error: err.Error()})
			return
		}
		if errors.Is(err, app.ErrCaseNotFound) {
			writeJSON(w, http.StatusNotFound, ErrorResponse{Error: err.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, UpdateStatusResponse{
		ID:     output.ID,
		Status: output.Status,
	})
}
