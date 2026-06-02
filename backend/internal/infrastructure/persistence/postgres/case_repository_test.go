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

func cleanupCases(t *testing.T, db *sql.DB) {
	t.Helper()
	_, err := db.Exec("DELETE FROM cases")
	require.NoError(t, err)
}
