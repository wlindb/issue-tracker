package project_test

import (
	"testing"

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
