package timeline

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	app "legalflow/internal/application/timeline"
)

type GetCaseTimelineUseCase interface {
	Execute(ctx context.Context, input app.GetCaseTimelineInput) (*app.GetCaseTimelineOutput, error)
}

type Handler struct {
	useCase GetCaseTimelineUseCase
}

func NewHandler(useCase GetCaseTimelineUseCase) *Handler {
	return &Handler{useCase: useCase}
}

func (h *Handler) Register(r chi.Router) {
	r.Get("/cases/{id}/timeline", h.GetTimeline)
}

func (h *Handler) GetTimeline(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	output, err := h.useCase.Execute(r.Context(), app.GetCaseTimelineInput{CaseID: id})
	if err != nil {
		if errors.Is(err, app.ErrInvalidCaseID) {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}
		if errors.Is(err, app.ErrCaseNotFound) {
			writeJSON(w, http.StatusNotFound, ErrorResponse{Error: err.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}

	events := make([]EventItem, len(output.Events))
	for i, e := range output.Events {
		events[i] = EventItem{
			ID:          e.ID,
			Type:        e.Type,
			Description: e.Description,
			CreatedAt:   e.CreatedAt,
		}
	}

	writeJSON(w, http.StatusOK, TimelineResponse{Events: events})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
