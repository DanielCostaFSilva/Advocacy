package client

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	app "legalflow/internal/application/client"
)

type GetClientByIDUseCase interface {
	Execute(ctx context.Context, input app.GetClientByIDInput) (*app.GetClientByIDOutput, error)
}

type GetHandler struct {
	useCase GetClientByIDUseCase
}

func NewGetHandler(useCase GetClientByIDUseCase) *GetHandler {
	return &GetHandler{useCase: useCase}
}

func (h *GetHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	output, err := h.useCase.Execute(r.Context(), app.GetClientByIDInput{ID: id})
	if err != nil {
		if errors.Is(err, app.ErrInvalidClientID) {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}
		if errors.Is(err, app.ErrClientNotFound) {
			writeJSON(w, http.StatusNotFound, ErrorResponse{Error: err.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, GetClientResponse{
		ID:        output.ID,
		Name:      output.Name,
		CPF:       output.CPF,
		Email:     output.Email,
		Phone:     output.Phone,
		CreatedAt: output.CreatedAt,
		UpdatedAt: output.UpdatedAt,
	})
}
