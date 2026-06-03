package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	contractDomain "legalflow/internal/domain/contract"
)

func TestContractRepository_Create_ShouldSaveSuccessfully(t *testing.T) {
	db := connectDB(t)
	contractRepo := NewContractRepository(db)
	clientRepo := NewClientRepository(db)
	caseRepo := NewCaseRepository(db)
	ctx := context.Background()

	cleanupClients(t, db)
	cleanupCases(t, db)

	client := createTestClient(t, clientRepo, ctx)
	c := createTestCase(t, caseRepo, client.ID, ctx)

	startDate := time.Now().AddDate(0, 0, 1)
	endDate := startDate.AddDate(1, 0, 0)
	contract, err := contractDomain.NewContract(client.ID, c.ID, "Honorários Advocatícios", "Descrição completa", contractDomain.ContractTypeFixedFee, decimal.NewFromFloat(5000.00), startDate, &endDate)
	require.NoError(t, err)
	require.NotNil(t, contract)

	err = contractRepo.Create(ctx, contract)
	require.NoError(t, err)
	assert.NotEmpty(t, contract.ID)
	assert.False(t, contract.CreatedAt.IsZero())
	assert.False(t, contract.UpdatedAt.IsZero())
	assert.True(t, contract.Active)
}

func TestContractRepository_FindByID_ShouldReturnContract(t *testing.T) {
	db := connectDB(t)
	contractRepo := NewContractRepository(db)
	clientRepo := NewClientRepository(db)
	caseRepo := NewCaseRepository(db)
	ctx := context.Background()

	cleanupClients(t, db)
	cleanupCases(t, db)

	client := createTestClient(t, clientRepo, ctx)
	c := createTestCase(t, caseRepo, client.ID, ctx)

	startDate := time.Now().AddDate(0, 0, 1)
	contract, err := contractDomain.NewContract(client.ID, c.ID, "Contrato Mensal", "", contractDomain.ContractTypeMonthly, decimal.NewFromFloat(2500.00), startDate, nil)
	require.NoError(t, err)
	require.NoError(t, contractRepo.Create(ctx, contract))

	found, err := contractRepo.FindByID(ctx, contract.ID)
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, contract.ID, found.ID)
	assert.Equal(t, "Contrato Mensal", found.Title)
	assert.Equal(t, contractDomain.ContractTypeMonthly, found.Type)
	assert.True(t, found.Amount.Equal(decimal.NewFromFloat(2500.00)))
	assert.True(t, found.Active)
}

func TestContractRepository_FindByID_ShouldReturnNilWhenNotFound(t *testing.T) {
	db := connectDB(t)
	contractRepo := NewContractRepository(db)
	ctx := context.Background()

	found, err := contractRepo.FindByID(ctx, uuid.New())
	require.NoError(t, err)
	assert.Nil(t, found)
}

func TestContractRepository_ListByCaseID_ShouldReturnContracts(t *testing.T) {
	db := connectDB(t)
	contractRepo := NewContractRepository(db)
	clientRepo := NewClientRepository(db)
	caseRepo := NewCaseRepository(db)
	ctx := context.Background()

	cleanupClients(t, db)
	cleanupCases(t, db)

	client := createTestClient(t, clientRepo, ctx)
	c := createTestCase(t, caseRepo, client.ID, ctx)

	startDate := time.Now().AddDate(0, 0, 1)
	c1, _ := contractDomain.NewContract(client.ID, c.ID, "Contrato Principal", "", contractDomain.ContractTypeFixedFee, decimal.NewFromFloat(10000), startDate, nil)
	c2, _ := contractDomain.NewContract(client.ID, c.ID, "Contrato Adicional", "", contractDomain.ContractTypeHourly, decimal.NewFromFloat(5000), startDate, nil)
	require.NoError(t, contractRepo.Create(ctx, c1))
	require.NoError(t, contractRepo.Create(ctx, c2))

	contracts, err := contractRepo.ListByCaseID(ctx, c.ID)
	require.NoError(t, err)
	assert.Len(t, contracts, 2)
}

func TestContractRepository_ListByCaseID_ShouldReturnOnlyMatchingCase(t *testing.T) {
	db := connectDB(t)
	contractRepo := NewContractRepository(db)
	clientRepo := NewClientRepository(db)
	caseRepo := NewCaseRepository(db)
	ctx := context.Background()

	cleanupClients(t, db)
	cleanupCases(t, db)

	client := createTestClient(t, clientRepo, ctx)
	c1 := createTestCase(t, caseRepo, client.ID, ctx)
	c2 := createTestCase(t, caseRepo, client.ID, ctx)

	startDate := time.Now().AddDate(0, 0, 1)
	contract1, _ := contractDomain.NewContract(client.ID, c1.ID, "Case 1 Contract", "", contractDomain.ContractTypeFixedFee, decimal.NewFromFloat(3000), startDate, nil)
	contract2, _ := contractDomain.NewContract(client.ID, c2.ID, "Case 2 Contract", "", contractDomain.ContractTypeFixedFee, decimal.NewFromFloat(4000), startDate, nil)
	require.NoError(t, contractRepo.Create(ctx, contract1))
	require.NoError(t, contractRepo.Create(ctx, contract2))

	contracts, err := contractRepo.ListByCaseID(ctx, c1.ID)
	require.NoError(t, err)
	assert.Len(t, contracts, 1)
	assert.Equal(t, "Case 1 Contract", contracts[0].Title)
}


