package client

import (
	"bytes"
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

type mockUpdateClientUseCase struct {
	ExecuteFunc func(ctx context.Context, input app.UpdateClientInput) (*app.UpdateClientOutput, error)
}

func (m *mockUpdateClientUseCase) Execute(ctx context.Context, input app.UpdateClientInput) (*app.UpdateClientOutput, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, input)
	}
	return nil, nil
}

func setupUpdateHandler(uc *mockUpdateClientUseCase, tokenProv *mockTokenValidator) http.Handler {
	r := chi.NewRouter()
	h := NewUpdateHandler(uc)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(tokenProv))
		h.Register(r)
	})
	return r
}

func TestUpdateClientHandler_ShouldReturn200(t *testing.T) {
	uc := &mockUpdateClientUseCase{
		ExecuteFunc: func(ctx context.Context, input app.UpdateClientInput) (*app.UpdateClientOutput, error) {
			return &app.UpdateClientOutput{
				ID:    "550e8400-e29b-41d4-a716-446655440000",
				Name:  "Maria da Silva",
				CPF:   "12345678900",
				Email: "maria.silva@example.com",
				Phone: "81988887777",
			}, nil
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupUpdateHandler(uc, tokenProv)
	body := `{"name":"Maria da Silva","email":"maria.silva@example.com","phone":"81988887777"}`
	req := httptest.NewRequest(http.MethodPut, "/clients/550e8400-e29b-41d4-a716-446655440000", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp CreateClientResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", resp.ID)
	assert.Equal(t, "Maria da Silva", resp.Name)
	assert.Equal(t, "12345678900", resp.CPF)
	assert.Equal(t, "maria.silva@example.com", resp.Email)
	assert.Equal(t, "81988887777", resp.Phone)
}

func TestUpdateClientHandler_ShouldReturn400WhenInvalidUUID(t *testing.T) {
	uc := &mockUpdateClientUseCase{
		ExecuteFunc: func(ctx context.Context, input app.UpdateClientInput) (*app.UpdateClientOutput, error) {
			return nil, app.ErrInvalidClientID
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupUpdateHandler(uc, tokenProv)
	body := `{"name":"Maria","email":"maria@example.com","phone":"81999999999"}`
	req := httptest.NewRequest(http.MethodPut, "/clients/invalid-uuid", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateClientHandler_ShouldReturn404WhenClientNotFound(t *testing.T) {
	uc := &mockUpdateClientUseCase{
		ExecuteFunc: func(ctx context.Context, input app.UpdateClientInput) (*app.UpdateClientOutput, error) {
			return nil, app.ErrClientNotFound
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupUpdateHandler(uc, tokenProv)
	body := `{"name":"Maria","email":"maria@example.com","phone":"81999999999"}`
	req := httptest.NewRequest(http.MethodPut, "/clients/550e8400-e29b-41d4-a716-446655440000", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdateClientHandler_ShouldReturn400WhenMalformedJSON(t *testing.T) {
	uc := &mockUpdateClientUseCase{}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupUpdateHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodPut, "/clients/550e8400-e29b-41d4-a716-446655440000", bytes.NewReader([]byte(`{invalid json`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateClientHandler_ShouldReturn401WhenUnauthorized(t *testing.T) {
	uc := &mockUpdateClientUseCase{}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return nil, errors.New("invalid token")
		},
	}

	handler := setupUpdateHandler(uc, tokenProv)
	body := `{"name":"Maria","email":"maria@example.com","phone":"81999999999"}`
	req := httptest.NewRequest(http.MethodPut, "/clients/550e8400-e29b-41d4-a716-446655440000", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer invalid.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestUpdateClientHandler_ShouldReturn500OnInternalError(t *testing.T) {
	uc := &mockUpdateClientUseCase{
		ExecuteFunc: func(ctx context.Context, input app.UpdateClientInput) (*app.UpdateClientOutput, error) {
			return nil, errors.New("unexpected error")
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupUpdateHandler(uc, tokenProv)
	body := `{"name":"Maria","email":"maria@example.com","phone":"81999999999"}`
	req := httptest.NewRequest(http.MethodPut, "/clients/550e8400-e29b-41d4-a716-446655440000", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
