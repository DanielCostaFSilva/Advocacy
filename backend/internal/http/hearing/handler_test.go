package hearing

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
	app "legalflow/internal/application/hearing"
	"legalflow/internal/middleware"
)

type mockCreateHearingUseCase struct {
	ExecuteFunc func(ctx context.Context, input app.CreateHearingInput) (*app.CreateHearingOutput, error)
}

func (m *mockCreateHearingUseCase) Execute(ctx context.Context, input app.CreateHearingInput) (*app.CreateHearingOutput, error) {
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

func setupHandler(uc *mockCreateHearingUseCase, tokenProv *mockTokenValidator) http.Handler {
	r := chi.NewRouter()
	h := NewHandler(uc)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(tokenProv))
		h.Register(r)
	})
	return r
}

func TestHandler_ShouldReturn201(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	caseID := uuid.New().String()
	hearingID := uuid.New().String()

	uc := &mockCreateHearingUseCase{
		ExecuteFunc: func(ctx context.Context, input app.CreateHearingInput) (*app.CreateHearingOutput, error) {
			return &app.CreateHearingOutput{
				ID:          hearingID,
				CaseID:      caseID,
				Title:       "Audiência de Conciliação",
				Description: "Descrição",
				Type:        "conciliation",
				Location:    "Fórum Central",
				ScheduledAt: now,
			}, nil
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	body, _ := json.Marshal(CreateHearingRequest{
		CaseID:      caseID,
		Title:       "Audiência de Conciliação",
		Description: "Descrição",
		Type:        "conciliation",
		Location:    "Fórum Central",
		ScheduledAt: now.Format(time.RFC3339),
	})
	req := httptest.NewRequest(http.MethodPost, "/hearings", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	var resp CreateHearingResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, hearingID, resp.ID)
	assert.Equal(t, caseID, resp.CaseID)
	assert.Equal(t, "conciliation", resp.Type)
}

func TestHandler_ShouldReturn400WhenCaseNotFound(t *testing.T) {
	uc := &mockCreateHearingUseCase{
		ExecuteFunc: func(ctx context.Context, input app.CreateHearingInput) (*app.CreateHearingOutput, error) {
			return nil, app.ErrCaseNotFound
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	body, _ := json.Marshal(CreateHearingRequest{
		CaseID: uuid.New().String(), Title: "Audiência", Type: "conciliation",
		Location: "Local", ScheduledAt: time.Now().Add(48 * time.Hour).Format(time.RFC3339),
	})
	req := httptest.NewRequest(http.MethodPost, "/hearings", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestHandler_ShouldReturn400WhenInvalidCaseID(t *testing.T) {
	uc := &mockCreateHearingUseCase{
		ExecuteFunc: func(ctx context.Context, input app.CreateHearingInput) (*app.CreateHearingOutput, error) {
			return nil, app.ErrInvalidCaseID
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	body, _ := json.Marshal(CreateHearingRequest{
		CaseID: "invalid", Title: "Audiência", Type: "conciliation",
		Location: "Local", ScheduledAt: time.Now().Add(48 * time.Hour).Format(time.RFC3339),
	})
	req := httptest.NewRequest(http.MethodPost, "/hearings", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandler_ShouldReturn400WhenInvalidInput(t *testing.T) {
	uc := &mockCreateHearingUseCase{
		ExecuteFunc: func(ctx context.Context, input app.CreateHearingInput) (*app.CreateHearingOutput, error) {
			return nil, app.ErrInvalidInput
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	body, _ := json.Marshal(CreateHearingRequest{
		CaseID: uuid.New().String(), Title: "AB", Type: "conciliation",
		Location: "Local", ScheduledAt: time.Now().Add(48 * time.Hour).Format(time.RFC3339),
	})
	req := httptest.NewRequest(http.MethodPost, "/hearings", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandler_ShouldReturn400WhenMalformedJSON(t *testing.T) {
	uc := &mockCreateHearingUseCase{}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodPost, "/hearings", bytes.NewReader([]byte("{invalid")))
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandler_ShouldReturn401WhenMissingJWT(t *testing.T) {
	uc := &mockCreateHearingUseCase{}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return nil, errors.New("missing token")
		},
	}

	handler := setupHandler(uc, tokenProv)
	body, _ := json.Marshal(CreateHearingRequest{CaseID: uuid.New().String()})
	req := httptest.NewRequest(http.MethodPost, "/hearings", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestHandler_ShouldReturn401WhenInvalidJWT(t *testing.T) {
	uc := &mockCreateHearingUseCase{}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return nil, errors.New("invalid token")
		},
	}

	handler := setupHandler(uc, tokenProv)
	body, _ := json.Marshal(CreateHearingRequest{CaseID: uuid.New().String()})
	req := httptest.NewRequest(http.MethodPost, "/hearings", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer invalid.token")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestHandler_ShouldReturn500OnInternalError(t *testing.T) {
	uc := &mockCreateHearingUseCase{
		ExecuteFunc: func(ctx context.Context, input app.CreateHearingInput) (*app.CreateHearingOutput, error) {
			return nil, errors.New("unexpected error")
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	body, _ := json.Marshal(CreateHearingRequest{
		CaseID: uuid.New().String(), Title: "Audiência", Type: "conciliation",
		Location: "Local", ScheduledAt: time.Now().Add(48 * time.Hour).Format(time.RFC3339),
	})
	req := httptest.NewRequest(http.MethodPost, "/hearings", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestHandler_ShouldInvokeUseCaseWithCorrectInput(t *testing.T) {
	caseID := uuid.New().String()
	var capturedInput app.CreateHearingInput

	uc := &mockCreateHearingUseCase{
		ExecuteFunc: func(ctx context.Context, input app.CreateHearingInput) (*app.CreateHearingOutput, error) {
			capturedInput = input
			return &app.CreateHearingOutput{
				ID: uuid.New().String(), CaseID: caseID, Title: "Test",
				Type: "conciliation", Location: "Local", ScheduledAt: time.Now(),
			}, nil
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	body, _ := json.Marshal(CreateHearingRequest{
		CaseID: caseID, Title: "Audiência Teste", Description: "Desc",
		Type: "instruction", Location: "Fórum", ScheduledAt: "2026-08-15T14:00:00Z",
	})
	req := httptest.NewRequest(http.MethodPost, "/hearings", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, caseID, capturedInput.CaseID)
	assert.Equal(t, "Audiência Teste", capturedInput.Title)
	assert.Equal(t, "instruction", capturedInput.Type)
}
