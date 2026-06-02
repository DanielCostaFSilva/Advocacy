package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errInvalidToken = errors.New("invalid token")
var errExpiredToken = errors.New("expired token")

type mockTokenProvider struct {
	ValidateFunc func(token string) (*TokenClaims, error)
}

func (m *mockTokenProvider) Validate(token string) (*TokenClaims, error) {
	if m.ValidateFunc != nil {
		return m.ValidateFunc(token)
	}
	return nil, nil
}

func TestAuthMiddleware_ShouldAllowWhenValidToken(t *testing.T) {
	provider := &mockTokenProvider{
		ValidateFunc: func(token string) (*TokenClaims, error) {
			return &TokenClaims{
				UserID:    "user-123",
				Email:     "john@example.com",
				IssuedAt:  time.Now(),
				ExpiresAt: time.Now().Add(time.Hour),
			}, nil
		},
	}

	handler := AuthMiddleware(provider)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "user-123", GetUserID(r.Context()))
		assert.Equal(t, "john@example.com", GetUserEmail(r.Context()))
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAuthMiddleware_ShouldReturn401WhenMissingHeader(t *testing.T) {
	provider := &mockTokenProvider{}

	handler := AuthMiddleware(provider)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	var body map[string]string
	json.NewDecoder(rec.Body).Decode(&body)
	assert.Equal(t, "unauthorized", body["error"])
}

func TestAuthMiddleware_ShouldReturn401WhenMalformedHeader(t *testing.T) {
	provider := &mockTokenProvider{}

	handler := AuthMiddleware(provider)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthMiddleware_ShouldReturn401WhenInvalidToken(t *testing.T) {
	provider := &mockTokenProvider{
		ValidateFunc: func(token string) (*TokenClaims, error) {
			return nil, errInvalidToken
		},
	}

	handler := AuthMiddleware(provider)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer invalid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthMiddleware_ShouldReturn401WhenExpiredToken(t *testing.T) {
	provider := &mockTokenProvider{
		ValidateFunc: func(token string) (*TokenClaims, error) {
			return nil, errExpiredToken
		},
	}

	handler := AuthMiddleware(provider)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer expired.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthMiddleware_ShouldPutClaimsInContext(t *testing.T) {
	claims := &TokenClaims{
		UserID:    "user-456",
		Email:     "jane@example.com",
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(time.Hour),
	}
	provider := &mockTokenProvider{
		ValidateFunc: func(token string) (*TokenClaims, error) {
			return claims, nil
		},
	}

	var capturedClaims *TokenClaims
	handler := AuthMiddleware(provider)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedClaims = GetClaims(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.NotNil(t, capturedClaims)
	assert.Equal(t, "user-456", capturedClaims.UserID)
	assert.Equal(t, "jane@example.com", capturedClaims.Email)
}
