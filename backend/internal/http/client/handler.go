package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	app "legalflow/internal/application/client"
)

type CreateClientUseCase interface {
	Execute(ctx context.Context, input app.CreateClientInput) (*app.CreateClientOutput, error)
}

type Handler struct {
	useCase CreateClientUseCase
}

func NewHandler(useCase CreateClientUseCase) *Handler {
	return &Handler{useCase: useCase}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateClientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	input := app.CreateClientInput{
		Name:  req.Name,
		CPF:   req.CPF,
		Email: req.Email,
		Phone: req.Phone,
	}

	output, err := h.useCase.Execute(r.Context(), input)
	if err != nil {
		if errors.Is(err, app.ErrClientAlreadyExists) {
			writeJSON(w, http.StatusConflict, ErrorResponse{Error: err.Error()})
			return
		}
		if errors.Is(err, app.ErrInvalidInput) {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}

	writeJSON(w, http.StatusCreated, CreateClientResponse{
		ID:    output.ID,
		Name:  output.Name,
		CPF:   output.CPF,
		Email: output.Email,
		Phone: output.Phone,
	})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
