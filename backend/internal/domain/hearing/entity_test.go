package hearing

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHearing_ShouldCreateSuccessfully(t *testing.T) {
	caseID := uuid.New()
	future := time.Now().Add(24 * time.Hour)

	h, err := NewHearing(caseID, "Audiência de Conciliação", "Descrição", HearingTypeConciliation, "Fórum Central", future)

	require.NoError(t, err)
	require.NotNil(t, h)
	assert.NotEmpty(t, h.ID)
	assert.Equal(t, caseID, h.CaseID)
	assert.Equal(t, "Audiência de Conciliação", h.Title)
	assert.Equal(t, "Descrição", h.Description)
	assert.Equal(t, HearingTypeConciliation, h.Type)
	assert.Equal(t, "Fórum Central", h.Location)
	assert.Equal(t, future, h.ScheduledAt)
	assert.False(t, h.CreatedAt.IsZero())
	assert.False(t, h.UpdatedAt.IsZero())
}

func TestNewHearing_ShouldReturnErrInvalidCaseID(t *testing.T) {
	_, err := NewHearing(uuid.Nil, "Audiência", "Desc", HearingTypeConciliation, "Local", time.Now().Add(24*time.Hour))
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidCaseID)
}

func TestNewHearing_ShouldReturnErrInvalidTitleWhenEmpty(t *testing.T) {
	_, err := NewHearing(uuid.New(), "", "Desc", HearingTypeConciliation, "Local", time.Now().Add(24*time.Hour))
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidTitle)
}

func TestNewHearing_ShouldReturnErrInvalidTitleWhenLessThan5(t *testing.T) {
	_, err := NewHearing(uuid.New(), "ABC", "Desc", HearingTypeConciliation, "Local", time.Now().Add(24*time.Hour))
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidTitle)
}

func TestNewHearing_ShouldReturnErrInvalidHearingType(t *testing.T) {
	_, err := NewHearing(uuid.New(), "Audiência", "Desc", "invalid", "Local", time.Now().Add(24*time.Hour))
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidHearingType)
}

func TestNewHearing_ShouldReturnErrInvalidLocation(t *testing.T) {
	_, err := NewHearing(uuid.New(), "Audiência", "Desc", HearingTypeConciliation, "", time.Now().Add(24*time.Hour))
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidLocation)
}

func TestNewHearing_ShouldReturnErrInvalidScheduledAtWhenPast(t *testing.T) {
	_, err := NewHearing(uuid.New(), "Audiência", "Desc", HearingTypeConciliation, "Local", time.Now().Add(-24*time.Hour))
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidScheduledAt)
}

func TestNewHearing_ShouldAllowEmptyDescription(t *testing.T) {
	caseID := uuid.New()
	h, err := NewHearing(caseID, "Audiência Completa", "", HearingTypeJudgment, "Fórum", time.Now().Add(24*time.Hour))

	require.NoError(t, err)
	assert.Empty(t, h.Description)
}

func TestNewHearing_ShouldAcceptAllTypes(t *testing.T) {
	future := time.Now().Add(24 * time.Hour)
	caseID := uuid.New()

	types := []HearingType{HearingTypeConciliation, HearingTypeInstruction, HearingTypeJudgment, HearingTypeVirtual}
	for _, ht := range types {
		h, err := NewHearing(caseID, "Audiência Teste", "Desc", ht, "Local", future)
		require.NoError(t, err)
		assert.Equal(t, ht, h.Type)
	}
}
