package user

import (
	"context"
	"fmt"
)

type UserService struct {
	repository UserRepository
}

func NewUserService(repository UserRepository) *UserService {
	return &UserService{repository: repository}
}

// Upsert creates the user if one does not already exist, or updates its
// mutable profile fields if it does. The operation is idempotent: repeated
// calls for the same ID never create duplicate records.
func (s *UserService) Upsert(ctx context.Context, command UpsertUserCommand) (User, error) {
	user, err := command.ToUser()
	if err != nil {
		return User{}, fmt.Errorf("upsert user: %w", err)
	}
	result, err := s.repository.Upsert(ctx, user)
	if err != nil {
		return User{}, fmt.Errorf("upsert user: %w", err)
	}
	return result, nil
}
