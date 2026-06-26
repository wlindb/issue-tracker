// Package tracker
package tracker

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
)

type UnitOfWork struct {
	db *pgxpool.Pool
}

func NewUoW(db *pgxpool.Pool) *UnitOfWork {
	return &UnitOfWork{db: db}
}

func (u *UnitOfWork) RunInTx(ctx context.Context, fn func(issue.Repositories) error) error {
	tx, err := u.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	stores := issue.Repositories{
		Issues: NewIssueRepository(tx),
	}

	if err := fn(stores); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}
