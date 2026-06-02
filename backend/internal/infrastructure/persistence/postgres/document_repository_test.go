package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domain "legalflow/internal/domain/document"
)

func TestDocumentRepository_Create_ShouldSaveSuccessfully(t *testing.T) {
	db := connectDB(t)
	docRepo := NewDocumentRepository(db)
	caseRepo := NewCaseRepository(db)
	clientRepo := NewClientRepository(db)
	ctx := context.Background()

	cleanupCases(t, db)
	cleanupClients(t, db)
	client := createTestClient(t, clientRepo, ctx)
	c := createTestCase(t, caseRepo, client.ID, ctx)

	doc, err := domain.NewDocument(c.ID, "Contrato Social", "Descrição", domain.DocumentTypeContract, "contrato.pdf", "application/pdf", 1024, "cases/"+c.ID.String()+"/contrato.pdf")
	require.NoError(t, err)

	err = docRepo.Create(ctx, doc)
	require.NoError(t, err)
	assert.NotEmpty(t, doc.ID)
	assert.False(t, doc.CreatedAt.IsZero())
	assert.False(t, doc.UpdatedAt.IsZero())
}

func TestDocumentRepository_FindByID_ShouldReturnDocument(t *testing.T) {
	db := connectDB(t)
	docRepo := NewDocumentRepository(db)
	caseRepo := NewCaseRepository(db)
	clientRepo := NewClientRepository(db)
	ctx := context.Background()

	cleanupCases(t, db)
	cleanupClients(t, db)
	client := createTestClient(t, clientRepo, ctx)
	c := createTestCase(t, caseRepo, client.ID, ctx)

	original, err := domain.NewDocument(c.ID, "Petição Inicial", "Desc", domain.DocumentTypePetition, "peticao.pdf", "application/pdf", 2048, "key-peticao")
	require.NoError(t, err)
	require.NoError(t, docRepo.Create(ctx, original))

	found, err := docRepo.FindByID(ctx, original.ID)
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, original.ID, found.ID)
	assert.Equal(t, "Petição Inicial", found.Name)
	assert.Equal(t, domain.DocumentTypePetition, found.Type)
	assert.Equal(t, int64(2048), found.FileSize)
}

func TestDocumentRepository_FindByID_ShouldReturnNilWhenNotFound(t *testing.T) {
	db := connectDB(t)
	docRepo := NewDocumentRepository(db)
	ctx := context.Background()

	cleanupCases(t, db)

	found, err := docRepo.FindByID(ctx, uuid.New())
	require.NoError(t, err)
	assert.Nil(t, found)
}

func TestDocumentRepository_ListByCaseID_ShouldReturnDocuments(t *testing.T) {
	db := connectDB(t)
	docRepo := NewDocumentRepository(db)
	caseRepo := NewCaseRepository(db)
	clientRepo := NewClientRepository(db)
	ctx := context.Background()

	cleanupCases(t, db)
	cleanupClients(t, db)
	client := createTestClient(t, clientRepo, ctx)
	c := createTestCase(t, caseRepo, client.ID, ctx)

	d1, _ := domain.NewDocument(c.ID, "Documento 1", "Desc", domain.DocumentTypeEvidence, "ev1.pdf", "application/pdf", 100, "key1")
	d2, _ := domain.NewDocument(c.ID, "Documento 2", "Desc", domain.DocumentTypeEvidence, "ev2.pdf", "application/pdf", 200, "key2")
	require.NoError(t, docRepo.Create(ctx, d1))
	require.NoError(t, docRepo.Create(ctx, d2))

	docs, err := docRepo.ListByCaseID(ctx, c.ID)
	require.NoError(t, err)
	assert.Len(t, docs, 2)
}

func TestDocumentRepository_ListByCaseID_ShouldReturnOnlyMatchingCase(t *testing.T) {
	db := connectDB(t)
	docRepo := NewDocumentRepository(db)
	caseRepo := NewCaseRepository(db)
	clientRepo := NewClientRepository(db)
	ctx := context.Background()

	cleanupCases(t, db)
	cleanupClients(t, db)
	client := createTestClient(t, clientRepo, ctx)
	c1 := createTestCase(t, caseRepo, client.ID, ctx)
	c2 := createTestCase(t, caseRepo, client.ID, ctx)

	d1, _ := domain.NewDocument(c1.ID, "Doc Case 1", "Desc", domain.DocumentTypeContract, "c1.pdf", "application/pdf", 100, "key1")
	d2, _ := domain.NewDocument(c2.ID, "Doc Case 2", "Desc", domain.DocumentTypeContract, "c2.pdf", "application/pdf", 200, "key2")
	require.NoError(t, docRepo.Create(ctx, d1))
	require.NoError(t, docRepo.Create(ctx, d2))

	docs, err := docRepo.ListByCaseID(ctx, c1.ID)
	require.NoError(t, err)
	assert.Len(t, docs, 1)
	assert.Equal(t, "Doc Case 1", docs[0].Name)
}

