package contract

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
	app "legalflow/internal/application/contract"
	"legalflow/internal/middleware"
)

type mockCreateContractUseCase struct {
	ExecuteFunc func(ctx context.Context, input app.CreateContractInput) (*app.CreateContractOutput, error)
}

func (m *mockCreateContractUseCase) Execute(ctx context.Context, input app.CreateContractInput) (*app.CreateContractOutput, error) {
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

func setupHandler(uc *mockCreateContractUseCase, tokenProv *mockTokenValidator) http.Handler {
	r := chi.NewRouter()
	h := NewCreateHandler(uc)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(tokenProv))
		h.Register(r)
	})
	return r
}

func validPayload() []byte {
	body := `{
		"client_id": "550e8400-e29b-41d4-a716-446655440000",
		"case_id": "660e8400-e29b-41d4-a716-446655440001",
		"title": "Contrato de Honorários",
		"type": "fixed_fee",
		"amount": "15000.00",
		"start_date": "2026-06-01T00:00:00Z"
	}`
	return []byte(body)
}

func TestCreateContractHandler_ShouldReturn201(t *testing.T) {
	uc := &mockCreateContractUseCase{
		ExecuteFunc: func(ctx context.Context, input app.CreateContractInput) (*app.CreateContractOutput, error) {
			return &app.CreateContractOutput{
				ID:        "new-uuid",
				ClientID:  input.ClientID,
				CaseID:    input.CaseID,
				Title:     input.Title,
				Type:      input.Type,
				Amount:    input.Amount,
				StartDate: input.StartDate,
				Active:    true,
			}, nil
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodPost, "/contracts", bytes.NewReader(validPayload()))
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	var resp CreateContractResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, "new-uuid", resp.ID)
	assert.Equal(t, "Contrato de Honorários", resp.Title)
	assert.Equal(t, "fixed_fee", resp.Type)
	assert.True(t, resp.Active)
}

func TestCreateContractHandler_ShouldReturn400WhenMalformedJSON(t *testing.T) {
	uc := &mockCreateContractUseCase{}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodPost, "/contracts", bytes.NewReader([]byte("{invalid")))
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateContractHandler_ShouldReturn400WhenMissingClientID(t *testing.T) {
	uc := &mockCreateContractUseCase{}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	body := `{"case_id": "uuid", "title": "Test", "type": "fixed_fee", "amount": "100", "start_date": "2026-06-01T00:00:00Z"}`
	req := httptest.NewRequest(http.MethodPost, "/contracts", bytes.NewReader([]byte(body)))
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var errResp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	assert.Equal(t, "client_id is required", errResp.Error)
}

func TestCreateContractHandler_ShouldReturn400WhenMissingCaseID(t *testing.T) {
	uc := &mockCreateContractUseCase{}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	body := `{"client_id": "uuid", "title": "Test", "type": "fixed_fee", "amount": "100", "start_date": "2026-06-01T00:00:00Z"}`
	req := httptest.NewRequest(http.MethodPost, "/contracts", bytes.NewReader([]byte(body)))
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var errResp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	assert.Contains(t, errResp.Error, "case_id")
}

