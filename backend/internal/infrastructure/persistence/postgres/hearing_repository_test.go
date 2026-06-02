package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domain "legalflow/internal/domain/hearing"
	legalcaseDomain "legalflow/internal/domain/legalcase"
)

func TestHearingRepository_Create_ShouldSaveSuccessfully(t *testing.T) {
	db := connectDB(t)
	hearingRepo := NewHearingRepository(db)
	caseRepo := NewCaseRepository(db)
	clientRepo := NewClientRepository(db)
	ctx := context.Background()

	cleanupCases(t, db)
	cleanupClients(t, db)
	client := createTestClient(t, clientRepo, ctx)
	c := createTestCase(t, caseRepo, client.ID, ctx)

	future := time.Now().Add(48 * time.Hour)
	h, err := domain.NewHearing(c.ID, "Audiência de Conciliação", "Descrição", domain.HearingTypeConciliation, "Fórum Central", future)
	require.NoError(t, err)

	err = hearingRepo.Create(ctx, h)
	require.NoError(t, err)
	assert.NotEmpty(t, h.ID)
	assert.False(t, h.CreatedAt.IsZero())
	assert.False(t, h.UpdatedAt.IsZero())
}

func TestHearingRepository_FindByID_ShouldReturnHearing(t *testing.T) {
	db := connectDB(t)
	hearingRepo := NewHearingRepository(db)
	caseRepo := NewCaseRepository(db)
	clientRepo := NewClientRepository(db)
	ctx := context.Background()

	cleanupCases(t, db)
	cleanupClients(t, db)
	client := createTestClient(t, clientRepo, ctx)
	c := createTestCase(t, caseRepo, client.ID, ctx)

	future := time.Now().Add(48 * time.Hour)
	original, err := domain.NewHearing(c.ID, "Audiência Instrução", "Desc", domain.HearingTypeInstruction, "Fórum", future)
	require.NoError(t, err)
	require.NoError(t, hearingRepo.Create(ctx, original))

	found, err := hearingRepo.FindByID(ctx, original.ID)
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, original.ID, found.ID)
	assert.Equal(t, "Audiência Instrução", found.Title)
	assert.Equal(t, domain.HearingTypeInstruction, found.Type)
}

func TestHearingRepository_ListByCaseID_ShouldReturnHearings(t *testing.T) {
	db := connectDB(t)
	hearingRepo := NewHearingRepository(db)
	caseRepo := NewCaseRepository(db)
	clientRepo := NewClientRepository(db)
	ctx := context.Background()

	cleanupCases(t, db)
	cleanupClients(t, db)
	client := createTestClient(t, clientRepo, ctx)
	c := createTestCase(t, caseRepo, client.ID, ctx)

	future := time.Now().Add(48 * time.Hour)
	h1, _ := domain.NewHearing(c.ID, "Primeira Audiência", "Desc", domain.HearingTypeConciliation, "Local A", future)
	h2, _ := domain.NewHearing(c.ID, "Segunda Audiência", "Desc", domain.HearingTypeJudgment, "Local B", future.Add(2*time.Hour))
	require.NoError(t, hearingRepo.Create(ctx, h1))
	require.NoError(t, hearingRepo.Create(ctx, h2))

	hearings, err := hearingRepo.ListByCaseID(ctx, c.ID)
	require.NoError(t, err)
	assert.Len(t, hearings, 2)
}

func TestHearingRepository_ListByCaseID_ShouldReturnOnlyMatchingCase(t *testing.T) {
	db := connectDB(t)
	hearingRepo := NewHearingRepository(db)
	caseRepo := NewCaseRepository(db)
	clientRepo := NewClientRepository(db)
	ctx := context.Background()

	cleanupCases(t, db)
	cleanupClients(t, db)
	client := createTestClient(t, clientRepo, ctx)
	c1 := createTestCase(t, caseRepo, client.ID, ctx)
	c2 := createTestCase(t, caseRepo, client.ID, ctx)

	future := time.Now().Add(48 * time.Hour)
	h1, _ := domain.NewHearing(c1.ID, "Audiência Case 1", "Desc", domain.HearingTypeVirtual, "Online", future)
	h2, _ := domain.NewHearing(c2.ID, "Audiência Case 2", "Desc", domain.HearingTypeVirtual, "Online", future)
	require.NoError(t, hearingRepo.Create(ctx, h1))
	require.NoError(t, hearingRepo.Create(ctx, h2))

	hearings, err := hearingRepo.ListByCaseID(ctx, c1.ID)
	require.NoError(t, err)
	assert.Len(t, hearings, 1)
	assert.Equal(t, "Audiência Case 1", hearings[0].Title)
}

func TestHearingRepository_FindByID_ShouldReturnNilWhenNotFound(t *testing.T) {
	db := connectDB(t)
	hearingRepo := NewHearingRepository(db)
	ctx := context.Background()

	cleanupCases(t, db)

	found, err := hearingRepo.FindByID(ctx, uuid.New())
	require.NoError(t, err)
	assert.Nil(t, found)
}

func createTestCase(t *testing.T, repo *CaseRepository, clientID uuid.UUID, ctx context.Context) *legalcaseDomain.Case {
	t.Helper()
	c, err := legalcaseDomain.NewCase(clientID, "HLP-"+uuid.New().String()[:8], "Hearing Helper Case", "Desc", "TJSP")
	require.NoError(t, err)
	require.NoError(t, repo.Create(ctx, c))
	return c
}
