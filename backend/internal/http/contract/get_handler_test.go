package contract

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
	app "legalflow/internal/application/contract"
	"legalflow/internal/middleware"
)

type mockGetContractUseCase struct {
	ExecuteFunc func(ctx context.Context, input app.GetContractInput) (*app.GetContractOutput, error)
}

func (m *mockGetContractUseCase) Execute(ctx context.Context, input app.GetContractInput) (*app.GetContractOutput, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, input)
	}
	return nil, nil
}

func setupGetHandler(uc *mockGetContractUseCase, tokenProv *mockTokenValidator) http.Handler {
	r := chi.NewRouter()
	h := NewGetHandler(uc)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(tokenProv))
		h.Register(r)
	})
	return r
}

func TestGetContractHandler_ShouldReturn200(t *testing.T) {
	now := time.Now().UTC()
	contractID := uuid.New().String()
	clientID := uuid.New().String()
	caseID := uuid.New().String()

	uc := &mockGetContractUseCase{
		ExecuteFunc: func(ctx context.Context, input app.GetContractInput) (*app.GetContractOutput, error) {
			return &app.GetContractOutput{
				ID:        contractID,
				ClientID:  clientID,
				CaseID:    caseID,
				Title:     "Contrato de Honorários",
				Type:      "fixed_fee",
				Amount:    "15000.00",
				StartDate: "2026-06-01T00:00:00Z",
				Active:    true,
				CreatedAt: now,
				UpdatedAt: now,
			}, nil
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupGetHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/contracts/"+contractID, nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp GetContractResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, contractID, resp.ID)
	assert.Equal(t, clientID, resp.ClientID)
	assert.Equal(t, caseID, resp.CaseID)
	assert.Equal(t, "fixed_fee", resp.Type)
	assert.True(t, resp.Active)
}

func TestGetContractHandler_ShouldReturn400WhenInvalidID(t *testing.T) {
	uc := &mockGetContractUseCase{
		ExecuteFunc: func(ctx context.Context, input app.GetContractInput) (*app.GetContractOutput, error) {
			return nil, app.ErrInvalidContractID
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupGetHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/contracts/invalid-uuid", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var errResp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	assert.Equal(t, "invalid contract id", errResp.Error)
}

func TestGetContractHandler_ShouldReturn404WhenNotFound(t *testing.T) {
	uc := &mockGetContractUseCase{
		ExecuteFunc: func(ctx context.Context, input app.GetContractInput) (*app.GetContractOutput, error) {
			return nil, app.ErrContractNotFound
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupGetHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/contracts/"+uuid.New().String(), nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	var errResp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	assert.Equal(t, "contract not found", errResp.Error)
}

func TestGetContractHandler_ShouldReturn401WhenUnauthorized(t *testing.T) {
	uc := &mockGetContractUseCase{}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return nil, errors.New("invalid token")
		},
	}

	handler := setupGetHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/contracts/"+uuid.New().String(), nil)
	req.Header.Set("Authorization", "Bearer invalid.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestGetContractHandler_ShouldReturn500OnInternalError(t *testing.T) {
	uc := &mockGetContractUseCase{
		ExecuteFunc: func(ctx context.Context, input app.GetContractInput) (*app.GetContractOutput, error) {
			return nil, errors.New("unexpected error")
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupGetHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/contracts/"+uuid.New().String(), nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
