//go:build !integration

package workspace_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wlindb/issue-tracker/internal/domain/tracker/workspace"
)

func Test_New_NilID_ReturnsError(t *testing.T) {
	_, err := workspace.New(uuid.Nil, "Acme", uuid.New())

	require.Error(t, err)
	assert.ErrorIs(t, err, workspace.ErrInvalidWorkspace)
}

func Test_New_EmptyName_ReturnsError(t *testing.T) {
	_, err := workspace.New(uuid.New(), "", uuid.New())

	require.Error(t, err)
	assert.ErrorIs(t, err, workspace.ErrInvalidWorkspace)
}

func Test_New_NilOwnerID_ReturnsError(t *testing.T) {
	_, err := workspace.New(uuid.New(), "Acme", uuid.Nil)

	require.Error(t, err)
	assert.ErrorIs(t, err, workspace.ErrInvalidWorkspace)
}
