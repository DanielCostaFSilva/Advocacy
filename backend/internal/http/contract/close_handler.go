package contract

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	app "legalflow/internal/application/contract"
)

type CloseContractUseCase interface {
	Execute(ctx context.Context, input app.CloseContractInput) error
}

type CloseHandler struct {
	useCase CloseContractUseCase
}

func NewCloseHandler(useCase CloseContractUseCase) *CloseHandler {
	return &CloseHandler{useCase: useCase}
}

func (h *CloseHandler) Register(r chi.Router) {
	r.Delete("/contracts/{id}", h.Close)
}

func (h *CloseHandler) Close(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.useCase.Execute(r.Context(), app.CloseContractInput{ContractID: id})
	if err != nil {
		if errors.Is(err, app.ErrInvalidContractID) {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}
		if errors.Is(err, app.ErrContractNotFound) {
			writeJSON(w, http.StatusNotFound, ErrorResponse{Error: err.Error()})
			return
		}
		if errors.Is(err, app.ErrContractAlreadyClosed) {
			writeJSON(w, http.StatusConflict, ErrorResponse{Error: err.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
