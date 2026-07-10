//go:build !integration

package user_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wlindb/issue-tracker/internal/domain/tracker/user"
)

func Test_ToUser_NilID_ReturnsError(t *testing.T) {
	command := user.UpsertUserCommand{ID: uuid.Nil, Email: "jane@example.com", Name: "Jane Doe"}

	actual, err := command.ToUser()

	require.Error(t, err)
	assert.ErrorIs(t, err, user.ErrInvalidUser)
	assert.Zero(t, actual)
}

func Test_ToUser_EmptyEmail_ReturnsError(t *testing.T) {
	command := user.UpsertUserCommand{ID: uuid.New(), Email: "", Name: "Jane Doe"}

	actual, err := command.ToUser()

	require.Error(t, err)
	assert.ErrorIs(t, err, user.ErrInvalidUser)
	assert.Zero(t, actual)
}

func Test_ToUser_NameProvided_UsesName(t *testing.T) {
	id := uuid.New()
	command := user.UpsertUserCommand{ID: id, Email: "jane@example.com", Name: "Jane Doe"}

	actual, err := command.ToUser()

	require.NoError(t, err)
	assert.Equal(t, id, actual.ID)
	assert.Equal(t, "jane@example.com", actual.Email)
	assert.Equal(t, "Jane Doe", actual.Name)
}

func Test_ToUser_NameMissing_SuccessfulUserWithEmptyName(t *testing.T) {
	command := user.UpsertUserCommand{ID: uuid.New(), Email: "jane@example.com"}

	actual, err := command.ToUser()

	require.NoError(t, err)
	assert.Equal(t, "", actual.Name)
}
