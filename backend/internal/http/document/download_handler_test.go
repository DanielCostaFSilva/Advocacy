package document

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	app "legalflow/internal/application/document"
	"legalflow/internal/middleware"
)

type mockDownloadDocumentUseCase struct {
	ExecuteFunc func(ctx context.Context, input app.GetDocumentDownloadInput) (*app.GetDocumentDownloadOutput, error)
}

func (m *mockDownloadDocumentUseCase) Execute(ctx context.Context, input app.GetDocumentDownloadInput) (*app.GetDocumentDownloadOutput, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, input)
	}
	return nil, nil
}

func setupDownloadHandler(uc DownloadDocumentUseCase, tokenProv *mockTokenValidator) http.Handler {
	r := chi.NewRouter()
	h := NewDownloadHandler(uc)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(tokenProv))
		h.Register(r)
	})
	return r
}

func TestDownloadDocumentHandler_ShouldReturn200(t *testing.T) {
	now := time.Now()

	uc := &mockDownloadDocumentUseCase{
		ExecuteFunc: func(ctx context.Context, input app.GetDocumentDownloadInput) (*app.GetDocumentDownloadOutput, error) {
			return &app.GetDocumentDownloadOutput{
				DocumentID:  input.DocumentID,
				FileName:    "contrato.pdf",
				DownloadURL: "https://temporary-download.local/documents/" + input.DocumentID,
				ExpiresAt:   now.Add(30 * time.Minute),
			}, nil
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupDownloadHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/documents/550e8400-e29b-41d4-a716-446655440000/download", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp DownloadDocumentResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", resp.DocumentID)
	assert.Equal(t, "contrato.pdf", resp.FileName)
	assert.Contains(t, resp.DownloadURL, "temporary-download.local")
	assert.False(t, resp.ExpiresAt.IsZero())
}

func TestDownloadDocumentHandler_ShouldReturn400WhenInvalidDocumentID(t *testing.T) {
	uc := &mockDownloadDocumentUseCase{
		ExecuteFunc: func(ctx context.Context, input app.GetDocumentDownloadInput) (*app.GetDocumentDownloadOutput, error) {
			return nil, app.ErrInvalidDocumentID
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupDownloadHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/documents/invalid-uuid/download", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var errResp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	assert.Equal(t, "invalid document id", errResp.Error)
}

func TestDownloadDocumentHandler_ShouldReturn404WhenDocumentNotFound(t *testing.T) {
	uc := &mockDownloadDocumentUseCase{
		ExecuteFunc: func(ctx context.Context, input app.GetDocumentDownloadInput) (*app.GetDocumentDownloadOutput, error) {
			return nil, app.ErrDocumentNotFound
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupDownloadHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/documents/550e8400-e29b-41d4-a716-446655440000/download", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	var errResp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	assert.Equal(t, "document not found", errResp.Error)
}

func TestDownloadDocumentHandler_ShouldReturn401WhenMissingJWT(t *testing.T) {
	uc := &mockDownloadDocumentUseCase{}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return nil, errors.New("missing token")
		},
	}

	handler := setupDownloadHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/documents/550e8400-e29b-41d4-a716-446655440000/download", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestDownloadDocumentHandler_ShouldReturn500OnInternalError(t *testing.T) {
	uc := &mockDownloadDocumentUseCase{
		ExecuteFunc: func(ctx context.Context, input app.GetDocumentDownloadInput) (*app.GetDocumentDownloadOutput, error) {
			return nil, errors.New("unexpected error")
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupDownloadHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/documents/550e8400-e29b-41d4-a716-446655440000/download", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
