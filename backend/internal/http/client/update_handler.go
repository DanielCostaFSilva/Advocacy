package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	app "legalflow/internal/application/client"
)

type UpdateClientUseCase interface {
	Execute(ctx context.Context, input app.UpdateClientInput) (*app.UpdateClientOutput, error)
}

type UpdateHandler struct {
	useCase UpdateClientUseCase
}

func NewUpdateHandler(useCase UpdateClientUseCase) *UpdateHandler {
	return &UpdateHandler{useCase: useCase}
}

func (h *UpdateHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req CreateClientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	output, err := h.useCase.Execute(r.Context(), app.UpdateClientInput{
		ID:    id,
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	})
	if err != nil {
		if errors.Is(err, app.ErrInvalidClientID) {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}
		if errors.Is(err, app.ErrInvalidInput) {
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

	writeJSON(w, http.StatusOK, CreateClientResponse{
		ID:    output.ID,
		Name:  output.Name,
		CPF:   output.CPF,
		Email: output.Email,
		Phone: output.Phone,
	})
}
