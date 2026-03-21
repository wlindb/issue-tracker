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

func Test_Create_NoDescription_SuccessfulProjectCreation(t *testing.T) {
	repo := tracker.NewProjectRepository(testPool)
	id, ownerID := uuid.New(), uuid.New()

	actual, err := repo.Create(context.Background(), id, ownerID, "Acme", nil)

	require.NoError(t, err)
	require.NotNil(t, actual)
	assert.Equal(t, id, actual.ID)
	assert.Equal(t, ownerID, actual.OwnerID)
	assert.Equal(t, "Acme", actual.Name)
	assert.Nil(t, actual.Description)
	assert.False(t, actual.CreatedAt.IsZero())
	assert.False(t, actual.UpdatedAt.IsZero())
}

func Test_Create_WithDescription_SuccessfulProjectCreation(t *testing.T) {
	repo := tracker.NewProjectRepository(testPool)
	desc := "My description"

	actual, err := repo.Create(context.Background(), uuid.New(), uuid.New(), "Described", &desc)

	require.NoError(t, err)
	require.NotNil(t, actual.Description)
	assert.Equal(t, desc, *actual.Description)
}

func Test_Create_DuplicateID_ReturnsError(t *testing.T) {
	repo := tracker.NewProjectRepository(testPool)
	ctx := context.Background()
	id := uuid.New()

	_, err := repo.Create(ctx, id, uuid.New(), "First", nil)
	require.NoError(t, err)

	_, err = repo.Create(ctx, id, uuid.New(), "Second", nil)
	require.Error(t, err) // PK violation
}
