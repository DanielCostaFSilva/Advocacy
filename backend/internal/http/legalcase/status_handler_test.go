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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	app "legalflow/internal/application/legalcase"
	"legalflow/internal/middleware"
)

type mockUpdateCaseStatusUseCase struct {
	ExecuteFunc func(ctx context.Context, input app.UpdateCaseStatusInput) (*app.UpdateCaseStatusOutput, error)
}

func (m *mockUpdateCaseStatusUseCase) Execute(ctx context.Context, input app.UpdateCaseStatusInput) (*app.UpdateCaseStatusOutput, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, input)
	}
	return nil, nil
}

func setupStatusHandler(uc *mockUpdateCaseStatusUseCase, tokenProv *mockTokenValidator) http.Handler {
	r := chi.NewRouter()
	h := NewStatusHandler(uc)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(tokenProv))
		h.Register(r)
	})
	return r
}

func TestStatusHandler_ShouldReturn200(t *testing.T) {
	caseID := uuid.New().String()
	uc := &mockUpdateCaseStatusUseCase{
		ExecuteFunc: func(ctx context.Context, input app.UpdateCaseStatusInput) (*app.UpdateCaseStatusOutput, error) {
			return &app.UpdateCaseStatusOutput{
				ID:     caseID,
				Status: "active",
			}, nil
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupStatusHandler(uc, tokenProv)
	body, _ := json.Marshal(map[string]string{"status": "active"})
	req := httptest.NewRequest(http.MethodPatch, "/cases/"+caseID+"/status", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp UpdateStatusResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, caseID, resp.ID)
	assert.Equal(t, "active", resp.Status)
}

func TestStatusHandler_ShouldReturn400WhenInvalidStatus(t *testing.T) {
	uc := &mockUpdateCaseStatusUseCase{
		ExecuteFunc: func(ctx context.Context, input app.UpdateCaseStatusInput) (*app.UpdateCaseStatusOutput, error) {
			return nil, app.ErrInvalidStatus
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupStatusHandler(uc, tokenProv)
	body, _ := json.Marshal(map[string]string{"status": "invalid"})
	req := httptest.NewRequest(http.MethodPatch, "/cases/"+uuid.New().String()+"/status", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	assert.Equal(t, "invalid status", resp.Error)
}

func TestStatusHandler_ShouldReturn422WhenInvalidTransition(t *testing.T) {
	uc := &mockUpdateCaseStatusUseCase{
		ExecuteFunc: func(ctx context.Context, input app.UpdateCaseStatusInput) (*app.UpdateCaseStatusOutput, error) {
			return nil, app.ErrInvalidStatusTransition
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupStatusHandler(uc, tokenProv)
	body, _ := json.Marshal(map[string]string{"status": "suspended"})
	req := httptest.NewRequest(http.MethodPatch, "/cases/"+uuid.New().String()+"/status", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	var resp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	assert.Equal(t, "invalid status transition", resp.Error)
}

func TestStatusHandler_ShouldReturn404WhenCaseNotFound(t *testing.T) {
	uc := &mockUpdateCaseStatusUseCase{
		ExecuteFunc: func(ctx context.Context, input app.UpdateCaseStatusInput) (*app.UpdateCaseStatusOutput, error) {
			return nil, app.ErrCaseNotFound
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupStatusHandler(uc, tokenProv)
	body, _ := json.Marshal(map[string]string{"status": "active"})
	req := httptest.NewRequest(http.MethodPatch, "/cases/"+uuid.New().String()+"/status", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	var resp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	assert.Equal(t, "case not found", resp.Error)
}

func TestStatusHandler_ShouldReturn401WhenUnauthorized(t *testing.T) {
	uc := &mockUpdateCaseStatusUseCase{}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return nil, errors.New("invalid token")
		},
	}

	handler := setupStatusHandler(uc, tokenProv)
	body, _ := json.Marshal(map[string]string{"status": "active"})
	req := httptest.NewRequest(http.MethodPatch, "/cases/"+uuid.New().String()+"/status", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer invalid.token")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestStatusHandler_ShouldReturn500OnError(t *testing.T) {
	uc := &mockUpdateCaseStatusUseCase{
		ExecuteFunc: func(ctx context.Context, input app.UpdateCaseStatusInput) (*app.UpdateCaseStatusOutput, error) {
			return nil, errors.New("unexpected error")
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupStatusHandler(uc, tokenProv)
	body, _ := json.Marshal(map[string]string{"status": "active"})
	req := httptest.NewRequest(http.MethodPatch, "/cases/"+uuid.New().String()+"/status", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
