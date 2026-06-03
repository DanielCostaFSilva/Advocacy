package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domain "legalflow/internal/domain/legalcase"
	clientdomain "legalflow/internal/domain/client"
)

var cpfCounter int64

func TestCaseRepository_Create_ShouldSaveSuccessfully(t *testing.T) {
	db := connectDB(t)
	caseRepo := NewCaseRepository(db)
	clientRepo := NewClientRepository(db)
	ctx := context.Background()

	cleanupCases(t, db)
	cleanupClients(t, db)
	client := createTestClient(t, clientRepo, ctx)

	c, err := domain.NewCase(client.ID, "ABC-12345", "Indenizatória", "Descrição", "TJSP")
	require.NoError(t, err)

	err = caseRepo.Create(ctx, c)
	require.NoError(t, err)
	assert.NotEmpty(t, c.ID)
	assert.Equal(t, domain.CaseStatusDraft, c.Status)
}

func TestCaseRepository_Create_ShouldReturnDuplicateError(t *testing.T) {
	db := connectDB(t)
	caseRepo := NewCaseRepository(db)
	clientRepo := NewClientRepository(db)
	ctx := context.Background()

	cleanupCases(t, db)
	cleanupClients(t, db)
	client := createTestClient(t, clientRepo, ctx)

	c1, err := domain.NewCase(client.ID, "UNIQ-12345", "Title One", "Desc", "TJSP")
	require.NoError(t, err)
	err = caseRepo.Create(ctx, c1)
	require.NoError(t, err)

	c2, err := domain.NewCase(client.ID, "UNIQ-12345", "Title Two", "Desc", "TJSP")
	require.NoError(t, err)
	err = caseRepo.Create(ctx, c2)
	require.Error(t, err)
}

func TestCaseRepository_FindByID_ShouldReturnCase(t *testing.T) {
	db := connectDB(t)
	caseRepo := NewCaseRepository(db)
	clientRepo := NewClientRepository(db)
	ctx := context.Background()

	cleanupCases(t, db)
	cleanupClients(t, db)
	client := createTestClient(t, clientRepo, ctx)

	original, err := domain.NewCase(client.ID, "FIND-12345", "Find By ID", "Desc", "TJSP")
	require.NoError(t, err)
	err = caseRepo.Create(ctx, original)
	require.NoError(t, err)

	found, err := caseRepo.FindByID(ctx, original.ID)
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, original.ID, found.ID)
	assert.Equal(t, "FIND-12345", found.Number)
	assert.Equal(t, "Find By ID", found.Title)
	assert.Equal(t, client.ID, found.ClientID)
}

func TestCaseRepository_FindByNumber_ShouldReturnCase(t *testing.T) {
	db := connectDB(t)
	caseRepo := NewCaseRepository(db)
	clientRepo := NewClientRepository(db)
	ctx := context.Background()

	cleanupCases(t, db)
	cleanupClients(t, db)
	client := createTestClient(t, clientRepo, ctx)

	original, err := domain.NewCase(client.ID, "NUMB-67890", "Find By Number", "Desc", "TJSP")
	require.NoError(t, err)
	err = caseRepo.Create(ctx, original)
	require.NoError(t, err)

	found, err := caseRepo.FindByNumber(ctx, "NUMB-67890")
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, original.ID, found.ID)
	assert.Equal(t, "Find By Number", found.Title)
}

func TestCaseRepository_FindByID_ShouldReturnNilWhenNotFound(t *testing.T) {
	db := connectDB(t)
	caseRepo := NewCaseRepository(db)
	ctx := context.Background()

	cleanupCases(t, db)

	found, err := caseRepo.FindByID(ctx, uuid.New())
	require.NoError(t, err)
	assert.Nil(t, found)
}

func TestCaseRepository_FindByNumber_ShouldReturnNilWhenNotFound(t *testing.T) {
	db := connectDB(t)
	caseRepo := NewCaseRepository(db)
	ctx := context.Background()

	cleanupCases(t, db)

	found, err := caseRepo.FindByNumber(ctx, "NONEXISTENT")
	require.NoError(t, err)
	assert.Nil(t, found)
}

func createTestClient(t *testing.T, repo *ClientRepository, ctx context.Context) *clientdomain.Client {
	t.Helper()
	n := atomic.AddInt64(&cpfCounter, 1)
	cpf := fmt.Sprintf("%011d", 10000000000+n)
	client, err := clientdomain.NewClient("Test Client", cpf, "test@example.com", "11999999999")
	require.NoError(t, err)
	err = repo.Create(ctx, client)
	require.NoError(t, err)
	return client
}

