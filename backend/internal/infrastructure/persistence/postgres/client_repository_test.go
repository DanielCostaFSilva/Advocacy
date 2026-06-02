package postgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domain "legalflow/internal/domain/client"
)

func TestClientRepository_Create_ShouldSaveSuccessfully(t *testing.T) {
	db := connectDB(t)
	repo := NewClientRepository(db)
	ctx := context.Background()

	cleanupClients(t, db)
	client, err := domain.NewClient("Jane Doe", "11122233344", "jane@example.com", "11999999999")
	require.NoError(t, err)

	err = repo.Create(ctx, client)
	require.NoError(t, err)
	assert.NotEmpty(t, client.ID)
}

func TestClientRepository_Create_ShouldReturnDuplicateCPFError(t *testing.T) {
	db := connectDB(t)
	repo := NewClientRepository(db)
	ctx := context.Background()

	cleanupClients(t, db)
	client1, err := domain.NewClient("John Doe", "99988877766", "john@example.com", "11999999999")
	require.NoError(t, err)
	err = repo.Create(ctx, client1)
	require.NoError(t, err)

	client2, err := domain.NewClient("Jane Doe", "99988877766", "jane@example.com", "11988888888")
	require.NoError(t, err)
	err = repo.Create(ctx, client2)
	require.Error(t, err)
}

func TestClientRepository_FindByCPF_ShouldReturnClient(t *testing.T) {
	db := connectDB(t)
	repo := NewClientRepository(db)
	ctx := context.Background()

	cleanupClients(t, db)
	original, err := domain.NewClient("Find CPF", "55544433322", "findcpf@example.com", "11999999999")
	require.NoError(t, err)
	err = repo.Create(ctx, original)
	require.NoError(t, err)

	found, err := repo.FindByCPF(ctx, "55544433322")
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, original.ID, found.ID)
	assert.Equal(t, "Find CPF", found.Name)
	assert.Equal(t, "55544433322", found.CPF)
}

func TestClientRepository_FindByID_ShouldReturnClient(t *testing.T) {
	db := connectDB(t)
	repo := NewClientRepository(db)
	ctx := context.Background()

	cleanupClients(t, db)
	original, err := domain.NewClient("Find ID", "66677788899", "findid@example.com", "11999999999")
	require.NoError(t, err)
	err = repo.Create(ctx, original)
	require.NoError(t, err)

	found, err := repo.FindByID(ctx, original.ID)
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, original.ID, found.ID)
	assert.Equal(t, "Find ID", found.Name)
	assert.Equal(t, "66677788899", found.CPF)
}

func TestClientRepository_Update_ShouldSaveChanges(t *testing.T) {
	db := connectDB(t)
	repo := NewClientRepository(db)
	ctx := context.Background()

	cleanupClients(t, db)
	original, err := domain.NewClient("Maria Silva", "12345678900", "maria@example.com", "81999999999")
	require.NoError(t, err)
	err = repo.Create(ctx, original)
	require.NoError(t, err)

	require.NoError(t, original.Update("Maria da Silva", "maria.silva@example.com", "81988887777"))
	err = repo.Update(ctx, original)
	require.NoError(t, err)

	found, err := repo.FindByID(ctx, original.ID)
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, "Maria da Silva", found.Name)
	assert.Equal(t, "maria.silva@example.com", found.Email)
	assert.Equal(t, "81988887777", found.Phone)
	assert.Equal(t, "12345678900", found.CPF)
	assert.True(t, found.UpdatedAt.After(found.CreatedAt))
}

func TestClientRepository_FindByID_ShouldReturnNilWhenNotFound(t *testing.T) {
	db := connectDB(t)
	repo := NewClientRepository(db)
	ctx := context.Background()

	cleanupClients(t, db)

	found, err := repo.FindByID(ctx, uuid.New())
	require.NoError(t, err)
	assert.Nil(t, found)
}

func TestClientRepository_List_ShouldReturnPaginatedResult(t *testing.T) {
	db := connectDB(t)
	repo := NewClientRepository(db)
	ctx := context.Background()

	cleanupClients(t, db)
	assert.NoError(t, repo.Create(ctx, mustClient(t, "Alice", "11111111111")))
	assert.NoError(t, repo.Create(ctx, mustClient(t, "Bob", "22222222222")))

	clients, total, err := repo.List(ctx, domain.ListClientsParams{
		Offset: 0,
		Limit:  10,
	})
	require.NoError(t, err)
	assert.Equal(t, 2, len(clients))
	assert.Equal(t, int64(2), total)
}

func TestClientRepository_List_ShouldFilterByName(t *testing.T) {
	db := connectDB(t)
	repo := NewClientRepository(db)
	ctx := context.Background()

	cleanupClients(t, db)
	assert.NoError(t, repo.Create(ctx, mustClient(t, "Alice", "33333333333")))
	assert.NoError(t, repo.Create(ctx, mustClient(t, "Bob", "44444444444")))

	clients, total, err := repo.List(ctx, domain.ListClientsParams{
		Offset: 0,
		Limit:  10,
		Name:   "Ali",
	})
	require.NoError(t, err)
	assert.Len(t, clients, 1)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, "Alice", clients[0].Name)
}

func TestClientRepository_List_ShouldFilterByCPF(t *testing.T) {
	db := connectDB(t)
	repo := NewClientRepository(db)
	ctx := context.Background()

	cleanupClients(t, db)
	assert.NoError(t, repo.Create(ctx, mustClient(t, "Alice", "55555555555")))
	assert.NoError(t, repo.Create(ctx, mustClient(t, "Bob", "66666666666")))

	clients, total, err := repo.List(ctx, domain.ListClientsParams{
		Offset: 0,
		Limit:  10,
		CPF:    "55555555555",
	})
	require.NoError(t, err)
	assert.Len(t, clients, 1)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, "Alice", clients[0].Name)
}

func TestClientRepository_List_ShouldSortByName(t *testing.T) {
	db := connectDB(t)
	repo := NewClientRepository(db)
	ctx := context.Background()

	cleanupClients(t, db)
	assert.NoError(t, repo.Create(ctx, mustClient(t, "Zara", "77777777777")))
	assert.NoError(t, repo.Create(ctx, mustClient(t, "Anna", "88888888888")))

	clients, _, err := repo.List(ctx, domain.ListClientsParams{
		Offset: 0,
		Limit:  10,
		Sort:   "name",
		Order:  "asc",
	})
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(clients), 2)
	assert.Equal(t, "Anna", clients[0].Name)
}

func mustClient(t *testing.T, name, cpf string) *domain.Client {
	t.Helper()
	client, err := domain.NewClient(name, cpf, name+"@example.com", "11999999999")
	require.NoError(t, err)
	return client
}

func cleanupClients(t *testing.T, db *sql.DB) {
	t.Helper()
	_, err := db.Exec("DELETE FROM clients")
	require.NoError(t, err)
}
