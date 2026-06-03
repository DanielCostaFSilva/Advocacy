package hearing

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	app "legalflow/internal/application/hearing"
)

type CreateHearingUseCase interface {
	Execute(ctx context.Context, input app.CreateHearingInput) (*app.CreateHearingOutput, error)
}

type Handler struct {
	useCase CreateHearingUseCase
}

func NewHandler(useCase CreateHearingUseCase) *Handler {
	return &Handler{useCase: useCase}
}

func (h *Handler) Register(r chi.Router) {
	r.Post("/hearings", h.Create)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateHearingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	scheduledAt, err := time.Parse(time.RFC3339, req.ScheduledAt)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid scheduled_at format"})
		return
	}

	output, err := h.useCase.Execute(r.Context(), app.CreateHearingInput{
		CaseID:      req.CaseID,
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		Location:    req.Location,
		ScheduledAt: scheduledAt,
	})
	if err != nil {
		if errors.Is(err, app.ErrCaseNotFound) {
			writeJSON(w, http.StatusNotFound, ErrorResponse{Error: err.Error()})
			return
		}
		if errors.Is(err, app.ErrInvalidCaseID) || errors.Is(err, app.ErrInvalidInput) {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}

	writeJSON(w, http.StatusCreated, CreateHearingResponse{
		ID:          output.ID,
		CaseID:      output.CaseID,
		Title:       output.Title,
		Description: output.Description,
		Type:        output.Type,
		Location:    output.Location,
		ScheduledAt: output.ScheduledAt,
	})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
