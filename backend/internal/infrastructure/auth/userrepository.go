package auth

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	authdomain "github.com/wlindb/issue-tracker/internal/domain/auth"
	authdb "github.com/wlindb/issue-tracker/internal/infrastructure/auth/generated"
)

type UserRepository struct {
	q *authdb.Queries
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{q: authdb.New(pool)}
}

func (r *UserRepository) Create(ctx context.Context, email, name, passwordHash string) (*authdomain.User, error) {
	row, err := r.q.CreateUser(ctx, authdb.CreateUserParams{
		Email:        email,
		Name:         name,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return toDomainUser(row), nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*authdomain.User, error) {
	row, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return toDomainUser(row), nil
}

func toDomainUser(u authdb.User) *authdomain.User {
	return &authdomain.User{
		ID:           u.ID,
		Email:        u.Email,
		Name:         u.Name,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}
