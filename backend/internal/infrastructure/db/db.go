package db

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// WithAppSessionVars returns a pool configurator that sets app.workspace_id and
// app.user_id on each acquired connection and clears them on release.
// extractWorkspaceID and extractUserID return the UUID for each ID, or an error
// if the ID is absent or invalid (e.g. uuid.Nil). The session variable is only
// set when the extractor returns a nil error.
func WithAppSessionVars(
	extractWorkspaceID func(context.Context) (uuid.UUID, error),
	extractUserID func(context.Context) (uuid.UUID, error),
) func(*pgxpool.Config) {
	return func(config *pgxpool.Config) {
		config.PrepareConn = func(ctx context.Context, conn *pgx.Conn) (bool, error) {
			if id, err := extractWorkspaceID(ctx); err == nil {
				if _, err := conn.Exec(ctx, "SELECT set_config('app.workspace_id', $1, false)", id.String()); err != nil {
					return false, fmt.Errorf("set app.workspace_id: %w", err)
				}
			}
			if id, err := extractUserID(ctx); err == nil {
				if _, err := conn.Exec(ctx, "SELECT set_config('app.user_id', $1, false)", id.String()); err != nil {
					return false, fmt.Errorf("set app.user_id: %w", err)
				}
			}
			return true, nil
		}
		config.AfterRelease = func(conn *pgx.Conn) bool {
			for _, name := range []string{"app.workspace_id", "app.user_id"} {
				if _, err := conn.Exec(context.Background(), "SELECT set_config($1, '', false)", name); err != nil {
					return false
				}
			}
			return true
		}
	}
}

// WithAppRole returns a pool configurator that switches to the given role on
// every acquired connection and resets it on release. This ensures queries run
// as a non-superuser role so that PostgreSQL RLS policies are enforced.
// Role names cannot be parameterised in PostgreSQL; the value comes from
// internal config, not user input, so string concatenation is safe here.
func WithAppRole(role string) func(*pgxpool.Config) {
	return func(config *pgxpool.Config) {
		existingPrepare := config.PrepareConn
		config.PrepareConn = func(ctx context.Context, conn *pgx.Conn) (bool, error) {
			if existingPrepare != nil {
				if ok, err := existingPrepare(ctx, conn); !ok || err != nil {
					return ok, err
				}
			}
			if _, err := conn.Exec(ctx, "SET ROLE "+role); err != nil {
				return false, fmt.Errorf("set role %s: %w", role, err)
			}
			return true, nil
		}
		existingRelease := config.AfterRelease
		config.AfterRelease = func(conn *pgx.Conn) bool {
			if _, err := conn.Exec(context.Background(), "RESET ROLE"); err != nil {
				return false
			}
			if existingRelease != nil {
				return existingRelease(conn)
			}
			return true
		}
	}
}

func New(ctx context.Context, databaseURL string, configure ...func(*pgxpool.Config)) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse pool config: %w", err)
	}
	for _, fn := range configure {
		fn(config)
	}
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}
	return pool, nil
}
