package search

import (
	"context"
	"embed"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrations embed.FS

// Migrate runs goose migrations for the search module against the given pool.
func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	db := stdlib.OpenDBFromPool(pool)

	goose.SetBaseFS(migrations)
	defer goose.SetBaseFS(nil)

	if err := goose.SetDialect("postgres"); err != nil {
		return errors.Join(fmt.Errorf("search migrate set dialect: %w", err), db.Close())
	}
	if err := goose.UpContext(ctx, db, "migrations"); err != nil {
		return errors.Join(fmt.Errorf("search migrate up: %w", err), db.Close())
	}
	if err := db.Close(); err != nil {
		return fmt.Errorf("search migrate close db: %w", err)
	}
	return nil
}
