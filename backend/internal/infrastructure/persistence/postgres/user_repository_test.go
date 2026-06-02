package postgres

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domain "legalflow/internal/user/domain/user"
)

func connectDB(t *testing.T) *sql.DB {
	if os.Getenv("DB_HOST") == "" {
		t.Skip("DB_HOST not set, skipping integration test")
	}

	dsn := "host=" + os.Getenv("DB_HOST") +
		" port=" + os.Getenv("DB_PORT") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)

	t.Cleanup(func() { db.Close() })

	return db
}

func TestUserRepository_Create_ShouldSaveSuccessfully(t *testing.T) {
	db := connectDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	user, err := domain.NewUser("Jane Doe", "jane@example.com", "hash123")
	require.NoError(t, err)

	err = repo.Create(ctx, user)
	require.NoError(t, err)
	assert.NotEmpty(t, user.ID)
}

func TestUserRepository_Create_ShouldReturnDuplicateEmail(t *testing.T) {
	db := connectDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	user1, err := domain.NewUser("John Doe", "duplicate@example.com", "hash123")
	require.NoError(t, err)
	err = repo.Create(ctx, user1)
	require.NoError(t, err)

	user2, err := domain.NewUser("Jane Doe", "duplicate@example.com", "hash456")
	require.NoError(t, err)
	err = repo.Create(ctx, user2)
	require.Error(t, err)
}

func TestUserRepository_FindByEmail_ShouldReturnUser(t *testing.T) {
	db := connectDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	original, err := domain.NewUser("Find Email", "findemail@example.com", "hash123")
	require.NoError(t, err)
	err = repo.Create(ctx, original)
	require.NoError(t, err)

	found, err := repo.FindByEmail(ctx, "findemail@example.com")
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, original.ID, found.ID)
	assert.Equal(t, "Find Email", found.Name)
}

func TestUserRepository_FindByID_ShouldReturnUser(t *testing.T) {
	db := connectDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	original, err := domain.NewUser("Find ID", "findid@example.com", "hash123")
	require.NoError(t, err)
	err = repo.Create(ctx, original)
	require.NoError(t, err)

	found, err := repo.FindByID(ctx, original.ID)
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, original.ID, found.ID)
	assert.Equal(t, "Find ID", found.Name)
}

func TestUserRepository_FindByEmail_ShouldReturnNilWhenNotFound(t *testing.T) {
	db := connectDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	found, err := repo.FindByEmail(ctx, "nonexistent@example.com")
	require.NoError(t, err)
	assert.Nil(t, found)
}

func TestUserRepository_FindByID_ShouldReturnNilWhenNotFound(t *testing.T) {
	db := connectDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	found, err := repo.FindByID(ctx, uuid.MustParse("00000000-0000-0000-0000-000000000000"))
	require.NoError(t, err)
	assert.Nil(t, found)
}
