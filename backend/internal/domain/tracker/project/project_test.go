//go:build !integration

package project_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wlindb/issue-tracker/internal/domain/tracker/project"
)

func Test_NewListProjectQuery_NilLimit_UsesDefaultLimit(t *testing.T) {
	actual := project.NewListProjectQuery(nil, nil)

	require.NotNil(t, actual.Limit)
	assert.Equal(t, 20, *actual.Limit)
}

func Test_NewListProjectQuery_ExplicitLimit_UsesProvidedLimit(t *testing.T) {
	limit := 5

	actual := project.NewListProjectQuery(nil, &limit)

	require.NotNil(t, actual.Limit)
	assert.Equal(t, 5, *actual.Limit)
}

func Test_NewListProjectQuery_NilCursor_CursorIsNil(t *testing.T) {
	actual := project.NewListProjectQuery(nil, nil)

	assert.Nil(t, actual.Cursor)
}

func Test_New_NilID_ReturnsError(t *testing.T) {
	actual, err := project.New(uuid.Nil, "my-project", "My Project", nil, uuid.New())

	require.Error(t, err)
	assert.ErrorIs(t, err, project.ErrInvalidProject)
	assert.Zero(t, actual)
}

func Test_New_EmptyIdentifier_ReturnsError(t *testing.T) {
	actual, err := project.New(uuid.New(), "", "My Project", nil, uuid.New())

	require.Error(t, err)
	assert.ErrorIs(t, err, project.ErrInvalidProject)
	assert.Zero(t, actual)
}

func Test_New_EmptyName_ReturnsError(t *testing.T) {
	actual, err := project.New(uuid.New(), "my-project", "", nil, uuid.New())

	require.Error(t, err)
	assert.ErrorIs(t, err, project.ErrInvalidProject)
	assert.Zero(t, actual)
}

func Test_New_NilOwnerID_ReturnsError(t *testing.T) {
	actual, err := project.New(uuid.New(), "my-project", "My Project", nil, uuid.Nil)

	require.Error(t, err)
	assert.ErrorIs(t, err, project.ErrInvalidProject)
	assert.Zero(t, actual)
}

func Test_New_ValidArguments_ReturnsProject(t *testing.T) {
	id := uuid.New()
	ownerID := uuid.New()

	actual, err := project.New(id, "my-project", "My Project", nil, ownerID)

	require.NoError(t, err)
	assert.Equal(t, id, actual.ID)
	assert.Equal(t, "my-project", actual.Identifier)
	assert.Equal(t, "My Project", actual.Name)
	assert.Nil(t, actual.Description)
	assert.Equal(t, ownerID, actual.OwnerID)
	assert.False(t, actual.CreatedAt.IsZero())
	assert.False(t, actual.UpdatedAt.IsZero())
}

func Test_ToProject_ValidCommand_ReturnsProject(t *testing.T) {
	ownerID := uuid.New()
	id := uuid.New()
	cmd := project.CreateProjectCommand{Name: "My Project", OwnerID: ownerID}

	actual, err := cmd.ToProject(id, cmd.Slugify)

	require.NoError(t, err)
	assert.Equal(t, id, actual.ID)
	assert.Equal(t, "my-project", actual.Identifier)
	assert.Equal(t, "My Project", actual.Name)
	assert.Equal(t, ownerID, actual.OwnerID)
	assert.False(t, actual.CreatedAt.IsZero())
}

func Test_ToProject_EmptyName_ReturnsError(t *testing.T) {
	cmd := project.CreateProjectCommand{Name: "", OwnerID: uuid.New()}

	actual, err := cmd.ToProject(uuid.New(), cmd.Slugify)

	require.Error(t, err)
	assert.ErrorIs(t, err, project.ErrInvalidProject)
	assert.Zero(t, actual)
}

func Test_ToProject_NilOwnerID_ReturnsError(t *testing.T) {
	cmd := project.CreateProjectCommand{Name: "My Project", OwnerID: uuid.Nil}

	actual, err := cmd.ToProject(uuid.New(), cmd.Slugify)

	require.Error(t, err)
	assert.ErrorIs(t, err, project.ErrInvalidProject)
	assert.Zero(t, actual)
}
