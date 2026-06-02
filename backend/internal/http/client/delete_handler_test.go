package client

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	app "legalflow/internal/application/client"
	"legalflow/internal/middleware"
)

type mockDeleteClientUseCase struct {
	ExecuteFunc func(ctx context.Context, input app.DeleteClientInput) error
}

func (m *mockDeleteClientUseCase) Execute(ctx context.Context, input app.DeleteClientInput) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, input)
	}
	return nil
}

func setupDeleteHandler(uc *mockDeleteClientUseCase, tokenProv *mockTokenValidator) http.Handler {
	r := chi.NewRouter()
	h := NewDeleteHandler(uc)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(tokenProv))
		h.Register(r)
	})
	return r
}

func TestDeleteClientHandler_ShouldReturn204(t *testing.T) {
	uc := &mockDeleteClientUseCase{
		ExecuteFunc: func(ctx context.Context, input app.DeleteClientInput) error {
			return nil
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupDeleteHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodDelete, "/clients/550e8400-e29b-41d4-a716-446655440000", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestDeleteClientHandler_ShouldReturn400WhenInvalidUUID(t *testing.T) {
	uc := &mockDeleteClientUseCase{
		ExecuteFunc: func(ctx context.Context, input app.DeleteClientInput) error {
			return app.ErrInvalidClientID
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupDeleteHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodDelete, "/clients/invalid-uuid", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDeleteClientHandler_ShouldReturn404WhenClientNotFound(t *testing.T) {
	uc := &mockDeleteClientUseCase{
		ExecuteFunc: func(ctx context.Context, input app.DeleteClientInput) error {
			return app.ErrClientNotFound
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupDeleteHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodDelete, "/clients/550e8400-e29b-41d4-a716-446655440000", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteClientHandler_ShouldReturn401WhenUnauthorized(t *testing.T) {
	uc := &mockDeleteClientUseCase{}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return nil, errors.New("invalid token")
		},
	}

	handler := setupDeleteHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodDelete, "/clients/550e8400-e29b-41d4-a716-446655440000", nil)
	req.Header.Set("Authorization", "Bearer invalid.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestDeleteClientHandler_ShouldReturn500OnInternalError(t *testing.T) {
	uc := &mockDeleteClientUseCase{
		ExecuteFunc: func(ctx context.Context, input app.DeleteClientInput) error {
			return errors.New("unexpected error")
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupDeleteHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodDelete, "/clients/550e8400-e29b-41d4-a716-446655440000", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
