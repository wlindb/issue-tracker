package tracker

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	trackerdb "github.com/wlindb/issue-tracker/internal/infrastructure/tracker/generated"
)

func Test_RowToProject_NilDescription_SetsDescriptionNil(t *testing.T) {
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

	actual := rowToProject(row)

	require.NotNil(t, actual)
	assert.Equal(t, id, actual.ID)
	assert.Equal(t, ownerID, actual.OwnerID)
	assert.Equal(t, "Test", actual.Name)
	assert.Nil(t, actual.Description)
	assert.Equal(t, now, actual.CreatedAt)
}

func Test_RowToProject_WithDescription_SetsDescription(t *testing.T) {
	description := "hello"
	now := time.Now().UTC()
	row := trackerdb.Project{
		ID:          uuid.New(),
		OwnerID:     uuid.New(),
		Name:        "Test",
		Description: pgtype.Text{String: description, Valid: true},
		CreatedAt:   pgtype.Timestamptz{Time: now, Valid: true},
		UpdatedAt:   pgtype.Timestamptz{Time: now, Valid: true},
	}

	actual := rowToProject(row)

	require.NotNil(t, actual.Description)
	assert.Equal(t, description, *actual.Description)
}

func Test_RowsToProjects_Empty_ReturnsEmptySlice(t *testing.T) {
	actual := rowsToProjects([]trackerdb.Project{})

	assert.NotNil(t, actual)
	assert.Empty(t, actual)
}

func Test_RowsToProjects_MultipleRows_ReturnsMappedProjects(t *testing.T) {
	now := time.Now().UTC()
	id1, id2 := uuid.New(), uuid.New()
	rows := []trackerdb.Project{
		{
			ID:        id1,
			OwnerID:   uuid.New(),
			Name:      "Alpha",
			CreatedAt: pgtype.Timestamptz{Time: now, Valid: true},
			UpdatedAt: pgtype.Timestamptz{Time: now, Valid: true},
		},
		{
			ID:        id2,
			OwnerID:   uuid.New(),
			Name:      "Beta",
			CreatedAt: pgtype.Timestamptz{Time: now, Valid: true},
			UpdatedAt: pgtype.Timestamptz{Time: now, Valid: true},
		},
	}

	actual := rowsToProjects(rows)

	require.Len(t, actual, 2)
	assert.Equal(t, id1, actual[0].ID)
	assert.Equal(t, "Alpha", actual[0].Name)
	assert.Equal(t, id2, actual[1].ID)
	assert.Equal(t, "Beta", actual[1].Name)
}
