package document

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	app "legalflow/internal/application/document"
	"legalflow/internal/middleware"
)

type mockUploadDocumentUseCase struct {
	ExecuteFunc func(ctx context.Context, input app.UploadDocumentInput) (*app.UploadDocumentOutput, error)
}

func (m *mockUploadDocumentUseCase) Execute(ctx context.Context, input app.UploadDocumentInput) (*app.UploadDocumentOutput, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, input)
	}
	return nil, nil
}

type mockTokenValidator struct {
	ValidateFunc func(token string) (*middleware.TokenClaims, error)
}

func (m *mockTokenValidator) Validate(token string) (*middleware.TokenClaims, error) {
	if m.ValidateFunc != nil {
		return m.ValidateFunc(token)
	}
	return nil, nil
}

func setupHandler(uc *mockUploadDocumentUseCase, tokenProv *mockTokenValidator) http.Handler {
	r := chi.NewRouter()
	h := NewHandler(uc)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(tokenProv))
		h.Register(r)
	})
	return r
}

func multipartBody(caseID, name, docType, description, fileName, mimeType, content string) (bytes.Buffer, string, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	if caseID != "" {
		w.WriteField("case_id", caseID)
	}
	if name != "" {
		w.WriteField("name", name)
	}
	if docType != "" {
		w.WriteField("type", docType)
	}
	if description != "" {
		w.WriteField("description", description)
	}

	if fileName != "" {
		fw, err := w.CreateFormFile("file", fileName)
		if err != nil {
			return buf, "", err
		}
		fw.Write([]byte(content))
	}

	w.Close()
	return buf, w.FormDataContentType(), nil
}