func TestDocumentRepository_ListByCaseIDPaginated_ShouldReturnPaginatedResult(t *testing.T) {
	db := connectDB(t)
	docRepo := NewDocumentRepository(db)
	caseRepo := NewCaseRepository(db)
	clientRepo := NewClientRepository(db)
	ctx := context.Background()

	cleanupCases(t, db)
	cleanupClients(t, db)
	client := createTestClient(t, clientRepo, ctx)
	c := createTestCase(t, caseRepo, client.ID, ctx)

	for i := 0; i < 5; i++ {
		d, _ := domain.NewDocument(c.ID, fmt.Sprintf("Documento %d", i+1), "Desc", domain.DocumentTypeEvidence, fmt.Sprintf("doc%d.pdf", i+1), "application/pdf", int64(100*(i+1)), fmt.Sprintf("key%d", i+1))
		require.NoError(t, docRepo.Create(ctx, d))
	}

	docs, total, err := docRepo.ListByCaseIDPaginated(ctx, domain.ListDocumentsParams{
		CaseID: c.ID.String(),
		Limit:  2,
		Offset: 0,
	})
	require.NoError(t, err)
	assert.Len(t, docs, 2)
	assert.Equal(t, int64(5), total)

	docs, total, err = docRepo.ListByCaseIDPaginated(ctx, domain.ListDocumentsParams{
		CaseID: c.ID.String(),
		Limit:  2,
		Offset: 4,
	})
	require.NoError(t, err)
	assert.Len(t, docs, 1)
	assert.Equal(t, int64(5), total)
}

func TestDocumentRepository_ListByCaseIDPaginated_ShouldFilterByType(t *testing.T) {
	db := connectDB(t)
	docRepo := NewDocumentRepository(db)
	caseRepo := NewCaseRepository(db)
	clientRepo := NewClientRepository(db)
	ctx := context.Background()

	cleanupCases(t, db)
	cleanupClients(t, db)
	client := createTestClient(t, clientRepo, ctx)
	c := createTestCase(t, caseRepo, client.ID, ctx)

	d1, _ := domain.NewDocument(c.ID, "Contract Doc", "Desc", domain.DocumentTypeContract, "c.pdf", "application/pdf", 100, "k1")
	d2, _ := domain.NewDocument(c.ID, "Evidence Doc", "Desc", domain.DocumentTypeEvidence, "e.pdf", "application/pdf", 200, "k2")
	require.NoError(t, docRepo.Create(ctx, d1))
	require.NoError(t, docRepo.Create(ctx, d2))

	docs, total, err := docRepo.ListByCaseIDPaginated(ctx, domain.ListDocumentsParams{
		CaseID: c.ID.String(),
		Limit:  20,
		Offset: 0,
		Type:   "contract",
	})
	require.NoError(t, err)
	assert.Len(t, docs, 1)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, "Contract Doc", docs[0].Name)
}

func TestDocumentRepository_ListByCaseIDPaginated_ShouldSortByFileSize(t *testing.T) {
	db := connectDB(t)
	docRepo := NewDocumentRepository(db)
	caseRepo := NewCaseRepository(db)
	clientRepo := NewClientRepository(db)
	ctx := context.Background()

	cleanupCases(t, db)
	cleanupClients(t, db)
	client := createTestClient(t, clientRepo, ctx)
	c := createTestCase(t, caseRepo, client.ID, ctx)

	d1, _ := domain.NewDocument(c.ID, "Small", "Desc", domain.DocumentTypeContract, "s.pdf", "application/pdf", 100, "k1")
	d2, _ := domain.NewDocument(c.ID, "Large", "Desc", domain.DocumentTypeContract, "l.pdf", "application/pdf", 500, "k2")
	d3, _ := domain.NewDocument(c.ID, "Medium", "Desc", domain.DocumentTypeContract, "m.pdf", "application/pdf", 300, "k3")
	require.NoError(t, docRepo.Create(ctx, d1))
	require.NoError(t, docRepo.Create(ctx, d2))
	require.NoError(t, docRepo.Create(ctx, d3))

	docs, total, err := docRepo.ListByCaseIDPaginated(ctx, domain.ListDocumentsParams{
		CaseID: c.ID.String(),
		Limit:  20,
		Offset: 0,
		Sort:   "file_size",
		Order:  "desc",
	})
	require.NoError(t, err)
	assert.Len(t, docs, 3)
	assert.Equal(t, int64(3), total)
	assert.Equal(t, "Large", docs[0].Name)
	assert.Equal(t, "Medium", docs[1].Name)
	assert.Equal(t, "Small", docs[2].Name)
}
