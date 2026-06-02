package user

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
	app "legalflow/internal/application/user"
)

type mockUseCase struct {
	ExecuteFunc func(ctx context.Context, input app.CreateUserInput) (*app.CreateUserOutput, error)
}

func (m *mockUseCase) Execute(ctx context.Context, input app.CreateUserInput) (*app.CreateUserOutput, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, input)
	}
	return nil, nil
}

func setupHandler(mock *mockUseCase) http.Handler {
	r := chi.NewRouter()
	h := NewHandler(mock)
	h.Register(r)
	return r
}

func TestHandler_CreateUser_ShouldReturn201(t *testing.T) {
	mock := &mockUseCase{
		ExecuteFunc: func(ctx context.Context, input app.CreateUserInput) (*app.CreateUserOutput, error) {
			return &app.CreateUserOutput{
				ID:    "uuid-123",
				Name:  "John Doe",
				Email: "john@example.com",
			}, nil
		},
	}

	handler := setupHandler(mock)
	body := `{"name":"John Doe","email":"john@example.com","password":"12345678"}`
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	var resp CreateUserResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, "uuid-123", resp.ID)
	assert.Equal(t, "John Doe", resp.Name)
	assert.Equal(t, "john@example.com", resp.Email)
}

func TestHandler_CreateUser_ShouldReturn409WhenEmailExists(t *testing.T) {
	mock := &mockUseCase{
		ExecuteFunc: func(ctx context.Context, input app.CreateUserInput) (*app.CreateUserOutput, error) {
			return nil, app.ErrEmailAlreadyExists
		},
	}

	handler := setupHandler(mock)
	body := `{"name":"John Doe","email":"existing@example.com","password":"12345678"}`
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusConflict, rec.Code)

	var errResp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	assert.Equal(t, "email already exists", errResp.Error)
}

func TestHandler_CreateUser_ShouldReturn400WhenInvalidInput(t *testing.T) {
	mock := &mockUseCase{
		ExecuteFunc: func(ctx context.Context, input app.CreateUserInput) (*app.CreateUserOutput, error) {
			return nil, app.ErrInvalidInput
		},
	}

	handler := setupHandler(mock)
	body := `{"name":"","email":"","password":"123"}`
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var errResp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	assert.Equal(t, "invalid input", errResp.Error)
}

func TestHandler_CreateUser_ShouldReturn400WhenMalformedJSON(t *testing.T) {
	mock := &mockUseCase{}
	handler := setupHandler(mock)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(`{invalid json`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandler_CreateUser_ShouldReturn400WhenBodyMissing(t *testing.T) {
	mock := &mockUseCase{}
	handler := setupHandler(mock)

	req := httptest.NewRequest(http.MethodPost, "/users", nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandler_CreateUser_ShouldReturn500OnUnexpectedError(t *testing.T) {
	mock := &mockUseCase{
		ExecuteFunc: func(ctx context.Context, input app.CreateUserInput) (*app.CreateUserOutput, error) {
			return nil, errors.New("unexpected db error")
		},
	}

	handler := setupHandler(mock)
	body := `{"name":"John Doe","email":"john@example.com","password":"12345678"}`
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var errResp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	assert.Equal(t, "internal server error", errResp.Error)
}

func TestHandler_CreateUser_ShouldInvokeUseCaseWithCorrectInput(t *testing.T) {
	var capturedInput app.CreateUserInput
	mock := &mockUseCase{
		ExecuteFunc: func(ctx context.Context, input app.CreateUserInput) (*app.CreateUserOutput, error) {
			capturedInput = input
			return &app.CreateUserOutput{
				ID:    "uuid-123",
				Name:  input.Name,
				Email: input.Email,
			}, nil
		},
	}

	handler := setupHandler(mock)
	body := `{"name":"Jane Doe","email":"jane@example.com","password":"securepass123"}`
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "Jane Doe", capturedInput.Name)
	assert.Equal(t, "jane@example.com", capturedInput.Email)
	assert.Equal(t, "securepass123", capturedInput.Password)
}
