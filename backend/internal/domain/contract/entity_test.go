package contract

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewContract_ShouldCreateSuccessfully(t *testing.T) {
	clientID := uuid.New()
	caseID := uuid.New()
	amount := decimal.NewFromFloat(5000.00)
	startDate := time.Now().AddDate(0, 0, 1)
	endDate := startDate.AddDate(1, 0, 0)

	contract, err := NewContract(clientID, caseID, "Honorários Advocatícios", "Descrição", ContractTypeFixedFee, amount, startDate, &endDate)
	require.NoError(t, err)
	require.NotNil(t, contract)
	assert.NotEmpty(t, contract.ID)
	assert.Equal(t, clientID, contract.ClientID)
	assert.Equal(t, caseID, contract.CaseID)
	assert.Equal(t, "Honorários Advocatícios", contract.Title)
	assert.Equal(t, ContractTypeFixedFee, contract.Type)
	assert.True(t, contract.Amount.Equal(amount))
	assert.Equal(t, startDate, contract.StartDate)
	assert.Equal(t, endDate, *contract.EndDate)
	assert.True(t, contract.Active)
	assert.False(t, contract.CreatedAt.IsZero())
	assert.False(t, contract.UpdatedAt.IsZero())
}

func TestNewContract_ShouldBeActiveByDefault(t *testing.T) {
	clientID := uuid.New()
	caseID := uuid.New()
	startDate := time.Now().AddDate(0, 0, 1)

	contract, err := NewContract(clientID, caseID, "Contrato Ativo", "", ContractTypeMonthly, decimal.NewFromFloat(1000), startDate, nil)
	require.NoError(t, err)
	assert.True(t, contract.Active)
}

func TestNewContract_ShouldAllowNilEndDate(t *testing.T) {
	clientID := uuid.New()
	caseID := uuid.New()
	startDate := time.Now().AddDate(0, 0, 1)

	contract, err := NewContract(clientID, caseID, "Sem Prazo", "", ContractTypeHourly, decimal.NewFromFloat(200), startDate, nil)
	require.NoError(t, err)
	assert.Nil(t, contract.EndDate)
}

func TestNewContract_ShouldReturnErrInvalidClientID(t *testing.T) {
	_, err := NewContract(uuid.Nil, uuid.New(), "Title", "", ContractTypeFixedFee, decimal.NewFromFloat(100), time.Now(), nil)
	assert.ErrorIs(t, err, ErrInvalidClientID)
}

func TestNewContract_ShouldReturnErrInvalidCaseID(t *testing.T) {
	_, err := NewContract(uuid.New(), uuid.Nil, "Title", "", ContractTypeFixedFee, decimal.NewFromFloat(100), time.Now(), nil)
	assert.ErrorIs(t, err, ErrInvalidCaseID)
}

func TestNewContract_ShouldReturnErrInvalidTitleWhenEmpty(t *testing.T) {
	_, err := NewContract(uuid.New(), uuid.New(), "", "", ContractTypeFixedFee, decimal.NewFromFloat(100), time.Now(), nil)
	assert.ErrorIs(t, err, ErrInvalidTitle)
}

func TestNewContract_ShouldReturnErrInvalidTitleWhenLessThan3(t *testing.T) {
	_, err := NewContract(uuid.New(), uuid.New(), "AB", "", ContractTypeFixedFee, decimal.NewFromFloat(100), time.Now(), nil)
	assert.ErrorIs(t, err, ErrInvalidTitle)
}

func TestNewContract_ShouldReturnErrInvalidContractType(t *testing.T) {
	_, err := NewContract(uuid.New(), uuid.New(), "Title", "", "invalid", decimal.NewFromFloat(100), time.Now(), nil)
	assert.ErrorIs(t, err, ErrInvalidContractType)
}

func TestNewContract_ShouldReturnErrInvalidAmountWhenZero(t *testing.T) {
	_, err := NewContract(uuid.New(), uuid.New(), "Title", "", ContractTypeFixedFee, decimal.Zero, time.Now(), nil)
	assert.ErrorIs(t, err, ErrInvalidAmount)
}

func TestNewContract_ShouldReturnErrInvalidAmountWhenNegative(t *testing.T) {
	_, err := NewContract(uuid.New(), uuid.New(), "Title", "", ContractTypeFixedFee, decimal.NewFromFloat(-100), time.Now(), nil)
	assert.ErrorIs(t, err, ErrInvalidAmount)
}

func TestNewContract_ShouldReturnErrInvalidStartDate(t *testing.T) {
	_, err := NewContract(uuid.New(), uuid.New(), "Title", "", ContractTypeFixedFee, decimal.NewFromFloat(100), time.Time{}, nil)
	assert.ErrorIs(t, err, ErrInvalidStartDate)
}

func TestNewContract_ShouldReturnErrInvalidEndDate(t *testing.T) {
	startDate := time.Now().AddDate(0, 0, 1)
	endDate := startDate.AddDate(0, 0, -2)
	_, err := NewContract(uuid.New(), uuid.New(), "Title", "", ContractTypeFixedFee, decimal.NewFromFloat(100), startDate, &endDate)
	assert.ErrorIs(t, err, ErrInvalidEndDate)
}

func TestNewContract_ShouldAcceptAllTypes(t *testing.T) {
	clientID := uuid.New()
	caseID := uuid.New()
	startDate := time.Now().AddDate(0, 0, 1)
	types := []ContractType{ContractTypeFixedFee, ContractTypeHourly, ContractTypeSuccessFee, ContractTypeMonthly}
	for _, ct := range types {
		contract, err := NewContract(clientID, caseID, "Valid Title", "", ct, decimal.NewFromFloat(100), startDate, nil)
		require.NoError(t, err)
		assert.Equal(t, ct, contract.Type)
	}
}
