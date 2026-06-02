package legalcase

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
	app "legalflow/internal/application/legalcase"
	"legalflow/internal/middleware"
)

type mockCreateCaseUseCase struct {
	ExecuteFunc func(ctx context.Context, input app.CreateCaseInput) (*app.CreateCaseOutput, error)
}

func (m *mockCreateCaseUseCase) Execute(ctx context.Context, input app.CreateCaseInput) (*app.CreateCaseOutput, error) {
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

func setupHandler(uc *mockCreateCaseUseCase, tokenProv *mockTokenValidator) http.Handler {
	r := chi.NewRouter()
	h := NewHandler(uc)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(tokenProv))
		h.Register(r)
	})
	return r
}

func TestCreateCaseHandler_ShouldReturn201(t *testing.T) {
	uc := &mockCreateCaseUseCase{
		ExecuteFunc: func(ctx context.Context, input app.CreateCaseInput) (*app.CreateCaseOutput, error) {
			return &app.CreateCaseOutput{
				ID:          "uuid-123",
				ClientID:    "550e8400-e29b-41d4-a716-446655440000",
				Number:      "0001234-56.2026.8.17.0001",
				Title:       "Ação de Cobrança",
				Description: "Descrição",
				Court:       "TJPE",
				Status:      "draft",
			}, nil
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	body := `{"client_id":"550e8400-e29b-41d4-a716-446655440000","number":"0001234-56.2026.8.17.0001","title":"Ação de Cobrança","description":"Descrição","court":"TJPE"}`
	req := httptest.NewRequest(http.MethodPost, "/cases", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	var resp CreateCaseResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, "uuid-123", resp.ID)
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", resp.ClientID)
	assert.Equal(t, "0001234-56.2026.8.17.0001", resp.Number)
	assert.Equal(t, "Ação de Cobrança", resp.Title)
	assert.Equal(t, "Descrição", resp.Description)
	assert.Equal(t, "TJPE", resp.Court)
	assert.Equal(t, "draft", resp.Status)
}

func TestCreateCaseHandler_ShouldReturn409WhenDuplicateNumber(t *testing.T) {
	uc := &mockCreateCaseUseCase{
		ExecuteFunc: func(ctx context.Context, input app.CreateCaseInput) (*app.CreateCaseOutput, error) {
			return nil, app.ErrCaseAlreadyExists
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	body := `{"client_id":"550e8400-e29b-41d4-a716-446655440000","number":"0001234-56.2026.8.17.0001","title":"Ação","court":"TJPE"}`
	req := httptest.NewRequest(http.MethodPost, "/cases", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusConflict, rec.Code)

	var errResp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	assert.Equal(t, "case already exists", errResp.Error)
}

func TestCreateCaseHandler_ShouldReturn404WhenClientNotFound(t *testing.T) {
	uc := &mockCreateCaseUseCase{
		ExecuteFunc: func(ctx context.Context, input app.CreateCaseInput) (*app.CreateCaseOutput, error) {
			return nil, app.ErrClientNotFound
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	body := `{"client_id":"550e8400-e29b-41d4-a716-446655440000","number":"0001234-56.2026.8.17.0001","title":"Ação","court":"TJPE"}`
	req := httptest.NewRequest(http.MethodPost, "/cases", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	var errResp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	assert.Equal(t, "client not found", errResp.Error)
}

func TestCreateCaseHandler_ShouldReturn400WhenInvalidClientID(t *testing.T) {
	uc := &mockCreateCaseUseCase{
		ExecuteFunc: func(ctx context.Context, input app.CreateCaseInput) (*app.CreateCaseOutput, error) {
			return nil, app.ErrInvalidClientID
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	body := `{"client_id":"invalid","number":"0001234-56.2026.8.17.0001","title":"Ação","court":"TJPE"}`
	req := httptest.NewRequest(http.MethodPost, "/cases", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var errResp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	assert.Equal(t, "invalid client id", errResp.Error)
}

func TestCreateCaseHandler_ShouldReturn400WhenInvalidInput(t *testing.T) {
	uc := &mockCreateCaseUseCase{
		ExecuteFunc: func(ctx context.Context, input app.CreateCaseInput) (*app.CreateCaseOutput, error) {
			return nil, app.ErrInvalidInput
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	body := `{"client_id":"550e8400-e29b-41d4-a716-446655440000","number":"","title":"","court":""}`
	req := httptest.NewRequest(http.MethodPost, "/cases", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateCaseHandler_ShouldReturn400WhenMalformedJSON(t *testing.T) {
	uc := &mockCreateCaseUseCase{}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodPost, "/cases", bytes.NewReader([]byte(`{invalid json`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateCaseHandler_ShouldReturn401WhenMissingJWT(t *testing.T) {
	uc := &mockCreateCaseUseCase{}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return nil, errors.New("invalid token")
		},
	}

	handler := setupHandler(uc, tokenProv)
	body := `{"client_id":"550e8400-e29b-41d4-a716-446655440000","number":"0001234","title":"Ação","court":"TJPE"}`
	req := httptest.NewRequest(http.MethodPost, "/cases", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer invalid.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestCreateCaseHandler_ShouldReturn500OnInternalError(t *testing.T) {
	uc := &mockCreateCaseUseCase{
		ExecuteFunc: func(ctx context.Context, input app.CreateCaseInput) (*app.CreateCaseOutput, error) {
			return nil, errors.New("unexpected error")
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	body := `{"client_id":"550e8400-e29b-41d4-a716-446655440000","number":"0001234","title":"Ação","court":"TJPE"}`
	req := httptest.NewRequest(http.MethodPost, "/cases", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var errResp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	assert.Equal(t, "internal server error", errResp.Error)
}

func TestCreateCaseHandler_ShouldInvokeUseCaseWithCorrectInput(t *testing.T) {
	var capturedInput app.CreateCaseInput
	uc := &mockCreateCaseUseCase{
		ExecuteFunc: func(ctx context.Context, input app.CreateCaseInput) (*app.CreateCaseOutput, error) {
			capturedInput = input
			return &app.CreateCaseOutput{
				ID:     "uuid-456",
				Status: "draft",
			}, nil
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	body := `{"client_id":"550e8400-e29b-41d4-a716-446655440000","number":"0001234-56","title":"Ação de Cobrança","description":"Desc","court":"TJPE"}`
	req := httptest.NewRequest(http.MethodPost, "/cases", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", capturedInput.ClientID)
	assert.Equal(t, "0001234-56", capturedInput.Number)
	assert.Equal(t, "Ação de Cobrança", capturedInput.Title)
	assert.Equal(t, "Desc", capturedInput.Description)
	assert.Equal(t, "TJPE", capturedInput.Court)
}

func TestCreateCaseHandler_ShouldReturn400WhenBodyMissing(t *testing.T) {
	uc := &mockCreateCaseUseCase{}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodPost, "/cases", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
