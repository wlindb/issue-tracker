package api

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type contextKey string

const userIDKey contextKey = "userID"

// errMissingUserID is returned when no user ID is present in the context.
var errMissingUserID = errors.New("missing user ID in context")

func userIDFromContext(ctx context.Context) (uuid.UUID, error) {
	id, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok || id == uuid.Nil {
		return uuid.Nil, errMissingUserID
	}
	return id, nil
}

func withUserID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

// WithUserID is the exported form for use in tests and middleware.
func WithUserID(ctx context.Context, id uuid.UUID) context.Context {
	return withUserID(ctx, id)
}

// UserIDFromContext is the exported form for use in tests.
func UserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	return userIDFromContext(ctx)
}
