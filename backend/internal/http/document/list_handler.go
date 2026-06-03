package document

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	app "legalflow/internal/application/document"
)

type ListCaseDocumentsUseCase interface {
	Execute(ctx context.Context, input app.ListCaseDocumentsInput) (*app.ListCaseDocumentsOutput, error)
}

type ListHandler struct {
	useCase ListCaseDocumentsUseCase
}

func NewListHandler(useCase ListCaseDocumentsUseCase) *ListHandler {
	return &ListHandler{useCase: useCase}
}

func (h *ListHandler) Register(r chi.Router) {
	r.Get("/cases/{id}/documents", h.List)
}

func (h *ListHandler) List(w http.ResponseWriter, r *http.Request) {
	caseID := chi.URLParam(r, "id")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	output, err := h.useCase.Execute(r.Context(), app.ListCaseDocumentsInput{
		CaseID:   caseID,
		Page:     page,
		PageSize: pageSize,
		Type:     r.URL.Query().Get("type"),
		Sort:     r.URL.Query().Get("sort"),
		Order:    r.URL.Query().Get("order"),
	})
	if err != nil {
		switch err {
		case app.ErrInvalidCaseID:
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid case id"})
		case app.ErrCaseNotFound:
			writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "case not found"})
		default:
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		}
		return
	}

	resp := ListCaseDocumentsResponse{
		Data: make([]DocumentItem, len(output.Data)),
		Pagination: PaginationInfo{
			Page:       output.Page,
			PageSize:   output.PageSize,
			Total:      output.Total,
			TotalPages: output.TotalPages,
		},
	}
	for i, d := range output.Data {
		resp.Data[i] = DocumentItem{
			ID:         d.ID,
			CaseID:     d.CaseID,
			Name:       d.Name,
			Type:       d.Type,
			FileName:   d.FileName,
			MimeType:   d.MimeType,
			FileSize:   d.FileSize,
			StorageKey: d.StorageKey,
			CreatedAt:  d.CreatedAt,
		}
	}

	writeJSON(w, http.StatusOK, resp)
}

type ListCaseDocumentsResponse struct {
	Data       []DocumentItem `json:"data"`
	Pagination PaginationInfo `json:"pagination"`
}

type DocumentItem struct {
	ID         string    `json:"id"`
	CaseID     string    `json:"case_id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	FileName   string    `json:"file_name"`
	MimeType   string    `json:"mime_type"`
	FileSize   int64     `json:"file_size"`
	StorageKey string    `json:"storage_key"`
	CreatedAt  time.Time `json:"created_at"`
}

type PaginationInfo struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}
