package legalcase

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	app "legalflow/internal/application/legalcase"
)

type GetCaseByIDUseCase interface {
	Execute(ctx context.Context, input app.GetCaseByIDInput) (*app.GetCaseByIDOutput, error)
}

type GetHandler struct {
	useCase GetCaseByIDUseCase
}

func NewGetHandler(useCase GetCaseByIDUseCase) *GetHandler {
	return &GetHandler{useCase: useCase}
}

func (h *GetHandler) Register(r chi.Router) {
	r.Get("/cases/{id}", h.Get)
}

func (h *GetHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	output, err := h.useCase.Execute(r.Context(), app.GetCaseByIDInput{ID: id})
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

	writeJSON(w, http.StatusOK, GetCaseResponse{
		ID:          output.ID,
		ClientID:    output.ClientID,
		Number:      output.Number,
		Title:       output.Title,
		Description: output.Description,
		Court:       output.Court,
		Status:      output.Status,
		CreatedAt:   output.CreatedAt,
		UpdatedAt:   output.UpdatedAt,
	})
}