func TestUploadDocumentHandler_ShouldReturn201(t *testing.T) {
	caseID := uuid.New().String()
	docID := uuid.New().String()

	uc := &mockUploadDocumentUseCase{
		ExecuteFunc: func(ctx context.Context, input app.UploadDocumentInput) (*app.UploadDocumentOutput, error) {
			return &app.UploadDocumentOutput{
				ID:         docID,
				CaseID:     caseID,
				Name:       "Contrato de Honorários",
				Type:       "contract",
				FileName:   "contrato.pdf",
				MimeType:   "application/pdf",
				FileSize:   1024,
				StorageKey: "documents/uuid/contrato.pdf",
			}, nil
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	body, contentType, err := multipartBody(caseID, "Contrato de Honorários", "contract", "Desc", "contrato.pdf", "application/pdf", "fake-content")
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/documents/upload", &body)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	req.Header.Set("Content-Type", contentType)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	var resp UploadDocumentResponse
	err = json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, docID, resp.ID)
	assert.Equal(t, caseID, resp.CaseID)
	assert.Equal(t, "contract", resp.Type)
}

func TestUploadDocumentHandler_ShouldReturn400WhenFileMissing(t *testing.T) {
	uc := &mockUploadDocumentUseCase{}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	w.WriteField("case_id", uuid.New().String())
	w.WriteField("name", "Test")
	w.WriteField("type", "contract")
	w.Close()

	req := httptest.NewRequest(http.MethodPost, "/documents/upload", &buf)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	req.Header.Set("Content-Type", w.FormDataContentType())
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUploadDocumentHandler_ShouldReturn400WhenInvalidContentType(t *testing.T) {
	caseID := uuid.New().String()
	uc := &mockUploadDocumentUseCase{}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	body, contentType, err := multipartBody(caseID, "Test", "contract", "", "file.exe", "application/x-msdownload", "bad")
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/documents/upload", &body)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	req.Header.Set("Content-Type", contentType)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUploadDocumentHandler_ShouldReturn400WhenCaseNotFound(t *testing.T) {
	caseID := uuid.New().String()
	uc := &mockUploadDocumentUseCase{
		ExecuteFunc: func(ctx context.Context, input app.UploadDocumentInput) (*app.UploadDocumentOutput, error) {
			return nil, app.ErrCaseNotFound
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	body, contentType, err := multipartBody(caseID, "Test", "contract", "", "file.pdf", "application/pdf", "content")
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/documents/upload", &body)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	req.Header.Set("Content-Type", contentType)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUploadDocumentHandler_ShouldReturn400WhenInvalidCaseID(t *testing.T) {
	uc := &mockUploadDocumentUseCase{
		ExecuteFunc: func(ctx context.Context, input app.UploadDocumentInput) (*app.UploadDocumentOutput, error) {
			return nil, app.ErrInvalidCaseID
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	body, contentType, err := multipartBody("invalid", "Test", "contract", "", "file.pdf", "application/pdf", "content")
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/documents/upload", &body)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	req.Header.Set("Content-Type", contentType)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUploadDocumentHandler_ShouldReturn401WhenMissingJWT(t *testing.T) {
	uc := &mockUploadDocumentUseCase{}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return nil, errors.New("missing token")
		},
	}

	handler := setupHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodPost, "/documents/upload", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestUploadDocumentHandler_ShouldReturn500OnInternalError(t *testing.T) {
	caseID := uuid.New().String()
	uc := &mockUploadDocumentUseCase{
		ExecuteFunc: func(ctx context.Context, input app.UploadDocumentInput) (*app.UploadDocumentOutput, error) {
			return nil, errors.New("unexpected")
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	body, contentType, err := multipartBody(caseID, "Test", "contract", "", "file.pdf", "application/pdf", "content")
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/documents/upload", &body)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	req.Header.Set("Content-Type", contentType)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestUploadDocumentHandler_ShouldInvokeUseCaseWithCorrectInput(t *testing.T) {
	caseID := uuid.New().String()
	var capturedInput app.UploadDocumentInput

	uc := &mockUploadDocumentUseCase{
		ExecuteFunc: func(ctx context.Context, input app.UploadDocumentInput) (*app.UploadDocumentOutput, error) {
			capturedInput = input
			return &app.UploadDocumentOutput{
				ID:         uuid.New().String(),
				CaseID:     caseID,
				Name:       "Test",
				Type:       "contract",
				FileName:   "file.pdf",
				MimeType:   "application/pdf",
				FileSize:   7,
				StorageKey: "documents/uuid/file.pdf",
			}, nil
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	body, contentType, err := multipartBody(caseID, "Contrato Teste", "contract", "Descrição", "file.pdf", "application/pdf", "content")
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/documents/upload", &body)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	req.Header.Set("Content-Type", contentType)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, caseID, capturedInput.CaseID)
	assert.Equal(t, "Contrato Teste", capturedInput.Name)
	assert.Equal(t, "contract", capturedInput.Type)
	assert.Equal(t, "Descrição", capturedInput.Description)
	assert.Equal(t, "file.pdf", capturedInput.FileName)
	assert.Equal(t, "application/pdf", capturedInput.MimeType)
	assert.Equal(t, int64(7), capturedInput.FileSize)
	assert.Contains(t, capturedInput.StorageKey, "documents/")
}

func TestUploadDocumentHandler_ShouldReturn400WhenNameMissing(t *testing.T) {
	caseID := uuid.New().String()
	uc := &mockUploadDocumentUseCase{}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	w.WriteField("case_id", caseID)
	w.WriteField("type", "contract")
	fw, _ := w.CreateFormFile("file", "test.pdf")
	fw.Write([]byte("content"))
	w.Close()

	req := httptest.NewRequest(http.MethodPost, "/documents/upload", &buf)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	req.Header.Set("Content-Type", w.FormDataContentType())
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUploadDocumentHandler_ShouldReturn400WhenTypeMissing(t *testing.T) {
	caseID := uuid.New().String()
	uc := &mockUploadDocumentUseCase{}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	w.WriteField("case_id", caseID)
	w.WriteField("name", "Test")
	fw, _ := w.CreateFormFile("file", "test.pdf")
	fw.Write([]byte("content"))
	w.Close()

	req := httptest.NewRequest(http.MethodPost, "/documents/upload", &buf)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	req.Header.Set("Content-Type", w.FormDataContentType())
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
