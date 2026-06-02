package hearing

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
	app "legalflow/internal/application/hearing"
	"legalflow/internal/middleware"
)

type mockListHearingsUseCase struct {
	ExecuteFunc func(ctx context.Context, input app.ListHearingsInput) (*app.ListHearingsOutput, error)
}

func (m *mockListHearingsUseCase) Execute(ctx context.Context, input app.ListHearingsInput) (*app.ListHearingsOutput, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, input)
	}
	return nil, nil
}

func setupListHandler(uc *mockListHearingsUseCase, tokenProv *mockTokenValidator) http.Handler {
	r := chi.NewRouter()
	h := NewListHandler(uc)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(tokenProv))
		h.Register(r)
	})
	return r
}

func TestListHearingsHandler_ShouldReturn200(t *testing.T) {
	uc := &mockListHearingsUseCase{
		ExecuteFunc: func(ctx context.Context, input app.ListHearingsInput) (*app.ListHearingsOutput, error) {
			return &app.ListHearingsOutput{
				Data:       []app.HearingDTO{},
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
	req := httptest.NewRequest(http.MethodGet, "/hearings?page=1&page_size=20", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp ListHearingsResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, 1, resp.Pagination.Page)
	assert.Equal(t, 20, resp.Pagination.PageSize)
	assert.Equal(t, int64(0), resp.Pagination.Total)
	assert.Equal(t, 0, resp.Pagination.TotalPages)
	assert.Empty(t, resp.Data)
}

func TestListHearingsHandler_ShouldPassQueryParams(t *testing.T) {
	var capturedInput app.ListHearingsInput
	uc := &mockListHearingsUseCase{
		ExecuteFunc: func(ctx context.Context, input app.ListHearingsInput) (*app.ListHearingsOutput, error) {
			capturedInput = input
			return &app.ListHearingsOutput{
				Data:       []app.HearingDTO{},
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
	req := httptest.NewRequest(http.MethodGet, "/hearings?page=2&page_size=10&case_id="+uuid.New().String()+"&type=conciliation&sort=scheduled_at&order=desc", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, 2, capturedInput.Page)
	assert.Equal(t, 10, capturedInput.PageSize)
	assert.Equal(t, "conciliation", capturedInput.Type)
	assert.Equal(t, "scheduled_at", capturedInput.Sort)
	assert.Equal(t, "desc", capturedInput.Order)
}

func TestListHearingsHandler_ShouldReturn401WhenUnauthorized(t *testing.T) {
	uc := &mockListHearingsUseCase{}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return nil, errors.New("invalid token")
		},
	}

	handler := setupListHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/hearings", nil)
	req.Header.Set("Authorization", "Bearer invalid.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestListHearingsHandler_ShouldReturn500OnError(t *testing.T) {
	uc := &mockListHearingsUseCase{
		ExecuteFunc: func(ctx context.Context, input app.ListHearingsInput) (*app.ListHearingsOutput, error) {
			return nil, errors.New("unexpected error")
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupListHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodGet, "/hearings", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestListHearingsHandler_ShouldReturnData(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	hearingID := uuid.New().String()
	caseID := uuid.New().String()

	uc := &mockListHearingsUseCase{
		ExecuteFunc: func(ctx context.Context, input app.ListHearingsInput) (*app.ListHearingsOutput, error) {
			return &app.ListHearingsOutput{
				Data: []app.HearingDTO{
					{
						ID:          hearingID,
						CaseID:      caseID,
						Title:       "Audiência de Conciliação",
						Type:        "conciliation",
						Location:    "Fórum Central",
						ScheduledAt: now,
						CreatedAt:   now,
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
	req := httptest.NewRequest(http.MethodGet, "/hearings", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp ListHearingsResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)
	require.Len(t, resp.Data, 1)
	assert.Equal(t, hearingID, resp.Data[0].ID)
	assert.Equal(t, caseID, resp.Data[0].CaseID)
	assert.Equal(t, "conciliation", resp.Data[0].Type)
	assert.Equal(t, int64(1), resp.Pagination.Total)
	assert.Equal(t, 1, resp.Pagination.TotalPages)
}
