package client

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	app "legalflow/internal/application/client"
)

type DeleteClientUseCase interface {
	Execute(ctx context.Context, input app.DeleteClientInput) error
}

type DeleteHandler struct {
	useCase DeleteClientUseCase
}

func NewDeleteHandler(useCase DeleteClientUseCase) *DeleteHandler {
	return &DeleteHandler{useCase: useCase}
}

func (h *DeleteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.useCase.Execute(r.Context(), app.DeleteClientInput{ID: id})
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

	w.WriteHeader(http.StatusNoContent)
}
