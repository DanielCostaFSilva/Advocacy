package postgres

import (
	"context"
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
