package client

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
	app "legalflow/internal/application/client"
	"legalflow/internal/middleware"
)

type mockListClientsUseCase struct {
	ExecuteFunc func(ctx context.Context, input app.ListClientsInput) (*app.ListClientsOutput, error)
}

func (m *mockListClientsUseCase) Execute(ctx context.Context, input app.ListClientsInput) (*app.ListClientsOutput, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, input)
	}
	return nil, nil
}

func setupListHandler(uc *mockListClientsUseCase, tokenProv *mockTokenValidator) http.Handler {
	r := chi.NewRouter()
	h := NewListHandler(uc)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(tokenProv))
		h.Register(r)
	})
	return r
}

func TestListHandler_ShouldReturn200(t *testing.T) {
	uc := &mockListClientsUseCase{
		ExecuteFunc: func(ctx context.Context, input app.ListClientsInput) (*app.ListClientsOutput, error) {
			return &app.ListClientsOutput{
				Data:       []app.ClientDTO{},
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
	req := httptest.NewRequest(http.MethodGet, "/clients?page=1&page_size=20", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp ListClientsResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, 1, resp.Pagination.Page)
	assert.Equal(t, 20, resp.Pagination.PageSize)
	assert.Equal(t, int64(0), resp.Pagination.Total)
	assert.Equal(t, 0, resp.Pagination.TotalPages)
}

func TestListHandler_ShouldReturn401WhenUnauthorized(t *testing.T) {
	uc := &mockListClientsUseCase{}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return nil, errors.New("invalid token")
		},
	}

	handler := setupListHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/clients", nil)
	req.Header.Set("Authorization", "Bearer invalid.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestListHandler_ShouldReturn500OnError(t *testing.T) {
	uc := &mockListClientsUseCase{
		ExecuteFunc: func(ctx context.Context, input app.ListClientsInput) (*app.ListClientsOutput, error) {
			return nil, errors.New("unexpected error")
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupListHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/clients", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestListHandler_ShouldPassQueryParamsToUseCase(t *testing.T) {
	var capturedInput app.ListClientsInput
	uc := &mockListClientsUseCase{
		ExecuteFunc: func(ctx context.Context, input app.ListClientsInput) (*app.ListClientsOutput, error) {
			capturedInput = input
			return &app.ListClientsOutput{
				Data:       []app.ClientDTO{},
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
	req := httptest.NewRequest(http.MethodGet, "/clients?page=2&page_size=10&name=maria&cpf=12345678900&sort=name&order=desc", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, 2, capturedInput.Page)
	assert.Equal(t, 10, capturedInput.PageSize)
	assert.Equal(t, "maria", capturedInput.Name)
	assert.Equal(t, "12345678900", capturedInput.CPF)
	assert.Equal(t, "name", capturedInput.Sort)
	assert.Equal(t, "desc", capturedInput.Order)
}