func TestCreateContractHandler_ShouldReturn400WhenMissingRequiredFields(t *testing.T) {
	uc := &mockCreateContractUseCase{}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	body := `{"client_id": "uuid", "case_id": "uuid"}`
	req := httptest.NewRequest(http.MethodPost, "/contracts", bytes.NewReader([]byte(body)))
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateContractHandler_ShouldReturn400WhenErrInvalidClientID(t *testing.T) {
	uc := &mockCreateContractUseCase{
		ExecuteFunc: func(ctx context.Context, input app.CreateContractInput) (*app.CreateContractOutput, error) {
			return nil, app.ErrInvalidClientID
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodPost, "/contracts", bytes.NewReader(validPayload()))
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var errResp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	assert.Equal(t, "invalid client id", errResp.Error)
}

func TestCreateContractHandler_ShouldReturn400WhenErrInvalidCaseID(t *testing.T) {
	uc := &mockCreateContractUseCase{
		ExecuteFunc: func(ctx context.Context, input app.CreateContractInput) (*app.CreateContractOutput, error) {
			return nil, app.ErrInvalidCaseID
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodPost, "/contracts", bytes.NewReader(validPayload()))
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var errResp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	assert.Equal(t, "invalid case id", errResp.Error)
}

func TestCreateContractHandler_ShouldReturn400WhenErrInvalidInput(t *testing.T) {
	uc := &mockCreateContractUseCase{
		ExecuteFunc: func(ctx context.Context, input app.CreateContractInput) (*app.CreateContractOutput, error) {
			return nil, app.ErrInvalidInput
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodPost, "/contracts", bytes.NewReader(validPayload()))
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateContractHandler_ShouldReturn404WhenClientNotFound(t *testing.T) {
	uc := &mockCreateContractUseCase{
		ExecuteFunc: func(ctx context.Context, input app.CreateContractInput) (*app.CreateContractOutput, error) {
			return nil, app.ErrClientNotFound
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodPost, "/contracts", bytes.NewReader(validPayload()))
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	var errResp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	assert.Equal(t, "client not found", errResp.Error)
}

func TestCreateContractHandler_ShouldReturn404WhenCaseNotFound(t *testing.T) {
	uc := &mockCreateContractUseCase{
		ExecuteFunc: func(ctx context.Context, input app.CreateContractInput) (*app.CreateContractOutput, error) {
			return nil, app.ErrCaseNotFound
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodPost, "/contracts", bytes.NewReader(validPayload()))
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	var errResp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	assert.Equal(t, "case not found", errResp.Error)
}

func TestCreateContractHandler_ShouldReturn401WhenMissingJWT(t *testing.T) {
	uc := &mockCreateContractUseCase{}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return nil, errors.New("missing token")
		},
	}

	handler := setupHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodPost, "/contracts", bytes.NewReader(validPayload()))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestCreateContractHandler_ShouldReturn500OnInternalError(t *testing.T) {
	uc := &mockCreateContractUseCase{
		ExecuteFunc: func(ctx context.Context, input app.CreateContractInput) (*app.CreateContractOutput, error) {
			return nil, errors.New("unexpected")
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	req := httptest.NewRequest(http.MethodPost, "/contracts", bytes.NewReader(validPayload()))
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestCreateContractHandler_ShouldInvokeUseCaseWithCorrectInput(t *testing.T) {
	var capturedInput app.CreateContractInput

	uc := &mockCreateContractUseCase{
		ExecuteFunc: func(ctx context.Context, input app.CreateContractInput) (*app.CreateContractOutput, error) {
			capturedInput = input
			return &app.CreateContractOutput{
				ID: "new-uuid", ClientID: input.ClientID, CaseID: input.CaseID,
				Title: input.Title, Type: input.Type, Amount: input.Amount,
				StartDate: input.StartDate, Active: true,
			}, nil
		},
	}
	tokenProv := &mockTokenValidator{
		ValidateFunc: func(token string) (*middleware.TokenClaims, error) {
			return &middleware.TokenClaims{UserID: "user-123"}, nil
		},
	}

	handler := setupHandler(uc, tokenProv)
	payload := `{
		"client_id": "550e8400-e29b-41d4-a716-446655440000",
		"case_id": "660e8400-e29b-41d4-a716-446655440001",
		"title": "Contrato Teste",
		"description": "Desc",
		"type": "monthly",
		"amount": "5000.00",
		"start_date": "2026-06-01T00:00:00Z",
		"end_date": "2027-06-01T00:00:00Z"
	}`
	req := httptest.NewRequest(http.MethodPost, "/contracts", bytes.NewReader([]byte(payload)))
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", capturedInput.ClientID)
	assert.Equal(t, "660e8400-e29b-41d4-a716-446655440001", capturedInput.CaseID)
	assert.Equal(t, "Contrato Teste", capturedInput.Title)
	assert.Equal(t, "Desc", capturedInput.Description)
	assert.Equal(t, "monthly", capturedInput.Type)
	assert.Equal(t, "5000.00", capturedInput.Amount)
	assert.Equal(t, "2026-06-01T00:00:00Z", capturedInput.StartDate)
	require.NotNil(t, capturedInput.EndDate)
	assert.Equal(t, "2027-06-01T00:00:00Z", *capturedInput.EndDate)
}