func TestCaseRepository_List_ShouldReturnEmpty(t *testing.T) {
	db := connectDB(t)
	caseRepo := NewCaseRepository(db)
	ctx := context.Background()

	cleanupCases(t, db)
	cleanupClients(t, db)

	cases, total, err := caseRepo.List(ctx, domain.ListCasesParams{
		Limit:  20,
		Offset: 0,
	})
	require.NoError(t, err)
	assert.Empty(t, cases)
	assert.Equal(t, int64(0), total)
}

func TestCaseRepository_List_ShouldReturnAllCases(t *testing.T) {
	db := connectDB(t)
	caseRepo := NewCaseRepository(db)
	clientRepo := NewClientRepository(db)
	ctx := context.Background()

	cleanupCases(t, db)
	cleanupClients(t, db)
	client := createTestClient(t, clientRepo, ctx)

	c1, _ := domain.NewCase(client.ID, "LIST-001", "Alpha Case", "Desc", "TJSP")
	c2, _ := domain.NewCase(client.ID, "LIST-002", "Beta Case", "Desc", "TJPE")
	require.NoError(t, caseRepo.Create(ctx, c1))
	require.NoError(t, caseRepo.Create(ctx, c2))

	cases, total, err := caseRepo.List(ctx, domain.ListCasesParams{
		Limit:  20,
		Offset: 0,
	})
	require.NoError(t, err)
	assert.Len(t, cases, 2)
	assert.Equal(t, int64(2), total)
}

func TestCaseRepository_List_ShouldFilterByStatus(t *testing.T) {
	db := connectDB(t)
	caseRepo := NewCaseRepository(db)
	clientRepo := NewClientRepository(db)
	ctx := context.Background()

	cleanupCases(t, db)
	cleanupClients(t, db)
	client := createTestClient(t, clientRepo, ctx)

	c1, _ := domain.NewCase(client.ID, "STA-001", "Status 1", "Desc", "TJSP")
	c2, _ := domain.NewCase(client.ID, "STA-002", "Status 2", "Desc", "TJSP")
	require.NoError(t, caseRepo.Create(ctx, c1))
	require.NoError(t, caseRepo.Create(ctx, c2))

	cases, total, err := caseRepo.List(ctx, domain.ListCasesParams{
		Limit:  20,
		Offset: 0,
		Status: string(domain.CaseStatusDraft),
	})
	require.NoError(t, err)
	assert.Len(t, cases, 2)
	assert.Equal(t, int64(2), total)

	cases, total, err = caseRepo.List(ctx, domain.ListCasesParams{
		Limit:  20,
		Offset: 0,
		Status: string(domain.CaseStatusActive),
	})
	require.NoError(t, err)
	assert.Empty(t, cases)
	assert.Equal(t, int64(0), total)
}

func TestCaseRepository_List_ShouldFilterByNumber(t *testing.T) {
	db := connectDB(t)
	caseRepo := NewCaseRepository(db)
	clientRepo := NewClientRepository(db)
	ctx := context.Background()

	cleanupCases(t, db)
	cleanupClients(t, db)
	client := createTestClient(t, clientRepo, ctx)

	c1, _ := domain.NewCase(client.ID, "NUM-FILTER-001", "Filtered", "Desc", "TJSP")
	c2, _ := domain.NewCase(client.ID, "OTHER-002", "Not Found", "Desc", "TJSP")
	require.NoError(t, caseRepo.Create(ctx, c1))
	require.NoError(t, caseRepo.Create(ctx, c2))

	cases, total, err := caseRepo.List(ctx, domain.ListCasesParams{
		Limit:  20,
		Offset: 0,
		Number: "FILTER",
	})
	require.NoError(t, err)
	assert.Len(t, cases, 1)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, "NUM-FILTER-001", cases[0].Number)
}

