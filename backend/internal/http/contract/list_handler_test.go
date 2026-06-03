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

type mockListContractsUseCase struct {
	ExecuteFunc func(ctx context.Context, input app.ListContractsInput) (*app.ListContractsOutput, error)
}

func (m *mockListContractsUseCase) Execute(ctx context.Context, input app.ListContractsInput) (*app.ListContractsOutput, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, input)
	}
	return nil, nil
}

func setupListHandler(uc *mockListContractsUseCase, tokenProv *mockTokenValidator) http.Handler {
	r := chi.NewRouter()
	h := NewListHandler(uc)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(tokenProv))
		h.Register(r)
	})
	return r
}

func TestListContractsHandler_ShouldReturn200(t *testing.T) {
	uc := &mockListContractsUseCase{
		ExecuteFunc: func(ctx context.Context, input app.ListContractsInput) (*app.ListContractsOutput, error) {
			return &app.ListContractsOutput{
				Data:       []app.ContractDTO{},
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
	req := httptest.NewRequest(http.MethodGet, "/contracts?page=1&page_size=20", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp ListContractsResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, 1, resp.Pagination.Page)
	assert.Equal(t, 20, resp.Pagination.PageSize)
	assert.Equal(t, int64(0), resp.Pagination.Total)
	assert.Equal(t, 0, resp.Pagination.TotalPages)
	assert.Empty(t, resp.Data)
}

func TestListContractsHandler_ShouldReturnData(t *testing.T) {
	contractID := uuid.New().String()
	clientID := uuid.New().String()
	caseID := uuid.New().String()

	uc := &mockListContractsUseCase{
		ExecuteFunc: func(ctx context.Context, input app.ListContractsInput) (*app.ListContractsOutput, error) {
			return &app.ListContractsOutput{
				Data: []app.ContractDTO{
					{
						ID:        contractID,
						ClientID:  clientID,
						CaseID:    caseID,
						Title:     "Contrato de Honorários",
						Type:      "fixed_fee",
						Amount:    "15000.00",
						StartDate: "2026-06-01T00:00:00Z",
						Active:    true,
						CreatedAt: "2026-06-01T10:00:00Z",
					},
				},
				Page:       1,
				PageSize:   20,
				Total:      1,
				TotalPages: 1,
			}, nil
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupListHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/contracts", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp ListContractsResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)
	require.Len(t, resp.Data, 1)
	assert.Equal(t, contractID, resp.Data[0].ID)
	assert.Equal(t, clientID, resp.Data[0].ClientID)
	assert.Equal(t, caseID, resp.Data[0].CaseID)
	assert.Equal(t, "fixed_fee", resp.Data[0].Type)
	assert.Equal(t, "15000.00", resp.Data[0].Amount)
	assert.True(t, resp.Data[0].Active)
	assert.Equal(t, int64(1), resp.Pagination.Total)
	assert.Equal(t, 1, resp.Pagination.TotalPages)
}

func TestListContractsHandler_ShouldPassQueryParams(t *testing.T) {
	var capturedInput app.ListContractsInput
	uc := &mockListContractsUseCase{
		ExecuteFunc: func(ctx context.Context, input app.ListContractsInput) (*app.ListContractsOutput, error) {
			capturedInput = input
			return &app.ListContractsOutput{
				Data:       []app.ContractDTO{},
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
	req := httptest.NewRequest(http.MethodGet, "/contracts?page=2&page_size=10&client_id="+uuid.New().String()+"&type=fixed_fee&active=true&sort=amount&order=desc", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, 2, capturedInput.Page)
	assert.Equal(t, 10, capturedInput.PageSize)
	assert.Equal(t, "fixed_fee", capturedInput.Type)
	assert.Equal(t, "true", capturedInput.Active)
	assert.Equal(t, "amount", capturedInput.Sort)
	assert.Equal(t, "desc", capturedInput.Order)
}

func TestListContractsHandler_ShouldReturn401WhenUnauthorized(t *testing.T) {
	uc := &mockListContractsUseCase{}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return nil, errors.New("invalid token")
		},
	}

	handler := setupListHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/contracts", nil)
	req.Header.Set("Authorization", "Bearer invalid.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestListContractsHandler_ShouldReturn500OnError(t *testing.T) {
	uc := &mockListContractsUseCase{
		ExecuteFunc: func(ctx context.Context, input app.ListContractsInput) (*app.ListContractsOutput, error) {
			return nil, errors.New("unexpected error")
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupListHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/contracts", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
