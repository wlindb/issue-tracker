package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// WithAppSessionVars returns a pool configurator that sets app.workspace_id and
// app.user_id on each acquired connection and clears them on release.
// extractWorkspaceID and extractUserID return the string form of each ID and
// whether it is present in the context.
func WithAppSessionVars(
	extractWorkspaceID func(context.Context) (string, bool),
	extractUserID func(context.Context) (string, bool),
) func(*pgxpool.Config) {
	return func(config *pgxpool.Config) {
		config.PrepareConn = func(ctx context.Context, conn *pgx.Conn) (bool, error) {
			if id, ok := extractWorkspaceID(ctx); ok {
				if _, err := conn.Exec(ctx, "SELECT set_config('app.workspace_id', $1, false)", id); err != nil {
					return false, fmt.Errorf("set app.workspace_id: %w", err)
				}
			}
			if id, ok := extractUserID(ctx); ok {
				if _, err := conn.Exec(ctx, "SELECT set_config('app.user_id', $1, false)", id); err != nil {
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
