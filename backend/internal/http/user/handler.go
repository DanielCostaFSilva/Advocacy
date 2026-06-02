package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	app "legalflow/internal/application/user"
)

type CreateUserUseCase interface {
	Execute(ctx context.Context, input app.CreateUserInput) (*app.CreateUserOutput, error)
}

type Handler struct {
	useCase CreateUserUseCase
}

func NewHandler(useCase CreateUserUseCase) *Handler {
	return &Handler{useCase: useCase}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	input := app.CreateUserInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	output, err := h.useCase.Execute(r.Context(), input)
	if err != nil {
		if errors.Is(err, app.ErrEmailAlreadyExists) {
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

	writeJSON(w, http.StatusCreated, CreateUserResponse{
		ID:    output.ID,
		Name:  output.Name,
		Email: output.Email,
	})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
