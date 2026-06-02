package legalcase

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	clientdomain "legalflow/internal/domain/client"
	casedomain "legalflow/internal/domain/legalcase"
)

func TestCreateCase_ShouldSucceedWhenValidInput(t *testing.T) {
	clientID := uuid.New()
	caseRepo := &casedomain.MockCaseRepository{
		FindByNumberFunc: func(ctx context.Context, number string) (*casedomain.Case, error) {
			return nil, nil
		},
		CreateFunc: func(ctx context.Context, c *casedomain.Case) error {
			return nil
		},
	}
	clientRepo := &clientdomain.MockClientRepository{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*clientdomain.Client, error) {
			return &clientdomain.Client{ID: clientID, Name: "Maria"}, nil
		},
	}

	uc := NewCreateCaseUseCase(caseRepo, clientRepo)
	output, err := uc.Execute(context.Background(), CreateCaseInput{
		ClientID:    clientID.String(),
		Number:      "ABC-12345",
		Title:       "Indenizatória",
		Description: "Descrição do caso",
		Court:       "TJSP",
	})

	require.NoError(t, err)
	require.NotNil(t, output)
	assert.NotEmpty(t, output.ID)
	assert.Equal(t, clientID.String(), output.ClientID)
	assert.Equal(t, "ABC-12345", output.Number)
	assert.Equal(t, "Indenizatória", output.Title)
	assert.Equal(t, "Descrição do caso", output.Description)
	assert.Equal(t, "TJSP", output.Court)
	assert.Equal(t, "draft", output.Status)
}

func TestCreateCase_ShouldReturnErrCaseAlreadyExistsWhenDuplicateNumber(t *testing.T) {
	clientID := uuid.New()
	caseRepo := &casedomain.MockCaseRepository{
		FindByNumberFunc: func(ctx context.Context, number string) (*casedomain.Case, error) {
			return &casedomain.Case{Number: "ABC-12345"}, nil
		},
	}
	clientRepo := &clientdomain.MockClientRepository{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*clientdomain.Client, error) {
			return &clientdomain.Client{ID: clientID}, nil
		},
	}

	uc := NewCreateCaseUseCase(caseRepo, clientRepo)
	output, err := uc.Execute(context.Background(), CreateCaseInput{
		ClientID: clientID.String(),
		Number:   "ABC-12345",
		Title:    "Indenizatória",
		Court:    "TJSP",
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, ErrCaseAlreadyExists)
	assert.Nil(t, output)
}

func TestCreateCase_ShouldReturnErrClientNotFoundWhenClientNotExists(t *testing.T) {
	clientID := uuid.New()
	caseRepo := &casedomain.MockCaseRepository{}
	clientRepo := &clientdomain.MockClientRepository{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*clientdomain.Client, error) {
			return nil, nil
		},
	}

	uc := NewCreateCaseUseCase(caseRepo, clientRepo)
	output, err := uc.Execute(context.Background(), CreateCaseInput{
		ClientID: clientID.String(),
		Number:   "ABC-12345",
		Title:    "Indenizatória",
		Court:    "TJSP",
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, ErrClientNotFound)
	assert.Nil(t, output)
}

func TestCreateCase_ShouldReturnErrInvalidClientIDWhenInvalidUUID(t *testing.T) {
	caseRepo := &casedomain.MockCaseRepository{}
	clientRepo := &clientdomain.MockClientRepository{}

	uc := NewCreateCaseUseCase(caseRepo, clientRepo)
	output, err := uc.Execute(context.Background(), CreateCaseInput{
		ClientID: "not-a-uuid",
		Number:   "ABC-12345",
		Title:    "Indenizatória",
		Court:    "TJSP",
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidClientID)
	assert.Nil(t, output)
}

func TestCreateCase_ShouldReturnErrorWhenCaseDataInvalid(t *testing.T) {
	clientID := uuid.New()
	caseRepo := &casedomain.MockCaseRepository{
		FindByNumberFunc: func(ctx context.Context, number string) (*casedomain.Case, error) {
			return nil, nil
		},
	}
	clientRepo := &clientdomain.MockClientRepository{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*clientdomain.Client, error) {
			return &clientdomain.Client{ID: clientID}, nil
		},
	}

	uc := NewCreateCaseUseCase(caseRepo, clientRepo)
	output, err := uc.Execute(context.Background(), CreateCaseInput{
		ClientID: clientID.String(),
		Number:   "AB",
		Title:    "Indenizatória",
		Court:    "TJSP",
	})

	require.Error(t, err)
	assert.Nil(t, output)
}

