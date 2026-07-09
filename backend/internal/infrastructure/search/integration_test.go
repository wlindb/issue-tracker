//go:build integration

package search_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	infradb "github.com/wlindb/issue-tracker/internal/infrastructure/db"
	"github.com/wlindb/issue-tracker/internal/infrastructure/search"
)

var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()
	dsn, terminate, err := startPostgres(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "start postgres: %v\n", err)
		os.Exit(1)
	}

	migrationPool, err := infradb.New(ctx, dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "start postgres: connect migration pool: %v\n", err)
		os.Exit(1)
	}
	if err := search.Migrate(ctx, migrationPool); err != nil {
		fmt.Fprintf(os.Stderr, "migrate: %v\n", err)
		os.Exit(1)
	}
	migrationPool.Close()

	testPool, err = infradb.New(ctx, dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "start postgres: connect app pool: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()
	testPool.Close()
	terminate()
	os.Exit(code)
}

func startPostgres(ctx context.Context) (string, func(), error) {
	container, err := postgres.Run(ctx,
		"timescale/timescaledb-ha:pg18",
		postgres.WithDatabase("test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
		),
	)
	if err != nil {
		return "", nil, fmt.Errorf("start container: %w", err)
	}
	dsn, err := container.ConnectionString(ctx)
	if err != nil {
		return "", nil, errors.Join(fmt.Errorf("connection string: %w", err), container.Terminate(ctx))
	}
	return dsn, func() {
		if err := container.Terminate(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "terminate container: %v\n", err)
		}
	}, nil
}
