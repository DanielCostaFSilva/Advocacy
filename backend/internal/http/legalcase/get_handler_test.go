package legalcase

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
	app "legalflow/internal/application/legalcase"
	"legalflow/internal/middleware"
)

type mockGetCaseByIDUseCase struct {
	ExecuteFunc func(ctx context.Context, input app.GetCaseByIDInput) (*app.GetCaseByIDOutput, error)
}

func (m *mockGetCaseByIDUseCase) Execute(ctx context.Context, input app.GetCaseByIDInput) (*app.GetCaseByIDOutput, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, input)
	}
	return nil, nil
}

func setupGetHandler(uc *mockGetCaseByIDUseCase, tokenProv *mockTokenValidator) http.Handler {
	r := chi.NewRouter()
	h := NewGetHandler(uc)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(tokenProv))
		h.Register(r)
	})
	return r
}

func TestGetCaseByIDHandler_ShouldReturn200(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	caseID := uuid.New().String()
	clientID := uuid.New().String()

	uc := &mockGetCaseByIDUseCase{
		ExecuteFunc: func(ctx context.Context, input app.GetCaseByIDInput) (*app.GetCaseByIDOutput, error) {
			return &app.GetCaseByIDOutput{
				ID:          caseID,
				ClientID:    clientID,
				Number:      "ABC-123",
				Title:       "Test Case",
				Description: "Description",
				Court:       "TJSP",
				Status:      "active",
				CreatedAt:   now,
				UpdatedAt:   now,
			}, nil
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupGetHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/cases/"+caseID, nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp GetCaseResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, caseID, resp.ID)
	assert.Equal(t, clientID, resp.ClientID)
	assert.Equal(t, "ABC-123", resp.Number)
	assert.Equal(t, "Test Case", resp.Title)
	assert.Equal(t, "Description", resp.Description)
	assert.Equal(t, "TJSP", resp.Court)
	assert.Equal(t, "active", resp.Status)
}

func TestGetCaseByIDHandler_ShouldReturn400WhenInvalidUUID(t *testing.T) {
	uc := &mockGetCaseByIDUseCase{
		ExecuteFunc: func(ctx context.Context, input app.GetCaseByIDInput) (*app.GetCaseByIDOutput, error) {
			return nil, app.ErrInvalidCaseID
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupGetHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/cases/invalid-uuid", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	assert.Equal(t, "invalid case id", resp.Error)
}

func TestGetCaseByIDHandler_ShouldReturn404WhenNotFound(t *testing.T) {
	uc := &mockGetCaseByIDUseCase{
		ExecuteFunc: func(ctx context.Context, input app.GetCaseByIDInput) (*app.GetCaseByIDOutput, error) {
			return nil, app.ErrCaseNotFound
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupGetHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/cases/"+uuid.New().String(), nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	var resp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	assert.Equal(t, "case not found", resp.Error)
}

func TestGetCaseByIDHandler_ShouldReturn401WhenUnauthorized(t *testing.T) {
	uc := &mockGetCaseByIDUseCase{}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return nil, errors.New("invalid token")
		},
	}

	handler := setupGetHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/cases/"+uuid.New().String(), nil)
	req.Header.Set("Authorization", "Bearer invalid.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestGetCaseByIDHandler_ShouldReturn500OnError(t *testing.T) {
	uc := &mockGetCaseByIDUseCase{
		ExecuteFunc: func(ctx context.Context, input app.GetCaseByIDInput) (*app.GetCaseByIDOutput, error) {
			return nil, errors.New("unexpected error")
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupGetHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/cases/"+uuid.New().String(), nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
