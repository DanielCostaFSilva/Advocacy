package postgres

import (
	"context"
	"fmt"
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

func TestHearingRepository_List_ShouldReturnEmpty(t *testing.T) {
	db := connectDB(t)
	hearingRepo := NewHearingRepository(db)
	ctx := context.Background()

	cleanupCases(t, db)
	cleanupClients(t, db)

	hearings, total, err := hearingRepo.List(ctx, domain.ListHearingsParams{
		Limit:  20,
		Offset: 0,
	})
	require.NoError(t, err)
	assert.Empty(t, hearings)
	assert.Equal(t, int64(0), total)
}

func TestHearingRepository_List_ShouldReturnPaginatedResult(t *testing.T) {
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
	for i := 0; i < 5; i++ {
		h, _ := domain.NewHearing(c.ID, fmt.Sprintf("Audiência %d", i+1), "Desc", domain.HearingTypeConciliation, "Local", future.Add(time.Duration(i)*time.Hour))
		require.NoError(t, hearingRepo.Create(ctx, h))
	}

	hearings, total, err := hearingRepo.List(ctx, domain.ListHearingsParams{
		Limit:  2,
		Offset: 0,
	})
	require.NoError(t, err)
	assert.Len(t, hearings, 2)
	assert.Equal(t, int64(5), total)

	hearings, total, err = hearingRepo.List(ctx, domain.ListHearingsParams{
		Limit:  2,
		Offset: 4,
	})
	require.NoError(t, err)
	assert.Len(t, hearings, 1)
	assert.Equal(t, int64(5), total)
}

func TestHearingRepository_List_ShouldFilterByType(t *testing.T) {
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
	h1, _ := domain.NewHearing(c.ID, "Conciliação", "Desc", domain.HearingTypeConciliation, "Local", future)
	h2, _ := domain.NewHearing(c.ID, "Instrução", "Desc", domain.HearingTypeInstruction, "Local", future.Add(1*time.Hour))
	require.NoError(t, hearingRepo.Create(ctx, h1))
	require.NoError(t, hearingRepo.Create(ctx, h2))

	hearings, total, err := hearingRepo.List(ctx, domain.ListHearingsParams{
		Limit:  20,
		Offset: 0,
		Type:   "conciliation",
	})
	require.NoError(t, err)
	assert.Len(t, hearings, 1)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, "Conciliação", hearings[0].Title)
}

func TestHearingRepository_List_ShouldFilterByCaseID(t *testing.T) {
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

	hearings, total, err := hearingRepo.List(ctx, domain.ListHearingsParams{
		Limit:  20,
		Offset: 0,
		CaseID: c1.ID.String(),
	})
	require.NoError(t, err)
	assert.Len(t, hearings, 1)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, "Audiência Case 1", hearings[0].Title)
}

func TestHearingRepository_List_ShouldFilterByDateRange(t *testing.T) {
	db := connectDB(t)
	hearingRepo := NewHearingRepository(db)
	caseRepo := NewCaseRepository(db)
	clientRepo := NewClientRepository(db)
	ctx := context.Background()

	cleanupCases(t, db)
	cleanupClients(t, db)
	client := createTestClient(t, clientRepo, ctx)
	c := createTestCase(t, caseRepo, client.ID, ctx)

	base := time.Now().Add(48 * time.Hour)
	h1, _ := domain.NewHearing(c.ID, "Early", "Desc", domain.HearingTypeConciliation, "Local", base)
	h2, _ := domain.NewHearing(c.ID, "Middle", "Desc", domain.HearingTypeConciliation, "Local", base.Add(24*time.Hour))
	h3, _ := domain.NewHearing(c.ID, "Late", "Desc", domain.HearingTypeConciliation, "Local", base.Add(48*time.Hour))
	require.NoError(t, hearingRepo.Create(ctx, h1))
	require.NoError(t, hearingRepo.Create(ctx, h2))
	require.NoError(t, hearingRepo.Create(ctx, h3))

	startDate := base.Add(12 * time.Hour)
	endDate := base.Add(36 * time.Hour)

	hearings, total, err := hearingRepo.List(ctx, domain.ListHearingsParams{
		Limit:     20,
		Offset:    0,
		StartDate: &startDate,
		EndDate:   &endDate,
	})
	require.NoError(t, err)
	assert.Len(t, hearings, 1)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, "Middle", hearings[0].Title)
}

func TestHearingRepository_List_ShouldSortByTitleDesc(t *testing.T) {
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
	h1, _ := domain.NewHearing(c.ID, "Alpha", "Desc", domain.HearingTypeConciliation, "Local", future)
	h2, _ := domain.NewHearing(c.ID, "Beta", "Desc", domain.HearingTypeConciliation, "Local", future.Add(1*time.Hour))
	require.NoError(t, hearingRepo.Create(ctx, h1))
	require.NoError(t, hearingRepo.Create(ctx, h2))

	hearings, total, err := hearingRepo.List(ctx, domain.ListHearingsParams{
		Limit:  20,
		Offset: 0,
		Sort:   "title",
		Order:  "desc",
	})
	require.NoError(t, err)
	assert.Len(t, hearings, 2)
	assert.Equal(t, int64(2), total)
	assert.Equal(t, "Beta", hearings[0].Title)
	assert.Equal(t, "Alpha", hearings[1].Title)
}

func createTestCase(t *testing.T, repo *CaseRepository, clientID uuid.UUID, ctx context.Context) *legalcaseDomain.Case {
	t.Helper()
	c, err := legalcaseDomain.NewCase(clientID, "HLP-"+uuid.New().String()[:8], "Hearing Helper Case", "Desc", "TJSP")
	require.NoError(t, err)
	require.NoError(t, repo.Create(ctx, c))
	return c
}
