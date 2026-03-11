package infrastructure

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wlindb/issue-tracker/internal/domain"
	"github.com/wlindb/issue-tracker/internal/infrastructure/generated"
)

type UserRepository struct {
	q *generated.Queries
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{q: generated.New(pool)}
}

func (r *UserRepository) Create(ctx context.Context, email, name, passwordHash string) (*domain.User, error) {
	row, err := r.q.CreateUser(ctx, generated.CreateUserParams{
		Email:        email,
		Name:         name,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return toDomainUser(row), nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	row, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return toDomainUser(row), nil
}

func toDomainUser(u generated.User) *domain.User {
	return &domain.User{
		ID:           u.ID,
		Email:        u.Email,
		Name:         u.Name,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}
