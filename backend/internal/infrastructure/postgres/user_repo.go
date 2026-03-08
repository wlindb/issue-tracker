package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wlindb/issue-tracker/internal/domain"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

const createUserSQL = `
INSERT INTO users (email, name, password_hash)
VALUES ($1, $2, $3)
RETURNING id, email, name, password_hash, created_at, updated_at`

func (r *UserRepo) Create(ctx context.Context, email, name, passwordHash string) (*domain.User, error) {
	row := r.pool.QueryRow(ctx, createUserSQL, email, name, passwordHash)
	return scanUser(row)
}

const getUserByEmailSQL = `
SELECT id, email, name, password_hash, created_at, updated_at
FROM users WHERE email = $1`

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	row := r.pool.QueryRow(ctx, getUserByEmailSQL, email)
	return scanUser(row)
}

type scanner interface {
	Scan(dest ...any) error
}

func scanUser(row scanner) (*domain.User, error) {
	var u domain.User
	if err := row.Scan(&u.ID, &u.Email, &u.Name, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, fmt.Errorf("scanning user: %w", err)
	}
	return &u, nil
}
