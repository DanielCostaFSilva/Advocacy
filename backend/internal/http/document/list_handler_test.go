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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	app "legalflow/internal/application/document"
	"legalflow/internal/middleware"
)

type mockListCaseDocumentsUseCase struct {
	ExecuteFunc func(ctx context.Context, input app.ListCaseDocumentsInput) (*app.ListCaseDocumentsOutput, error)
}

func (m *mockListCaseDocumentsUseCase) Execute(ctx context.Context, input app.ListCaseDocumentsInput) (*app.ListCaseDocumentsOutput, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, input)
	}
	return nil, nil
}

func setupListHandler(uc *mockListCaseDocumentsUseCase, tokenProv *mockTokenValidator) http.Handler {
	r := chi.NewRouter()
	h := NewListHandler(uc)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(tokenProv))
		h.Register(r)
	})
	return r
}

func TestListDocumentsHandler_ShouldReturn200(t *testing.T) {
	uc := &mockListCaseDocumentsUseCase{
		ExecuteFunc: func(ctx context.Context, input app.ListCaseDocumentsInput) (*app.ListCaseDocumentsOutput, error) {
			return &app.ListCaseDocumentsOutput{
				Data:       []app.DocumentDTO{},
				Page:       1,
				PageSize:   20,
				Total:      0,
				TotalPages: 0,
			}, nil
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupListHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/cases/"+uuid.New().String()+"/documents?page=1&page_size=20", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp ListCaseDocumentsResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, 1, resp.Pagination.Page)
	assert.Equal(t, 20, resp.Pagination.PageSize)
	assert.Equal(t, int64(0), resp.Pagination.Total)
	assert.Equal(t, 0, resp.Pagination.TotalPages)
	assert.Empty(t, resp.Data)
}

func TestListDocumentsHandler_ShouldPassQueryParams(t *testing.T) {
	var capturedInput app.ListCaseDocumentsInput
	uc := &mockListCaseDocumentsUseCase{
		ExecuteFunc: func(ctx context.Context, input app.ListCaseDocumentsInput) (*app.ListCaseDocumentsOutput, error) {
			capturedInput = input
			return &app.ListCaseDocumentsOutput{
				Data:       []app.DocumentDTO{},
				Page:       1,
				PageSize:   20,
				Total:      0,
				TotalPages: 0,
			}, nil
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	caseID := uuid.New().String()
	handler := setupListHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/cases/"+caseID+"/documents?page=2&page_size=10&type=contract&sort=name&order=desc", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, caseID, capturedInput.CaseID)
	assert.Equal(t, 2, capturedInput.Page)
	assert.Equal(t, 10, capturedInput.PageSize)
	assert.Equal(t, "contract", capturedInput.Type)
	assert.Equal(t, "name", capturedInput.Sort)
	assert.Equal(t, "desc", capturedInput.Order)
}

func TestListDocumentsHandler_ShouldReturn400WhenInvalidCaseID(t *testing.T) {
	uc := &mockListCaseDocumentsUseCase{
		ExecuteFunc: func(ctx context.Context, input app.ListCaseDocumentsInput) (*app.ListCaseDocumentsOutput, error) {
			return nil, app.ErrInvalidCaseID
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupListHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/cases/invalid-uuid/documents", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var errResp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	assert.Equal(t, "invalid case id", errResp.Error)
}

func TestListDocumentsHandler_ShouldReturn404WhenCaseNotFound(t *testing.T) {
	uc := &mockListCaseDocumentsUseCase{
		ExecuteFunc: func(ctx context.Context, input app.ListCaseDocumentsInput) (*app.ListCaseDocumentsOutput, error) {
			return nil, app.ErrCaseNotFound
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupListHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/cases/"+uuid.New().String()+"/documents", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	var errResp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	assert.Equal(t, "case not found", errResp.Error)
}

func TestListDocumentsHandler_ShouldReturn401WhenUnauthorized(t *testing.T) {
	uc := &mockListCaseDocumentsUseCase{}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return nil, errors.New("invalid token")
		},
	}

	handler := setupListHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/cases/"+uuid.New().String()+"/documents", nil)
	req.Header.Set("Authorization", "Bearer invalid.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestListDocumentsHandler_ShouldReturn500OnError(t *testing.T) {
	uc := &mockListCaseDocumentsUseCase{
		ExecuteFunc: func(ctx context.Context, input app.ListCaseDocumentsInput) (*app.ListCaseDocumentsOutput, error) {
			return nil, errors.New("unexpected error")
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupListHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/cases/"+uuid.New().String()+"/documents", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestListDocumentsHandler_ShouldReturnData(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	docID := uuid.New().String()
	caseID := uuid.New().String()

	uc := &mockListCaseDocumentsUseCase{
		ExecuteFunc: func(ctx context.Context, input app.ListCaseDocumentsInput) (*app.ListCaseDocumentsOutput, error) {
			return &app.ListCaseDocumentsOutput{
				Data: []app.DocumentDTO{
					{
						ID:         docID,
						CaseID:     caseID,
						Name:       "Contrato de Honorários",
						Type:       "contract",
						FileName:   "contrato.pdf",
						MimeType:   "application/pdf",
						FileSize:   582144,
						StorageKey: "documents/uuid/contrato.pdf",
						CreatedAt:  now,
					},
				},
				Page:       1,
				PageSize:   20,
				Total:      1,
				TotalPages: 1,
			}, nil
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupListHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/cases/"+caseID+"/documents", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp ListCaseDocumentsResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)
	require.Len(t, resp.Data, 1)
	assert.Equal(t, docID, resp.Data[0].ID)
	assert.Equal(t, caseID, resp.Data[0].CaseID)
	assert.Equal(t, "contract", resp.Data[0].Type)
	assert.Equal(t, "contrato.pdf", resp.Data[0].FileName)
	assert.Equal(t, int64(582144), resp.Data[0].FileSize)
	assert.Equal(t, int64(1), resp.Pagination.Total)
	assert.Equal(t, 1, resp.Pagination.TotalPages)
}
