package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrInvalidUser = errors.New("invalid user")
)

// User is a local record of an authenticated Keycloak identity.
type User struct {
	ID    uuid.UUID
	Email string
	Name  string
}

// UpsertUserCommand holds the profile fields extracted from the authenticated
// user's JWT claims, keyed by the immutable `sub` claim.
type UpsertUserCommand struct {
	ID    uuid.UUID
	Email string
	Name  string
}

// ToUser builds a User from the command. Returns ErrInvalidUser if the
// command contains invalid data.
func (c UpsertUserCommand) ToUser() (User, error) {
	if c.ID == uuid.Nil {
		return User{}, ErrInvalidUser
	}
	if c.Email == "" {
		return User{}, ErrInvalidUser
	}
	return User(c), nil
}

// UserRepository upserts local user records keyed by their immutable ID.
type UserRepository interface {
	Upsert(ctx context.Context, user User) (User, error)
}
