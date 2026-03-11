package api

import (
	"context"

	"github.com/google/uuid"
)

type contextKey string

const userIDKey contextKey = "userID"

func userIDFromContext(ctx context.Context) uuid.UUID {
	id, _ := ctx.Value(userIDKey).(uuid.UUID)
	return id
}

func withUserID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

// WithUserID is the exported form for use in tests and middleware.
func WithUserID(ctx context.Context, id uuid.UUID) context.Context {
	return withUserID(ctx, id)
}
