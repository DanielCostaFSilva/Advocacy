package auth

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
	app "legalflow/internal/application/auth"
	"legalflow/internal/middleware"
)

type mockGetCurrentUserUseCase struct {
	ExecuteFunc func(ctx context.Context, input app.GetCurrentUserInput) (*app.GetCurrentUserOutput, error)
}

func (m *mockGetCurrentUserUseCase) Execute(ctx context.Context, input app.GetCurrentUserInput) (*app.GetCurrentUserOutput, error) {
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

func setupMeHandler(useCase *mockGetCurrentUserUseCase, tokenProv *mockTokenValidator) http.Handler {
	r := chi.NewRouter()
	h := NewMeHandler(useCase)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(tokenProv))
		r.Get("/me", h.Me)
	})
	return r
}

func TestMeHandler_ShouldReturn200WhenAuthenticated(t *testing.T) {
	uc := &mockGetCurrentUserUseCase{
		ExecuteFunc: func(ctx context.Context, input app.GetCurrentUserInput) (*app.GetCurrentUserOutput, error) {
			return &app.GetCurrentUserOutput{
				ID:    "user-123",
				Name:  "John Doe",
				Email: "john@example.com",
			}, nil
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123", Email: "john@example.com"}, nil
		},
	}

	handler := setupMeHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp app.GetCurrentUserOutput
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, "user-123", resp.ID)
	assert.Equal(t, "John Doe", resp.Name)
	assert.Equal(t, "john@example.com", resp.Email)
}

func TestMeHandler_ShouldReturn401WhenMissingClaims(t *testing.T) {
	uc := &mockGetCurrentUserUseCase{}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return nil, errors.New("invalid token")
		},
	}

	handler := setupMeHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer invalid.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestMeHandler_ShouldReturn404WhenUserNotFound(t *testing.T) {
	uc := &mockGetCurrentUserUseCase{
		ExecuteFunc: func(ctx context.Context, input app.GetCurrentUserInput) (*app.GetCurrentUserOutput, error) {
			return nil, app.ErrUserNotFound
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-999"}, nil
		},
	}

	handler := setupMeHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	var body map[string]string
	json.NewDecoder(rec.Body).Decode(&body)
	assert.Equal(t, "user not found", body["error"])
}

func TestMeHandler_ShouldReturn500WhenUseCaseFails(t *testing.T) {
	uc := &mockGetCurrentUserUseCase{
		ExecuteFunc: func(ctx context.Context, input app.GetCurrentUserInput) (*app.GetCurrentUserOutput, error) {
			return nil, errors.New("unexpected error")
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupMeHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var body map[string]string
	json.NewDecoder(rec.Body).Decode(&body)
	assert.Equal(t, "internal server error", body["error"])
}

func TestMeHandler_ShouldInvokeUseCaseWithCorrectUserID(t *testing.T) {
	var capturedInput app.GetCurrentUserInput
	uc := &mockGetCurrentUserUseCase{
		ExecuteFunc: func(ctx context.Context, input app.GetCurrentUserInput) (*app.GetCurrentUserOutput, error) {
			capturedInput = input
			return &app.GetCurrentUserOutput{
				ID:    "user-456",
				Name:  "Jane Doe",
				Email: "jane@example.com",
			}, nil
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-456"}, nil
		},
	}

	handler := setupMeHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "user-456", capturedInput.UserID)
}
