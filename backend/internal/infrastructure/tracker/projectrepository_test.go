package tracker

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	projectdomain "github.com/wlindb/issue-tracker/internal/domain/tracker/project"
	trackerdb "github.com/wlindb/issue-tracker/internal/infrastructure/tracker/generated"
)

// Compile-time: *ProjectRepository must satisfy domain interface.
var _ projectdomain.ProjectRepository = (*ProjectRepository)(nil)

func TestRowToProject_NilDescription(t *testing.T) {
	id, ownerID := uuid.New(), uuid.New()
	now := time.Now().UTC()
	row := trackerdb.Project{
		ID:          id,
		OwnerID:     ownerID,
		Name:        "Test",
		Description: pgtype.Text{Valid: false},
		CreatedAt:   pgtype.Timestamptz{Time: now, Valid: true},
		UpdatedAt:   pgtype.Timestamptz{Time: now, Valid: true},
	}
	got := rowToProject(row)
	require.NotNil(t, got)
	assert.Equal(t, id, got.ID)
	assert.Equal(t, ownerID, got.OwnerID)
	assert.Equal(t, "Test", got.Name)
	assert.Nil(t, got.Description)
	assert.Equal(t, now, got.CreatedAt)
}

func TestRowToProject_WithDescription(t *testing.T) {
	desc := "hello"
	now := time.Now().UTC()
	row := trackerdb.Project{
		ID:          uuid.New(),
		OwnerID:     uuid.New(),
		Name:        "Test",
		Description: pgtype.Text{String: desc, Valid: true},
		CreatedAt:   pgtype.Timestamptz{Time: now, Valid: true},
		UpdatedAt:   pgtype.Timestamptz{Time: now, Valid: true},
	}
	got := rowToProject(row)
	require.NotNil(t, got.Description)
	assert.Equal(t, desc, *got.Description)
}
