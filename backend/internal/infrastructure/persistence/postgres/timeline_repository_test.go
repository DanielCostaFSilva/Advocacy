package postgres

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domain "legalflow/internal/domain/timeline"
	legalcaseDomain "legalflow/internal/domain/legalcase"
)

func TestTimelineRepository_FindByCaseID_ShouldReturnOrderedEvents(t *testing.T) {
	db := connectDB(t)
	caseRepo := NewCaseRepository(db)
	clientRepo := NewClientRepository(db)
	timelineRepo := NewTimelineRepository(db)
	ctx := context.Background()

	cleanupCases(t, db)
	cleanupClients(t, db)
	client := createTestClient(t, clientRepo, ctx)

	c, _ := legalcaseDomain.NewCase(client.ID, "TL-001", "Timeline Test", "Desc", "TJSP")
	require.NoError(t, caseRepo.Create(ctx, c))

	e1 := domain.NewEvent(c.ID, domain.EventCaseCreated, "Processo criado")
	require.NoError(t, timelineRepo.Create(ctx, e1))

	e2 := domain.NewEvent(c.ID, domain.EventStatusChanged, "Status alterado para active")
	require.NoError(t, timelineRepo.Create(ctx, e2))

	events, err := timelineRepo.FindByCaseID(ctx, c.ID)
	require.NoError(t, err)
	assert.Len(t, events, 2)
	assert.Equal(t, domain.EventCaseCreated, events[0].Type)
	assert.Equal(t, domain.EventStatusChanged, events[1].Type)
	assert.True(t, events[1].CreatedAt.After(events[0].CreatedAt) || events[1].CreatedAt.Equal(events[0].CreatedAt))
}

func TestTimelineRepository_FindByCaseID_ShouldReturnEmptyForCaseWithoutEvents(t *testing.T) {
	db := connectDB(t)
	caseRepo := NewCaseRepository(db)
	clientRepo := NewClientRepository(db)
	timelineRepo := NewTimelineRepository(db)
	ctx := context.Background()

	cleanupCases(t, db)
	cleanupClients(t, db)
	client := createTestClient(t, clientRepo, ctx)

	c, _ := legalcaseDomain.NewCase(client.ID, "TL-002", "No Events", "Desc", "TJSP")
	require.NoError(t, caseRepo.Create(ctx, c))

	events, err := timelineRepo.FindByCaseID(ctx, c.ID)
	require.NoError(t, err)
	assert.Empty(t, events)
}

func TestTimelineRepository_FindByCaseID_ShouldReturnEmptyForNonExistentCase(t *testing.T) {
	db := connectDB(t)
	timelineRepo := NewTimelineRepository(db)
	ctx := context.Background()

	cleanupCases(t, db)

	events, err := timelineRepo.FindByCaseID(ctx, uuid.New())
	require.NoError(t, err)
	assert.Empty(t, events)
}
