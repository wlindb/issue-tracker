package api

import (
	"context"

	"github.com/google/uuid"
)

type contextKey string

const callerIDKey contextKey = "callerID"

func callerIDFromContext(ctx context.Context) uuid.UUID {
	id, _ := ctx.Value(callerIDKey).(uuid.UUID)
	return id
}

func withCallerID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, callerIDKey, id)
}

// WithCallerID is the exported form for use in tests and middleware.
func WithCallerID(ctx context.Context, id uuid.UUID) context.Context {
	return withCallerID(ctx, id)
}
