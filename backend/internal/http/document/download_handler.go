package document

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	app "legalflow/internal/application/document"
)

type DownloadDocumentUseCase interface {
	Execute(ctx context.Context, input app.GetDocumentDownloadInput) (*app.GetDocumentDownloadOutput, error)
}

type DownloadHandler struct {
	useCase DownloadDocumentUseCase
}

func NewDownloadHandler(useCase DownloadDocumentUseCase) *DownloadHandler {
	return &DownloadHandler{useCase: useCase}
}

func (h *DownloadHandler) Register(r chi.Router) {
	r.Get("/documents/{id}/download", h.Download)
}

func (h *DownloadHandler) Download(w http.ResponseWriter, r *http.Request) {
	docID := chi.URLParam(r, "id")

	output, err := h.useCase.Execute(r.Context(), app.GetDocumentDownloadInput{
		DocumentID: docID,
	})
	if err != nil {
		if errors.Is(err, app.ErrInvalidDocumentID) {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid document id"})
			return
		}
		if errors.Is(err, app.ErrDocumentNotFound) {
			writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "document not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, DownloadDocumentResponse{
		DocumentID:  output.DocumentID,
		FileName:    output.FileName,
		DownloadURL: output.DownloadURL,
		ExpiresAt:   output.ExpiresAt,
	})
}

type DownloadDocumentResponse struct {
	DocumentID  string    `json:"document_id"`
	FileName    string    `json:"file_name"`
	DownloadURL string    `json:"download_url"`
	ExpiresAt   time.Time `json:"expires_at"`
}
