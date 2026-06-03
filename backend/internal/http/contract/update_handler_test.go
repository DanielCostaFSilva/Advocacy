package contract

import (
	"bytes"
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
	app "legalflow/internal/application/contract"
	"legalflow/internal/middleware"
)

type mockUpdateContractUseCase struct {
	ExecuteFunc func(ctx context.Context, input app.UpdateContractInput) (*app.UpdateContractOutput, error)
}

func (m *mockUpdateContractUseCase) Execute(ctx context.Context, input app.UpdateContractInput) (*app.UpdateContractOutput, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, input)
	}
	return nil, nil
}

func setupUpdateHandler(uc *mockUpdateContractUseCase, tokenProv *mockTokenValidator) http.Handler {
	r := chi.NewRouter()
	h := NewUpdateHandler(uc)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(tokenProv))
		h.Register(r)
	})
	return r
}

func validUpdatePayload() []byte {
	body := `{
		"title": "Contrato Atualizado",
		"type": "fixed_fee",
		"amount": "20000.00",
		"start_date": "2026-06-01T00:00:00Z",
		"active": true
	}`
	return []byte(body)
}

func TestUpdateContractHandler_ShouldReturn200(t *testing.T) {
	now := time.Now().UTC()
	contractID := uuid.New().String()

	uc := &mockUpdateContractUseCase{
		ExecuteFunc: func(ctx context.Context, input app.UpdateContractInput) (*app.UpdateContractOutput, error) {
			return &app.UpdateContractOutput{
				ID:        contractID,
				ClientID:  uuid.New().String(),
				CaseID:    uuid.New().String(),
				Title:     "Contrato Atualizado",
				Type:      "fixed_fee",
				Amount:    "20000.00",
				StartDate: "2026-06-01T00:00:00Z",
				Active:    true,
				UpdatedAt: now,
			}, nil
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupUpdateHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodPut, "/contracts/"+contractID, bytes.NewReader(validUpdatePayload()))
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp UpdateContractResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, contractID, resp.ID)
	assert.Equal(t, "Contrato Atualizado", resp.Title)
	assert.True(t, resp.Active)
}

func TestUpdateContractHandler_ShouldReturn400WhenInvalidID(t *testing.T) {
	uc := &mockUpdateContractUseCase{
		ExecuteFunc: func(ctx context.Context, input app.UpdateContractInput) (*app.UpdateContractOutput, error) {
			return nil, app.ErrInvalidContractID
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupUpdateHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodPut, "/contracts/invalid-uuid", bytes.NewReader(validUpdatePayload()))
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var errResp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	assert.Equal(t, "invalid contract id", errResp.Error)
}

func TestUpdateContractHandler_ShouldReturn404WhenNotFound(t *testing.T) {
	uc := &mockUpdateContractUseCase{
		ExecuteFunc: func(ctx context.Context, input app.UpdateContractInput) (*app.UpdateContractOutput, error) {
			return nil, app.ErrContractNotFound
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupUpdateHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodPut, "/contracts/"+uuid.New().String(), bytes.NewReader(validUpdatePayload()))
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	var errResp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	assert.Equal(t, "contract not found", errResp.Error)
}

func TestUpdateContractHandler_ShouldReturn400WhenMalformedJSON(t *testing.T) {
	uc := &mockUpdateContractUseCase{}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupUpdateHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodPut, "/contracts/"+uuid.New().String(), bytes.NewReader([]byte("{invalid")))
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateContractHandler_ShouldReturn400WhenInvalidInput(t *testing.T) {
	uc := &mockUpdateContractUseCase{
		ExecuteFunc: func(ctx context.Context, input app.UpdateContractInput) (*app.UpdateContractOutput, error) {
			return nil, app.ErrInvalidInput
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupUpdateHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodPut, "/contracts/"+uuid.New().String(), bytes.NewReader(validUpdatePayload()))
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateContractHandler_ShouldReturn401WhenUnauthorized(t *testing.T) {
	uc := &mockUpdateContractUseCase{}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return nil, errors.New("invalid token")
		},
	}

	handler := setupUpdateHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodPut, "/contracts/"+uuid.New().String(), bytes.NewReader(validUpdatePayload()))
	req.Header.Set("Authorization", "Bearer invalid.token")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestUpdateContractHandler_ShouldReturn500OnInternalError(t *testing.T) {
	uc := &mockUpdateContractUseCase{
		ExecuteFunc: func(ctx context.Context, input app.UpdateContractInput) (*app.UpdateContractOutput, error) {
			return nil, errors.New("unexpected error")
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupUpdateHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodPut, "/contracts/"+uuid.New().String(), bytes.NewReader(validUpdatePayload()))
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
