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

type mockCreateClientUseCase struct {
	ExecuteFunc func(ctx context.Context, input app.CreateClientInput) (*app.CreateClientOutput, error)
}

func (m *mockCreateClientUseCase) Execute(ctx context.Context, input app.CreateClientInput) (*app.CreateClientOutput, error) {
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

func setupHandler(uc *mockCreateClientUseCase, tokenProv *mockTokenValidator) http.Handler {
	r := chi.NewRouter()
	h := NewHandler(uc)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(tokenProv))
		h.Register(r)
	})
	return r
}

func TestHandler_CreateClient_ShouldReturn201(t *testing.T) {
	uc := &mockCreateClientUseCase{
		ExecuteFunc: func(ctx context.Context, input app.CreateClientInput) (*app.CreateClientOutput, error) {
			return &app.CreateClientOutput{
				ID:    "uuid-123",
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

	handler := setupHandler(uc, tokenProv)
	body := `{"name":"Maria Silva","cpf":"12345678900","email":"maria@example.com","phone":"81999999999"}`
	req := httptest.NewRequest(http.MethodPost, "/clients", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	var resp CreateClientResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, "uuid-123", resp.ID)
	assert.Equal(t, "Maria Silva", resp.Name)
	assert.Equal(t, "12345678900", resp.CPF)
	assert.Equal(t, "maria@example.com", resp.Email)
	assert.Equal(t, "81999999999", resp.Phone)
}

func TestHandler_CreateClient_ShouldReturn409WhenDuplicateCPF(t *testing.T) {
	uc := &mockCreateClientUseCase{
		ExecuteFunc: func(ctx context.Context, input app.CreateClientInput) (*app.CreateClientOutput, error) {
			return nil, app.ErrClientAlreadyExists
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	body := `{"name":"Maria Silva","cpf":"12345678900","email":"maria@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/clients", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusConflict, rec.Code)

	var errResp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	assert.Equal(t, "client already exists", errResp.Error)
}

func TestHandler_CreateClient_ShouldReturn400WhenInvalidInput(t *testing.T) {
	uc := &mockCreateClientUseCase{
		ExecuteFunc: func(ctx context.Context, input app.CreateClientInput) (*app.CreateClientOutput, error) {
			return nil, app.ErrInvalidInput
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	body := `{"name":"","cpf":"","email":"","phone":""}`
	req := httptest.NewRequest(http.MethodPost, "/clients", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var errResp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	assert.Equal(t, "invalid input", errResp.Error)
}

func TestHandler_CreateClient_ShouldReturn400WhenMalformedJSON(t *testing.T) {
	uc := &mockCreateClientUseCase{}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodPost, "/clients", bytes.NewReader([]byte(`{invalid json`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandler_CreateClient_ShouldReturn400WhenBodyMissing(t *testing.T) {
	uc := &mockCreateClientUseCase{}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodPost, "/clients", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandler_CreateClient_ShouldReturn401WhenMissingJWT(t *testing.T) {
	uc := &mockCreateClientUseCase{}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return nil, errors.New("invalid token")
		},
	}

	handler := setupHandler(uc, tokenProv)
	body := `{"name":"Maria Silva","cpf":"12345678900"}`
	req := httptest.NewRequest(http.MethodPost, "/clients", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer invalid.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestHandler_CreateClient_ShouldReturn500OnUnexpectedError(t *testing.T) {
	uc := &mockCreateClientUseCase{
		ExecuteFunc: func(ctx context.Context, input app.CreateClientInput) (*app.CreateClientOutput, error) {
			return nil, errors.New("unexpected db error")
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	body := `{"name":"Maria Silva","cpf":"12345678900","email":"maria@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/clients", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var errResp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	assert.Equal(t, "internal server error", errResp.Error)
}

func TestHandler_CreateClient_ShouldInvokeUseCaseWithCorrectInput(t *testing.T) {
	var capturedInput app.CreateClientInput
	uc := &mockCreateClientUseCase{
		ExecuteFunc: func(ctx context.Context, input app.CreateClientInput) (*app.CreateClientOutput, error) {
			capturedInput = input
			return &app.CreateClientOutput{
				ID:    "uuid-456",
				Name:  input.Name,
				CPF:   input.CPF,
				Email: input.Email,
				Phone: input.Phone,
			}, nil
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	body := `{"name":"Jane Doe","cpf":"99988877766","email":"jane@example.com","phone":"11988888888"}`
	req := httptest.NewRequest(http.MethodPost, "/clients", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "Jane Doe", capturedInput.Name)
	assert.Equal(t, "99988877766", capturedInput.CPF)
	assert.Equal(t, "jane@example.com", capturedInput.Email)
	assert.Equal(t, "11988888888", capturedInput.Phone)
}