func TestCaseRepository_List_ShouldSortByTitleDesc(t *testing.T) {
	db := connectDB(t)
	caseRepo := NewCaseRepository(db)
	clientRepo := NewClientRepository(db)
	ctx := context.Background()

	cleanupCases(t, db)
	cleanupClients(t, db)
	client := createTestClient(t, clientRepo, ctx)

	c1, _ := domain.NewCase(client.ID, "SORT-001", "Alpha Case", "Desc", "TJSP")
	c2, _ := domain.NewCase(client.ID, "SORT-002", "Beta Case", "Desc", "TJPE")
	require.NoError(t, caseRepo.Create(ctx, c1))
	require.NoError(t, caseRepo.Create(ctx, c2))

	cases, total, err := caseRepo.List(ctx, domain.ListCasesParams{
		Limit:  20,
		Offset: 0,
		Sort:   "title",
		Order:  "desc",
	})
	require.NoError(t, err)
	assert.Len(t, cases, 2)
	assert.Equal(t, int64(2), total)
	assert.Equal(t, "Beta Case", cases[0].Title)
	assert.Equal(t, "Alpha Case", cases[1].Title)
}

func TestCaseRepository_UpdateStatus_ShouldPersist(t *testing.T) {
	db := connectDB(t)
	caseRepo := NewCaseRepository(db)
	clientRepo := NewClientRepository(db)
	ctx := context.Background()

	cleanupCases(t, db)
	cleanupClients(t, db)
	client := createTestClient(t, clientRepo, ctx)

	c, _ := domain.NewCase(client.ID, "STAT-001", "Status Update", "Desc", "TJSP")
	require.NoError(t, caseRepo.Create(ctx, c))
	assert.Equal(t, domain.CaseStatusDraft, c.Status)

	err := caseRepo.UpdateStatus(ctx, c.ID, domain.CaseStatusActive)
	require.NoError(t, err)

	updated, err := caseRepo.FindByID(ctx, c.ID)
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, domain.CaseStatusActive, updated.Status)
	assert.True(t, updated.UpdatedAt.After(c.UpdatedAt))
}

func TestCaseRepository_UpdateStatus_ShouldAllowMultipleUpdates(t *testing.T) {
	db := connectDB(t)
	caseRepo := NewCaseRepository(db)
	clientRepo := NewClientRepository(db)
	ctx := context.Background()

	cleanupCases(t, db)
	cleanupClients(t, db)
	client := createTestClient(t, clientRepo, ctx)

	c, _ := domain.NewCase(client.ID, "STAT-002", "Multiple Updates", "Desc", "TJSP")
	require.NoError(t, caseRepo.Create(ctx, c))

	require.NoError(t, caseRepo.UpdateStatus(ctx, c.ID, domain.CaseStatusActive))
	require.NoError(t, caseRepo.UpdateStatus(ctx, c.ID, domain.CaseStatusSuspended))
	require.NoError(t, caseRepo.UpdateStatus(ctx, c.ID, domain.CaseStatusClosed))

	updated, err := caseRepo.FindByID(ctx, c.ID)
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, domain.CaseStatusClosed, updated.Status)
}

func TestCaseRepository_List_ShouldPaginate(t *testing.T) {
	db := connectDB(t)
	caseRepo := NewCaseRepository(db)
	clientRepo := NewClientRepository(db)
	ctx := context.Background()

	cleanupCases(t, db)
	cleanupClients(t, db)
	client := createTestClient(t, clientRepo, ctx)

	for i := 0; i < 5; i++ {
		c, _ := domain.NewCase(client.ID, fmt.Sprintf("PAG-%03d", i+1), "Case Test", "Desc", "TJSP")
		require.NoError(t, caseRepo.Create(ctx, c))
	}

	cases, total, err := caseRepo.List(ctx, domain.ListCasesParams{
		Limit:  2,
		Offset: 0,
	})
	require.NoError(t, err)
	assert.Len(t, cases, 2)
	assert.Equal(t, int64(5), total)

	cases, total, err = caseRepo.List(ctx, domain.ListCasesParams{
		Limit:  2,
		Offset: 4,
	})
	require.NoError(t, err)
	assert.Len(t, cases, 1)
	assert.Equal(t, int64(5), total)
}

func cleanupCases(t *testing.T, db *sql.DB) {
	t.Helper()
	_, err := db.Exec("DELETE FROM documents")
	require.NoError(t, err)
	_, err = db.Exec("DELETE FROM timeline_events")
	require.NoError(t, err)
	_, err = db.Exec("DELETE FROM case_status_history")
	require.NoError(t, err)
	_, err = db.Exec("DELETE FROM hearings")
	require.NoError(t, err)
	_, err = db.Exec("DELETE FROM cases")
	require.NoError(t, err)
}
