package contract

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	app "legalflow/internal/application/contract"
)

type CreateContractUseCase interface {
	Execute(ctx context.Context, input app.CreateContractInput) (*app.CreateContractOutput, error)
}

type CreateHandler struct {
	useCase CreateContractUseCase
}

func NewCreateHandler(useCase CreateContractUseCase) *CreateHandler {
	return &CreateHandler{useCase: useCase}
}

func (h *CreateHandler) Register(r chi.Router) {
	r.Post("/contracts", h.Create)
}

func (h *CreateHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateContractRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	if req.ClientID == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "client_id is required"})
		return
	}
	if req.CaseID == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "case_id is required"})
		return
	}
	if req.Title == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "title is required"})
		return
	}
	if req.Type == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "type is required"})
		return
	}
	if req.Amount == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "amount is required"})
		return
	}
	if req.StartDate == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "start_date is required"})
		return
	}

	output, err := h.useCase.Execute(r.Context(), app.CreateContractInput{
		ClientID:    req.ClientID,
		CaseID:      req.CaseID,
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		Amount:      req.Amount,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	})
	if err != nil {
		if errors.Is(err, app.ErrClientNotFound) {
			writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "client not found"})
			return
		}
		if errors.Is(err, app.ErrCaseNotFound) {
			writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "case not found"})
			return
		}
		if errors.Is(err, app.ErrInvalidClientID) {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid client id"})
			return
		}
		if errors.Is(err, app.ErrInvalidCaseID) {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid case id"})
			return
		}
		if errors.Is(err, app.ErrInvalidInput) {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid input"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}

	writeJSON(w, http.StatusCreated, CreateContractResponse{
		ID:          output.ID,
		ClientID:    output.ClientID,
		CaseID:      output.CaseID,
		Title:       output.Title,
		Description: output.Description,
		Type:        output.Type,
		Amount:      output.Amount,
		StartDate:   output.StartDate,
		EndDate:     output.EndDate,
		Active:      output.Active,
	})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
