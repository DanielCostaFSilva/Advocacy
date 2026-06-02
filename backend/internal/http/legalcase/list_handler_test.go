package legalcase

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	app "legalflow/internal/application/legalcase"
	"legalflow/internal/middleware"
)

type mockListCasesUseCase struct {
	ExecuteFunc func(ctx context.Context, input app.ListCasesInput) (*app.ListCasesOutput, error)
}

func (m *mockListCasesUseCase) Execute(ctx context.Context, input app.ListCasesInput) (*app.ListCasesOutput, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, input)
	}
	return nil, nil
}

func setupListHandler(uc *mockListCasesUseCase, tokenProv *mockTokenValidator) http.Handler {
	r := chi.NewRouter()
	h := NewListHandler(uc)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(tokenProv))
		h.Register(r)
	})
	return r
}

func TestListCasesHandler_ShouldReturn200(t *testing.T) {
	uc := &mockListCasesUseCase{
		ExecuteFunc: func(ctx context.Context, input app.ListCasesInput) (*app.ListCasesOutput, error) {
			return &app.ListCasesOutput{
				Data:       []app.CaseDTO{},
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
	req := httptest.NewRequest(http.MethodGet, "/cases?page=1&page_size=20", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp ListCasesResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, 1, resp.Pagination.Page)
	assert.Equal(t, 20, resp.Pagination.PageSize)
	assert.Equal(t, int64(0), resp.Pagination.Total)
	assert.Equal(t, 0, resp.Pagination.TotalPages)
}

func TestListCasesHandler_ShouldPassQueryParams(t *testing.T) {
	var capturedInput app.ListCasesInput
	uc := &mockListCasesUseCase{
		ExecuteFunc: func(ctx context.Context, input app.ListCasesInput) (*app.ListCasesOutput, error) {
			capturedInput = input
			return &app.ListCasesOutput{
				Data:       []app.CaseDTO{},
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
	req := httptest.NewRequest(http.MethodGet, "/cases?page=2&page_size=10&number=ABC&status=active&sort=title&order=desc", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, 2, capturedInput.Page)
	assert.Equal(t, 10, capturedInput.PageSize)
	assert.Equal(t, "ABC", capturedInput.Number)
	assert.Equal(t, "active", capturedInput.Status)
	assert.Equal(t, "title", capturedInput.Sort)
	assert.Equal(t, "desc", capturedInput.Order)
}

func TestListCasesHandler_ShouldReturn401WhenUnauthorized(t *testing.T) {
	uc := &mockListCasesUseCase{}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return nil, errors.New("invalid token")
		},
	}

	handler := setupListHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/cases", nil)
	req.Header.Set("Authorization", "Bearer invalid.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestListCasesHandler_ShouldReturn500OnError(t *testing.T) {
	uc := &mockListCasesUseCase{
		ExecuteFunc: func(ctx context.Context, input app.ListCasesInput) (*app.ListCasesOutput, error) {
			return nil, errors.New("unexpected error")
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupListHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/cases", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
