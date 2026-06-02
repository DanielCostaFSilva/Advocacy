package legalcase

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	app "legalflow/internal/application/legalcase"
)

type CreateCaseUseCase interface {
	Execute(ctx context.Context, input app.CreateCaseInput) (*app.CreateCaseOutput, error)
}

type Handler struct {
	useCase CreateCaseUseCase
}

func NewHandler(useCase CreateCaseUseCase) *Handler {
	return &Handler{useCase: useCase}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateCaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	output, err := h.useCase.Execute(r.Context(), app.CreateCaseInput{
		ClientID:    req.ClientID,
		Number:      req.Number,
		Title:       req.Title,
		Description: req.Description,
		Court:       req.Court,
	})
	if err != nil {
		if errors.Is(err, app.ErrCaseAlreadyExists) {
			writeJSON(w, http.StatusConflict, ErrorResponse{Error: err.Error()})
			return
		}
		if errors.Is(err, app.ErrClientNotFound) {
			writeJSON(w, http.StatusNotFound, ErrorResponse{Error: err.Error()})
			return
		}
		if errors.Is(err, app.ErrInvalidClientID) ||
			errors.Is(err, app.ErrInvalidInput) {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}

	writeJSON(w, http.StatusCreated, CreateCaseResponse{
		ID:          output.ID,
		ClientID:    output.ClientID,
		Number:      output.Number,
		Title:       output.Title,
		Description: output.Description,
		Court:       output.Court,
		Status:      output.Status,
	})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
