package api

import (
	"context"
	"errors"

	key "github.com/wlindb/issue-tracker/internal/pkg/context"

	"github.com/google/uuid"
)

// errMissingUserID is returned when no user ID is present in the context.
var errMissingUserID = errors.New("missing user ID in context")

// errMissingWorkspaceID is returned when no workspace ID is present in the context.
var errMissingWorkspaceID = errors.New("missing workspace ID in context")

// WithUserID is the exported form for use in tests and middleware.
func WithUserID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, key.UserID, id)
}

// UserIDFromContext is the exported form for use in tests.
func UserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	id, ok := ctx.Value(key.UserID).(uuid.UUID)
	if !ok || id == uuid.Nil {
		return uuid.Nil, errMissingUserID
	}
	return id, nil
}

// WithWorkspaceID stores a workspace ID in the context for use by the pgx pool hooks.
func WithWorkspaceID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, key.WorkspaceID, id)
}

// WorkspaceIDFromContext retrieves the workspace ID stored by WithWorkspaceID.
// Returns errMissingWorkspaceID if not set or uuid.Nil.
func WorkspaceIDFromContext(ctx context.Context) (uuid.UUID, error) {
	id, ok := ctx.Value(key.WorkspaceID).(uuid.UUID)
	if !ok || id == uuid.Nil {
		return uuid.Nil, errMissingWorkspaceID
	}
	return id, nil
}
