package legalcase

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCase_ShouldCreateSuccessfully(t *testing.T) {
	clientID := uuid.New()
	c, err := NewCase(clientID, "ABC-12345", "Indenizatória", "Descrição do caso", "TJSP")

	require.NoError(t, err)
	require.NotNil(t, c)
	assert.NotEmpty(t, c.ID)
	assert.Equal(t, clientID, c.ClientID)
	assert.Equal(t, "ABC-12345", c.Number)
	assert.Equal(t, "Indenizatória", c.Title)
	assert.Equal(t, "Descrição do caso", c.Description)
	assert.Equal(t, "TJSP", c.Court)
	assert.Equal(t, CaseStatusDraft, c.Status)
	assert.False(t, c.CreatedAt.IsZero())
	assert.False(t, c.UpdatedAt.IsZero())
}

func TestNewCase_ShouldReturnErrInvalidCaseNumberWhenEmpty(t *testing.T) {
	_, err := NewCase(uuid.New(), "", "Title", "Desc", "TJSP")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidCaseNumber)
}

func TestNewCase_ShouldReturnErrInvalidCaseNumberWhenLessThan5(t *testing.T) {
	_, err := NewCase(uuid.New(), "AB", "Title", "Desc", "TJSP")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidCaseNumber)
}

func TestNewCase_ShouldReturnErrInvalidCaseTitleWhenEmpty(t *testing.T) {
	_, err := NewCase(uuid.New(), "ABC-12345", "", "Desc", "TJSP")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidCaseTitle)
}

func TestNewCase_ShouldReturnErrInvalidCaseTitleWhenLessThan5(t *testing.T) {
	_, err := NewCase(uuid.New(), "ABC-12345", "Tit", "Desc", "TJSP")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidCaseTitle)
}

func TestNewCase_ShouldReturnErrInvalidCourtWhenEmpty(t *testing.T) {
	_, err := NewCase(uuid.New(), "ABC-12345", "Title", "Desc", "")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidCourt)
}

func TestNewCase_ShouldAllowEmptyDescription(t *testing.T) {
	c, err := NewCase(uuid.New(), "ABC-12345", "Title Title", "", "TJSP")
	require.NoError(t, err)
	assert.Empty(t, c.Description)
}

func TestNewCase_ShouldReturnErrInvalidClientIDWhenNil(t *testing.T) {
	_, err := NewCase(uuid.Nil, "ABC-12345", "Title Title", "Desc", "TJSP")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidClientID)
}

func TestChangeStatus_DraftToActive_ShouldSucceed(t *testing.T) {
	c, _ := NewCase(uuid.New(), "ABC-12345", "Title Case", "Desc", "TJSP")
	err := c.ChangeStatus(CaseStatusActive)
	require.NoError(t, err)
	assert.Equal(t, CaseStatusActive, c.Status)
}

func TestChangeStatus_DraftToSuspended_ShouldReturnError(t *testing.T) {
	c, _ := NewCase(uuid.New(), "ABC-12345", "Title Case", "Desc", "TJSP")
	err := c.ChangeStatus(CaseStatusSuspended)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidStatusTransition)
	assert.Equal(t, CaseStatusDraft, c.Status)
}

func TestChangeStatus_ClosedToAny_ShouldReturnError(t *testing.T) {
	c, _ := NewCase(uuid.New(), "ABC-12345", "Title Case", "Desc", "TJSP")
	c.Status = CaseStatusClosed
	err := c.ChangeStatus(CaseStatusActive)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidStatusTransition)
	assert.Equal(t, CaseStatusClosed, c.Status)
}

func TestChangeStatus_SuspendedToClosed_ShouldSucceed(t *testing.T) {
	c, _ := NewCase(uuid.New(), "ABC-12345", "Title Case", "Desc", "TJSP")
	c.Status = CaseStatusSuspended
	err := c.ChangeStatus(CaseStatusClosed)
	require.NoError(t, err)
	assert.Equal(t, CaseStatusClosed, c.Status)
}

func TestChangeStatus_InvalidStatus_ShouldReturnError(t *testing.T) {
	c, _ := NewCase(uuid.New(), "ABC-12345", "Title Case", "Desc", "TJSP")
	err := c.ChangeStatus("invalid")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidStatus)
}

func TestChangeStatus_SameStatus_ShouldSucceed(t *testing.T) {
	c, _ := NewCase(uuid.New(), "ABC-12345", "Title Case", "Desc", "TJSP")
	err := c.ChangeStatus(CaseStatusDraft)
	require.NoError(t, err)
	assert.Equal(t, CaseStatusDraft, c.Status)
}
