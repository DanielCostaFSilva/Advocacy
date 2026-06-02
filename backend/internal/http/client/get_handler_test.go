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

type mockGetClientByIDUseCase struct {
	ExecuteFunc func(ctx context.Context, input app.GetClientByIDInput) (*app.GetClientByIDOutput, error)
}

func (m *mockGetClientByIDUseCase) Execute(ctx context.Context, input app.GetClientByIDInput) (*app.GetClientByIDOutput, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, input)
	}
	return nil, nil
}

func setupGetHandler(uc *mockGetClientByIDUseCase, tokenProv *mockTokenValidator) http.Handler {
	r := chi.NewRouter()
	h := NewGetHandler(uc)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(tokenProv))
		h.Register(r)
	})
	return r
}

func TestGetClientByIDHandler_ShouldReturn200(t *testing.T) {
	uc := &mockGetClientByIDUseCase{
		ExecuteFunc: func(ctx context.Context, input app.GetClientByIDInput) (*app.GetClientByIDOutput, error) {
			return &app.GetClientByIDOutput{
				ID:    "550e8400-e29b-41d4-a716-446655440000",
				Name:  "Maria Silva",
				CPF:   "12345678900",
				Email: "maria@example.com",
				Phone: "81999999999",
			}, nil
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupGetHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/clients/550e8400-e29b-41d4-a716-446655440000", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp GetClientResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", resp.ID)
	assert.Equal(t, "Maria Silva", resp.Name)
	assert.Equal(t, "12345678900", resp.CPF)
	assert.Equal(t, "maria@example.com", resp.Email)
	assert.Equal(t, "81999999999", resp.Phone)
}

func TestGetClientByIDHandler_ShouldReturn400WhenInvalidUUID(t *testing.T) {
	uc := &mockGetClientByIDUseCase{
		ExecuteFunc: func(ctx context.Context, input app.GetClientByIDInput) (*app.GetClientByIDOutput, error) {
			return nil, app.ErrInvalidClientID
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupGetHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/clients/invalid-uuid", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var errResp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	assert.Equal(t, "invalid client id", errResp.Error)
}

func TestGetClientByIDHandler_ShouldReturn404WhenClientNotFound(t *testing.T) {
	uc := &mockGetClientByIDUseCase{
		ExecuteFunc: func(ctx context.Context, input app.GetClientByIDInput) (*app.GetClientByIDOutput, error) {
			return nil, app.ErrClientNotFound
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupGetHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/clients/550e8400-e29b-41d4-a716-446655440000", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	var errResp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	assert.Equal(t, "client not found", errResp.Error)
}

func TestGetClientByIDHandler_ShouldReturn401WhenUnauthorized(t *testing.T) {
	uc := &mockGetClientByIDUseCase{}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return nil, errors.New("invalid token")
		},
	}

	handler := setupGetHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/clients/550e8400-e29b-41d4-a716-446655440000", nil)
	req.Header.Set("Authorization", "Bearer invalid.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestGetClientByIDHandler_ShouldReturn500OnInternalError(t *testing.T) {
	uc := &mockGetClientByIDUseCase{
		ExecuteFunc: func(ctx context.Context, input app.GetClientByIDInput) (*app.GetClientByIDOutput, error) {
			return nil, errors.New("unexpected error")
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupGetHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/clients/550e8400-e29b-41d4-a716-446655440000", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var errResp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	assert.Equal(t, "internal server error", errResp.Error)
}