func TestCreateCase_ShouldReturnErrorWhenRepoCreateFails(t *testing.T) {
	clientID := uuid.New()
	caseRepo := &casedomain.MockCaseRepository{
		FindByNumberFunc: func(ctx context.Context, number string) (*casedomain.Case, error) {
			return nil, nil
		},
		CreateFunc: func(ctx context.Context, c *casedomain.Case) error {
			return errors.New("db error")
		},
	}
	clientRepo := &clientdomain.MockClientRepository{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*clientdomain.Client, error) {
			return &clientdomain.Client{ID: clientID}, nil
		},
	}

	uc := NewCreateCaseUseCase(caseRepo, clientRepo)
	output, err := uc.Execute(context.Background(), CreateCaseInput{
		ClientID: clientID.String(),
		Number:   "ABC-12345",
		Title:    "Indenizatória",
		Court:    "TJSP",
	})

	require.Error(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "db error")
}

func TestCreateCase_ShouldReturnErrorWhenCaseRepoCheckFails(t *testing.T) {
	clientID := uuid.New()
	caseRepo := &casedomain.MockCaseRepository{
		FindByNumberFunc: func(ctx context.Context, number string) (*casedomain.Case, error) {
			return nil, errors.New("search error")
		},
	}
	clientRepo := &clientdomain.MockClientRepository{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*clientdomain.Client, error) {
			return &clientdomain.Client{ID: clientID}, nil
		},
	}

	uc := NewCreateCaseUseCase(caseRepo, clientRepo)
	output, err := uc.Execute(context.Background(), CreateCaseInput{
		ClientID: clientID.String(),
		Number:   "ABC-12345",
		Title:    "Indenizatória",
		Court:    "TJSP",
	})

	require.Error(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "search error")
}

func TestCreateCase_ShouldReturnErrorWhenClientRepoFails(t *testing.T) {
	clientID := uuid.New()
	caseRepo := &casedomain.MockCaseRepository{}
	clientRepo := &clientdomain.MockClientRepository{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*clientdomain.Client, error) {
			return nil, errors.New("client db error")
		},
	}

	uc := NewCreateCaseUseCase(caseRepo, clientRepo)
	output, err := uc.Execute(context.Background(), CreateCaseInput{
		ClientID: clientID.String(),
		Number:   "ABC-12345",
		Title:    "Indenizatória",
		Court:    "TJSP",
	})

	require.Error(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "client db error")
}

func TestCreateCase_ShouldInvokeRepoCreateWithCorrectData(t *testing.T) {
	clientID := uuid.New()
	var capturedCase *casedomain.Case
	caseRepo := &casedomain.MockCaseRepository{
		FindByNumberFunc: func(ctx context.Context, number string) (*casedomain.Case, error) {
			return nil, nil
		},
		CreateFunc: func(ctx context.Context, c *casedomain.Case) error {
			capturedCase = c
			return nil
		},
	}
	clientRepo := &clientdomain.MockClientRepository{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*clientdomain.Client, error) {
			return &clientdomain.Client{ID: clientID}, nil
		},
	}

	uc := NewCreateCaseUseCase(caseRepo, clientRepo)
	output, err := uc.Execute(context.Background(), CreateCaseInput{
		ClientID:    clientID.String(),
		Number:      "ABC-12345",
		Title:       "Indenizatória",
		Description: "Desc",
		Court:       "TJSP",
	})

	require.NoError(t, err)
	require.NotNil(t, output)
	require.NotNil(t, capturedCase)
	assert.Equal(t, "ABC-12345", capturedCase.Number)
	assert.Equal(t, "Indenizatória", capturedCase.Title)
	assert.Equal(t, "Desc", capturedCase.Description)
	assert.Equal(t, "TJSP", capturedCase.Court)
	assert.Equal(t, clientID, capturedCase.ClientID)
	assert.Equal(t, casedomain.CaseStatusDraft, capturedCase.Status)
}

func TestCreateCase_ShouldReturnStatusDraft(t *testing.T) {
	clientID := uuid.New()
	caseRepo := &casedomain.MockCaseRepository{
		FindByNumberFunc: func(ctx context.Context, number string) (*casedomain.Case, error) {
			return nil, nil
		},
		CreateFunc: func(ctx context.Context, c *casedomain.Case) error {
			return nil
		},
	}
	clientRepo := &clientdomain.MockClientRepository{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (*clientdomain.Client, error) {
			return &clientdomain.Client{ID: clientID}, nil
		},
	}

	uc := NewCreateCaseUseCase(caseRepo, clientRepo)
	output, err := uc.Execute(context.Background(), CreateCaseInput{
		ClientID: clientID.String(),
		Number:   "ABC-12345",
		Title:    "Indenizatória",
		Court:    "TJSP",
	})

	require.NoError(t, err)
	require.NotNil(t, output)
	assert.Equal(t, "draft", output.Status)
}
