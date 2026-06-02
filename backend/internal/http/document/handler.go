package document

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	app "legalflow/internal/application/document"
)

type UploadDocumentUseCase interface {
	Execute(ctx context.Context, input app.UploadDocumentInput) (*app.UploadDocumentOutput, error)
}

type Handler struct {
	useCase UploadDocumentUseCase
}

func NewHandler(useCase UploadDocumentUseCase) *Handler {
	return &Handler{useCase: useCase}
}

func (h *Handler) Register(r chi.Router) {
	r.Post("/documents/upload", h.Upload)
}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(MaxDocumentSize); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid multipart form"})
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "file is required"})
		return
	}
	defer file.Close()

	if header.Size == 0 {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "file is empty"})
		return
	}

	if header.Size > MaxDocumentSize {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "file exceeds maximum size of 20MB"})
		return
	}

	mimeType := header.Header.Get("Content-Type")
	if !allowedMimeTypes[mimeType] {
		mimeType = detectMimeType(header.Filename)
		if !allowedMimeTypes[mimeType] {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "file type not allowed"})
			return
		}
	}

	caseID := r.FormValue("case_id")
	if caseID == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "case_id is required"})
		return
	}

	name := r.FormValue("name")
	if name == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "name is required"})
		return
	}

	docType := r.FormValue("type")
	if docType == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "type is required"})
		return
	}

	storageKey := fmt.Sprintf("documents/%s/%s", uuid.New().String(), header.Filename)

	output, err := h.useCase.Execute(r.Context(), app.UploadDocumentInput{
		CaseID:      caseID,
		Name:        name,
		Description: r.FormValue("description"),
		Type:        docType,
		FileName:    header.Filename,
		MimeType:    mimeType,
		FileSize:    header.Size,
		StorageKey:  storageKey,
	})
	if err != nil {
		if errors.Is(err, app.ErrCaseNotFound) {
			writeJSON(w, http.StatusNotFound, ErrorResponse{Error: err.Error()})
			return
		}
		if errors.Is(err, app.ErrInvalidCaseID) || errors.Is(err, app.ErrInvalidInput) {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}

	writeJSON(w, http.StatusCreated, UploadDocumentResponse{
		ID:         output.ID,
		CaseID:     output.CaseID,
		Name:       output.Name,
		Type:       output.Type,
		FileName:   output.FileName,
		MimeType:   output.MimeType,
		FileSize:   output.FileSize,
		StorageKey: output.StorageKey,
	})
}

func detectMimeType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".pdf":
		return "application/pdf"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".doc":
		return "application/msword"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	default:
		return "application/octet-stream"
	}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
