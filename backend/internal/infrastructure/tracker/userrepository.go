package tracker

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	userdomain "github.com/wlindb/issue-tracker/internal/domain/tracker/user"
	trackerdb "github.com/wlindb/issue-tracker/internal/infrastructure/tracker/generated"
)

// Compile-time: *UserRepository must satisfy domain interface.
var _ userdomain.UserRepository = (*UserRepository)(nil)

// UserRepository is a PostgreSQL-backed implementation of userdomain.UserRepository.
type UserRepository struct {
	queries *trackerdb.Queries
}

// NewUserRepository returns a UserRepository backed by the given pool.
func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{queries: trackerdb.New(pool)}
}

// Upsert inserts a new user row, or updates its mutable profile fields if a
// row with the same ID already exists, and returns the resulting domain model.
func (r *UserRepository) Upsert(ctx context.Context, user userdomain.User) (userdomain.User, error) {
	row, err := r.queries.UpsertUser(ctx, upsertUserParamsFromDomain(user))
	if err != nil {
		return userdomain.User{}, fmt.Errorf("upsert user: %w", err)
	}
	return userToDomain(row), nil
}
