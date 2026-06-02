package timeline

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
	app "legalflow/internal/application/timeline"
	"legalflow/internal/middleware"
)

type mockGetCaseTimelineUseCase struct {
	ExecuteFunc func(ctx context.Context, input app.GetCaseTimelineInput) (*app.GetCaseTimelineOutput, error)
}

func (m *mockGetCaseTimelineUseCase) Execute(ctx context.Context, input app.GetCaseTimelineInput) (*app.GetCaseTimelineOutput, error) {
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

func setupHandler(uc *mockGetCaseTimelineUseCase, tokenProv *mockTokenValidator) http.Handler {
	r := chi.NewRouter()
	h := NewHandler(uc)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(tokenProv))
		h.Register(r)
	})
	return r
}

func TestTimelineHandler_ShouldReturn200(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	caseID := uuid.New().String()

	uc := &mockGetCaseTimelineUseCase{
		ExecuteFunc: func(ctx context.Context, input app.GetCaseTimelineInput) (*app.GetCaseTimelineOutput, error) {
			return &app.GetCaseTimelineOutput{
				Events: []app.TimelineEventDTO{
					{ID: uuid.New().String(), Type: "CASE_CREATED", Description: "Processo criado", CreatedAt: now},
				},
			}, nil
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/cases/"+caseID+"/timeline", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp TimelineResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Len(t, resp.Events, 1)
	assert.Equal(t, "CASE_CREATED", resp.Events[0].Type)
}

func TestTimelineHandler_ShouldIncludeMetadata(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	caseID := uuid.New().String()
	meta := json.RawMessage(`{"hearing_id":"abc","hearing_type":"conciliation"}`)

	uc := &mockGetCaseTimelineUseCase{
		ExecuteFunc: func(ctx context.Context, input app.GetCaseTimelineInput) (*app.GetCaseTimelineOutput, error) {
			return &app.GetCaseTimelineOutput{
				Events: []app.TimelineEventDTO{
					{ID: uuid.New().String(), Type: "HEARING_CREATED", Description: "Audiência agendada", Metadata: &meta, CreatedAt: now},
				},
			}, nil
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/cases/"+caseID+"/timeline", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp TimelineResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)
	require.Len(t, resp.Events, 1)
	require.NotNil(t, resp.Events[0].Metadata)
	var metaMap map[string]any
	err = json.Unmarshal(*resp.Events[0].Metadata, &metaMap)
	require.NoError(t, err)
	assert.Equal(t, "abc", metaMap["hearing_id"])
}

func TestTimelineHandler_ShouldReturn400WhenInvalidCaseID(t *testing.T) {
	uc := &mockGetCaseTimelineUseCase{
		ExecuteFunc: func(ctx context.Context, input app.GetCaseTimelineInput) (*app.GetCaseTimelineOutput, error) {
			return nil, app.ErrInvalidCaseID
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/cases/invalid-uuid/timeline", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestTimelineHandler_ShouldReturn404WhenCaseNotFound(t *testing.T) {
	uc := &mockGetCaseTimelineUseCase{
		ExecuteFunc: func(ctx context.Context, input app.GetCaseTimelineInput) (*app.GetCaseTimelineOutput, error) {
			return nil, app.ErrCaseNotFound
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/cases/"+uuid.New().String()+"/timeline", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestTimelineHandler_ShouldReturn401WhenUnauthorized(t *testing.T) {
	uc := &mockGetCaseTimelineUseCase{}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return nil, errors.New("invalid token")
		},
	}

	handler := setupHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/cases/"+uuid.New().String()+"/timeline", nil)
	req.Header.Set("Authorization", "Bearer invalid.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
