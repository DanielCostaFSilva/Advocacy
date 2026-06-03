package contract

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	app "legalflow/internal/application/contract"
	"legalflow/internal/middleware"
)

type mockCloseContractUseCase struct {
	ExecuteFunc func(ctx context.Context, input app.CloseContractInput) error
}

func (m *mockCloseContractUseCase) Execute(ctx context.Context, input app.CloseContractInput) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, input)
	}
	return nil
}

func setupCloseHandler(uc *mockCloseContractUseCase, tokenProv *mockTokenValidator) http.Handler {
	r := chi.NewRouter()
	h := NewCloseHandler(uc)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(tokenProv))
		h.Register(r)
	})
	return r
}

func TestCloseContractHandler_ShouldReturn204(t *testing.T) {
	contractID := uuid.New().String()

	uc := &mockCloseContractUseCase{
		ExecuteFunc: func(ctx context.Context, input app.CloseContractInput) error {
			return nil
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupCloseHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodDelete, "/contracts/"+contractID, nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Code)
}

func TestCloseContractHandler_ShouldReturn400WhenInvalidID(t *testing.T) {
	uc := &mockCloseContractUseCase{
		ExecuteFunc: func(ctx context.Context, input app.CloseContractInput) error {
			return app.ErrInvalidContractID
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupCloseHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodDelete, "/contracts/invalid-uuid", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var errResp ErrorResponse
	assert.NoError(t, json.NewDecoder(rec.Body).Decode(&errResp))
	assert.Equal(t, "invalid contract id", errResp.Error)
}

func TestCloseContractHandler_ShouldReturn404WhenNotFound(t *testing.T) {
	uc := &mockCloseContractUseCase{
		ExecuteFunc: func(ctx context.Context, input app.CloseContractInput) error {
			return app.ErrContractNotFound
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupCloseHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodDelete, "/contracts/"+uuid.New().String(), nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	var errResp ErrorResponse
	assert.NoError(t, json.NewDecoder(rec.Body).Decode(&errResp))
	assert.Equal(t, "contract not found", errResp.Error)
}

func TestCloseContractHandler_ShouldReturn409WhenAlreadyClosed(t *testing.T) {
	uc := &mockCloseContractUseCase{
		ExecuteFunc: func(ctx context.Context, input app.CloseContractInput) error {
			return app.ErrContractAlreadyClosed
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupCloseHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodDelete, "/contracts/"+uuid.New().String(), nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusConflict, rec.Code)

	var errResp ErrorResponse
	assert.NoError(t, json.NewDecoder(rec.Body).Decode(&errResp))
	assert.Equal(t, "contract already closed", errResp.Error)
}

func TestCloseContractHandler_ShouldReturn401WhenUnauthorized(t *testing.T) {
	uc := &mockCloseContractUseCase{}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return nil, errors.New("invalid token")
		},
	}

	handler := setupCloseHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodDelete, "/contracts/"+uuid.New().String(), nil)
	req.Header.Set("Authorization", "Bearer invalid.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestCloseContractHandler_ShouldReturn500OnInternalError(t *testing.T) {
	uc := &mockCloseContractUseCase{
		ExecuteFunc: func(ctx context.Context, input app.CloseContractInput) error {
			return errors.New("unexpected error")
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupCloseHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodDelete, "/contracts/"+uuid.New().String(), nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
