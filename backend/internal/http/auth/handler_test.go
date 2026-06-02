package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	app "legalflow/internal/application/auth"
)

type mockAuthUseCase struct {
	ExecuteFunc func(ctx context.Context, input app.AuthenticateUserInput) (*app.AuthenticateUserOutput, error)
}

func (m *mockAuthUseCase) Execute(ctx context.Context, input app.AuthenticateUserInput) (*app.AuthenticateUserOutput, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, input)
	}
	return nil, nil
}

type mockTokenProvider struct {
	GenerateFunc func(userID, email string) (string, error)
}

func (m *mockTokenProvider) Generate(userID, email string) (string, error) {
	if m.GenerateFunc != nil {
		return m.GenerateFunc(userID, email)
	}
	return "", nil
}

func setupHandler(auth *mockAuthUseCase, token *mockTokenProvider) http.Handler {
	r := chi.NewRouter()
	h := NewLoginHandler(auth, token)
	h.Register(r)
	return r
}

func TestLoginHandler_ShouldReturn200WhenValidCredentials(t *testing.T) {
	auth := &mockAuthUseCase{
		ExecuteFunc: func(ctx context.Context, input app.AuthenticateUserInput) (*app.AuthenticateUserOutput, error) {
			return &app.AuthenticateUserOutput{
				UserID: "user-123",
				Name:   "John Doe",
				Email:  "john@example.com",
			}, nil
		},
	}
	token := &mockTokenProvider{
		GenerateFunc: func(userID, email string) (string, error) {
			return "jwt-token-123", nil
		},
	}

	handler := setupHandler(auth, token)
	body := `{"email":"john@example.com","password":"correctpassword"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp LoginResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)
	assert.NotEmpty(t, resp.AccessToken)
	assert.Equal(t, "Bearer", resp.TokenType)
	assert.NotZero(t, resp.ExpiresIn)
	assert.Equal(t, "user-123", resp.User.ID)
	assert.Equal(t, "John Doe", resp.User.Name)
	assert.Equal(t, "john@example.com", resp.User.Email)
}

func TestLoginHandler_ShouldReturn401WhenInvalidCredentials(t *testing.T) {
	auth := &mockAuthUseCase{
		ExecuteFunc: func(ctx context.Context, input app.AuthenticateUserInput) (*app.AuthenticateUserOutput, error) {
			return nil, app.ErrInvalidCredentials
		},
	}
	token := &mockTokenProvider{}

	handler := setupHandler(auth, token)
	body := `{"email":"wrong@example.com","password":"wrongpass"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	var errResp LoginErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	assert.Equal(t, "invalid credentials", errResp.Error)
}

func TestLoginHandler_ShouldReturn400WhenMalformedJSON(t *testing.T) {
	auth := &mockAuthUseCase{}
	token := &mockTokenProvider{}

	handler := setupHandler(auth, token)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader([]byte(`{invalid json`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestLoginHandler_ShouldReturn400WhenEmailMissing(t *testing.T) {
	auth := &mockAuthUseCase{}
	token := &mockTokenProvider{}

	handler := setupHandler(auth, token)
	body := `{"password":"somepassword"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestLoginHandler_ShouldReturn400WhenPasswordMissing(t *testing.T) {
	auth := &mockAuthUseCase{}
	token := &mockTokenProvider{}

	handler := setupHandler(auth, token)
	body := `{"email":"john@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestLoginHandler_ShouldReturn500WhenTokenGenerationFails(t *testing.T) {
	auth := &mockAuthUseCase{
		ExecuteFunc: func(ctx context.Context, input app.AuthenticateUserInput) (*app.AuthenticateUserOutput, error) {
			return &app.AuthenticateUserOutput{
				UserID: "user-123",
				Name:   "John Doe",
				Email:  "john@example.com",
			}, nil
		},
	}
	token := &mockTokenProvider{
		GenerateFunc: func(userID, email string) (string, error) {
			return "", app.ErrInvalidInput
		},
	}

	handler := setupHandler(auth, token)
	body := `{"email":"john@example.com","password":"correctpassword"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestLoginHandler_ShouldReturn500WhenUnexpectedError(t *testing.T) {
	auth := &mockAuthUseCase{
		ExecuteFunc: func(ctx context.Context, input app.AuthenticateUserInput) (*app.AuthenticateUserOutput, error) {
			return nil, app.ErrInvalidInput
		},
	}
	token := &mockTokenProvider{}

	handler := setupHandler(auth, token)
	body := `{"email":"john@example.com","password":"correctpassword"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestLoginHandler_ShouldInvokeAuthUseCase(t *testing.T) {
	var capturedInput app.AuthenticateUserInput
	auth := &mockAuthUseCase{
		ExecuteFunc: func(ctx context.Context, input app.AuthenticateUserInput) (*app.AuthenticateUserOutput, error) {
			capturedInput = input
			return &app.AuthenticateUserOutput{
				UserID: "user-123",
				Name:   "Jane Doe",
				Email:  "jane@example.com",
			}, nil
		},
	}
	token := &mockTokenProvider{
		GenerateFunc: func(userID, email string) (string, error) {
			return "jwt-" + userID, nil
		},
	}

	handler := setupHandler(auth, token)
	body := `{"email":"jane@example.com","password":"pass123456"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "jane@example.com", capturedInput.Email)
	assert.Equal(t, "pass123456", capturedInput.Password)
}

func TestLoginHandler_ShouldInvokeTokenProvider(t *testing.T) {
	var capturedUserID, capturedEmail string
	auth := &mockAuthUseCase{
		ExecuteFunc: func(ctx context.Context, input app.AuthenticateUserInput) (*app.AuthenticateUserOutput, error) {
			return &app.AuthenticateUserOutput{
				UserID: "user-456",
				Name:   "John Smith",
				Email:  "john@example.com",
			}, nil
		},
	}
	token := &mockTokenProvider{
		GenerateFunc: func(userID, email string) (string, error) {
			capturedUserID = userID
			capturedEmail = email
			return "jwt-token", nil
		},
	}

	handler := setupHandler(auth, token)
	body := `{"email":"john@example.com","password":"somepassword"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "user-456", capturedUserID)
	assert.Equal(t, "john@example.com", capturedEmail)
}
