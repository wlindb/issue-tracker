//go:build !integration

package invitation_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wlindb/issue-tracker/internal/domain/tracker/invitation"
)

func Test_ToInvitation_EmptyEmail_ReturnsError(t *testing.T) {
	command := invitation.InviteCommand{
		WorkspaceID:     uuid.New(),
		InvitedByUserID: uuid.New(),
		Email:           "",
	}

	actual, err := command.ToInvitation(uuid.New(), "tokenhash", time.Now())

	require.Error(t, err)
	assert.ErrorIs(t, err, invitation.ErrInvalidInvitation)
	assert.Zero(t, actual)
}

func Test_ToInvitation_NilWorkspaceID_ReturnsError(t *testing.T) {
	command := invitation.InviteCommand{
		WorkspaceID:     uuid.Nil,
		InvitedByUserID: uuid.New(),
		Email:           "invitee@example.com",
	}

	actual, err := command.ToInvitation(uuid.New(), "tokenhash", time.Now())

	require.Error(t, err)
	assert.ErrorIs(t, err, invitation.ErrInvalidInvitation)
	assert.Zero(t, actual)
}

func Test_ToInvitation_NilInvitedByUserID_ReturnsError(t *testing.T) {
	command := invitation.InviteCommand{
		WorkspaceID:     uuid.New(),
		InvitedByUserID: uuid.Nil,
		Email:           "invitee@example.com",
	}

	actual, err := command.ToInvitation(uuid.New(), "tokenhash", time.Now())

	require.Error(t, err)
	assert.ErrorIs(t, err, invitation.ErrInvalidInvitation)
	assert.Zero(t, actual)
}

func Test_ToInvitation_ValidCommand_ReturnsPendingInvitationWithSevenDayExpiry(t *testing.T) {
	id := uuid.New()
	workspaceID := uuid.New()
	invitedByUserID := uuid.New()
	now := time.Now().UTC()
	command := invitation.InviteCommand{
		WorkspaceID:     workspaceID,
		InvitedByUserID: invitedByUserID,
		Email:           "invitee@example.com",
	}

	actual, err := command.ToInvitation(id, "tokenhash", now)

	require.NoError(t, err)
	assert.Equal(t, id, actual.ID)
	assert.Equal(t, workspaceID, actual.WorkspaceID)
	assert.Equal(t, invitedByUserID, actual.InvitedByUserID)
	assert.Equal(t, "invitee@example.com", actual.Email)
	assert.Equal(t, "tokenhash", actual.TokenHash)
	assert.Equal(t, invitation.StatusPending, actual.Status)
	assert.Equal(t, now, actual.CreatedAt)
	assert.Equal(t, now.Add(7*24*time.Hour), actual.ExpiresAt)
}
