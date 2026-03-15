//go:build integration

package tracker_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	infradb "github.com/wlindb/issue-tracker/internal/infrastructure/db"
	tracker "github.com/wlindb/issue-tracker/internal/infrastructure/tracker"
)

var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()
	pool, terminate, err := startPostgres(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "start postgres: %v\n", err)
		os.Exit(1)
	}
	defer terminate()

	if err := tracker.Migrate(ctx, pool); err != nil {
		fmt.Fprintf(os.Stderr, "migrate: %v\n", err)
		os.Exit(1)
	}
	testPool = pool
	os.Exit(m.Run())
}

func startPostgres(ctx context.Context) (*pgxpool.Pool, func(), error) {
	req := testcontainers.ContainerRequest{
		Image: "postgres:17-alpine",
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "test",
		},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
	}
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req, Started: true,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("start container: %w", err)
	}
	port, err := c.MappedPort(ctx, "5432")
	if err != nil {
		c.Terminate(ctx) //nolint:errcheck
		return nil, nil, fmt.Errorf("mapped port: %w", err)
	}
	host, err := c.Host(ctx)
	if err != nil {
		c.Terminate(ctx) //nolint:errcheck
		return nil, nil, fmt.Errorf("host: %w", err)
	}
	dsn := fmt.Sprintf("postgres://test:test@%s:%s/test", host, port.Port())
	pool, err := infradb.New(ctx, dsn)
	if err != nil {
		c.Terminate(ctx) //nolint:errcheck
		return nil, nil, fmt.Errorf("connect: %w", err)
	}
	return pool, func() { pool.Close(); c.Terminate(ctx) }, nil //nolint:errcheck
}

func TestProjectRepository_Create_Success(t *testing.T) {
	repo := tracker.NewProjectRepository(testPool)
	id, ownerID := uuid.New(), uuid.New()

	got, err := repo.Create(context.Background(), id, ownerID, "Acme", nil)

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, id, got.ID)
	assert.Equal(t, ownerID, got.OwnerID)
	assert.Equal(t, "Acme", got.Name)
	assert.Nil(t, got.Description)
	assert.False(t, got.CreatedAt.IsZero())
	assert.False(t, got.UpdatedAt.IsZero())
}

func TestProjectRepository_Create_WithDescription(t *testing.T) {
	repo := tracker.NewProjectRepository(testPool)
	desc := "My description"

	got, err := repo.Create(context.Background(), uuid.New(), uuid.New(), "Described", &desc)

	require.NoError(t, err)
	require.NotNil(t, got.Description)
	assert.Equal(t, desc, *got.Description)
}

func TestProjectRepository_Create_DuplicateID_ReturnsError(t *testing.T) {
	repo := tracker.NewProjectRepository(testPool)
	ctx := context.Background()
	id := uuid.New()

	_, err := repo.Create(ctx, id, uuid.New(), "First", nil)
	require.NoError(t, err)

	_, err = repo.Create(ctx, id, uuid.New(), "Second", nil)
	require.Error(t, err) // PK violation
}
