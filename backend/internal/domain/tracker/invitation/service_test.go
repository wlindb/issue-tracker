//go:build !integration

package invitation_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/wlindb/issue-tracker/internal/domain/tracker/invitation"
)

type mockInvitationRepository struct {
	mock.Mock
}

func (m *mockInvitationRepository) Create(ctx context.Context, i invitation.Invitation) (invitation.Invitation, error) {
	args := m.Called(ctx, i)
	return args.Get(0).(invitation.Invitation), args.Error(1)
}

func (m *mockInvitationRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*invitation.Invitation, error) {
	args := m.Called(ctx, tokenHash)
	if i, ok := args.Get(0).(*invitation.Invitation); ok {
		return i, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockInvitationRepository) MarkAccepted(ctx context.Context, id uuid.UUID, acceptedAt time.Time) error {
	args := m.Called(ctx, id, acceptedAt)
	return args.Error(0)
}

type mockWorkspaceMemberAdder struct {
	mock.Mock
}

func (m *mockWorkspaceMemberAdder) AddMember(ctx context.Context, workspaceID uuid.UUID, userID uuid.UUID) error {
	args := m.Called(ctx, workspaceID, userID)
	return args.Error(0)
}

func sha256Hex(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

func Test_Invite_ValidCommand_ReturnsInvitation(t *testing.T) {
	repository := &mockInvitationRepository{}
	memberAdder := &mockWorkspaceMemberAdder{}
	service := invitation.NewInvitationService(repository, memberAdder)

	workspaceID := uuid.New()
	invitedByUserID := uuid.New()
	command := invitation.InviteCommand{
		WorkspaceID:     workspaceID,
		InvitedByUserID: invitedByUserID,
		Email:           "invitee@example.com",
	}

	var captured invitation.Invitation
	repository.On("Create", mock.Anything, mock.MatchedBy(func(i invitation.Invitation) bool {
		captured = i
		return i.WorkspaceID == workspaceID && i.InvitedByUserID == invitedByUserID &&
			i.Email == command.Email && i.TokenHash != "" && i.Token != "" &&
			i.Status == invitation.StatusPending
	})).Return(invitation.Invitation{
		ID:              uuid.New(),
		WorkspaceID:     workspaceID,
		InvitedByUserID: invitedByUserID,
		Email:           command.Email,
		Status:          invitation.StatusPending,
	}, nil)

	actual, err := service.Invite(context.Background(), command)

	require.NoError(t, err)
	assert.NotEmpty(t, actual.Token)
	assert.Equal(t, sha256Hex(actual.Token), captured.TokenHash)
	assert.Equal(t, workspaceID, actual.WorkspaceID)
	assert.Equal(t, invitedByUserID, actual.InvitedByUserID)
	assert.Equal(t, command.Email, actual.Email)
	repository.AssertExpectations(t)
}

func Test_Invite_RepositoryError_ReturnsWrappedError(t *testing.T) {
	repository := &mockInvitationRepository{}
	memberAdder := &mockWorkspaceMemberAdder{}
	service := invitation.NewInvitationService(repository, memberAdder)

	command := invitation.InviteCommand{
		WorkspaceID:     uuid.New(),
		InvitedByUserID: uuid.New(),
		Email:           "invitee@example.com",
	}
	repoErr := errors.New("db down")

	repository.On("Create", mock.Anything, mock.Anything).Return(invitation.Invitation{}, repoErr)

	_, err := service.Invite(context.Background(), command)

	require.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
	repository.AssertExpectations(t)
}

func Test_Accept_ValidToken_AddsMemberAndMarksAccepted(t *testing.T) {
	repository := &mockInvitationRepository{}
	memberAdder := &mockWorkspaceMemberAdder{}
	service := invitation.NewInvitationService(repository, memberAdder)

	rawToken := "raw-token"
	tokenHash := sha256Hex(rawToken)
	workspaceID := uuid.New()
	invitationID := uuid.New()
	acceptingUserID := uuid.New()
	stored := &invitation.Invitation{
		ID:          invitationID,
		WorkspaceID: workspaceID,
		Email:       "invitee@example.com",
		TokenHash:   tokenHash,
		Status:      invitation.StatusPending,
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	}
	command := invitation.AcceptCommand{
		Token:              rawToken,
		AcceptingUserID:    acceptingUserID,
		AcceptingUserEmail: "invitee@example.com",
	}

	repository.On("GetByTokenHash", mock.Anything, tokenHash).Return(stored, nil)
	memberAdder.On("AddMember", mock.Anything, workspaceID, acceptingUserID).Return(nil)
	repository.On("MarkAccepted", mock.Anything, invitationID, mock.AnythingOfType("time.Time")).Return(nil)

	err := service.Accept(context.Background(), command)

	require.NoError(t, err)
	repository.AssertExpectations(t)
	memberAdder.AssertExpectations(t)
}

func Test_Accept_UnknownToken_ReturnsError(t *testing.T) {
	repository := &mockInvitationRepository{}
	memberAdder := &mockWorkspaceMemberAdder{}
	service := invitation.NewInvitationService(repository, memberAdder)

	rawToken := "unknown-token"
	tokenHash := sha256Hex(rawToken)
	command := invitation.AcceptCommand{
		Token:              rawToken,
		AcceptingUserID:    uuid.New(),
		AcceptingUserEmail: "invitee@example.com",
	}

	repository.On("GetByTokenHash", mock.Anything, tokenHash).Return(nil, nil)

	err := service.Accept(context.Background(), command)

	require.Error(t, err)
	assert.ErrorIs(t, err, invitation.ErrInvitationNotFound)
	repository.AssertExpectations(t)
	memberAdder.AssertNotCalled(t, "AddMember")
}

func Test_Accept_ExpiredToken_ReturnsError(t *testing.T) {
	repository := &mockInvitationRepository{}
	memberAdder := &mockWorkspaceMemberAdder{}
	service := invitation.NewInvitationService(repository, memberAdder)

	rawToken := "raw-token"
	tokenHash := sha256Hex(rawToken)
	stored := &invitation.Invitation{
		ID:          uuid.New(),
		WorkspaceID: uuid.New(),
		Email:       "invitee@example.com",
		TokenHash:   tokenHash,
		Status:      invitation.StatusPending,
		ExpiresAt:   time.Now().Add(-time.Hour),
	}
	command := invitation.AcceptCommand{
		Token:              rawToken,
		AcceptingUserID:    uuid.New(),
		AcceptingUserEmail: "invitee@example.com",
	}

	repository.On("GetByTokenHash", mock.Anything, tokenHash).Return(stored, nil)

	err := service.Accept(context.Background(), command)

	require.Error(t, err)
	assert.ErrorIs(t, err, invitation.ErrInvitationExpired)
	repository.AssertExpectations(t)
	memberAdder.AssertNotCalled(t, "AddMember")
}

func Test_Accept_AlreadyAccepted_ReturnsError(t *testing.T) {
	repository := &mockInvitationRepository{}
	memberAdder := &mockWorkspaceMemberAdder{}
	service := invitation.NewInvitationService(repository, memberAdder)

	rawToken := "raw-token"
	tokenHash := sha256Hex(rawToken)
	stored := &invitation.Invitation{
		ID:          uuid.New(),
		WorkspaceID: uuid.New(),
		Email:       "invitee@example.com",
		TokenHash:   tokenHash,
		Status:      invitation.StatusAccepted,
		ExpiresAt:   time.Now().Add(time.Hour),
	}
	command := invitation.AcceptCommand{
		Token:              rawToken,
		AcceptingUserID:    uuid.New(),
		AcceptingUserEmail: "invitee@example.com",
	}

	repository.On("GetByTokenHash", mock.Anything, tokenHash).Return(stored, nil)

	err := service.Accept(context.Background(), command)

	require.Error(t, err)
	assert.ErrorIs(t, err, invitation.ErrInvitationAlreadyAccepted)
	repository.AssertExpectations(t)
	memberAdder.AssertNotCalled(t, "AddMember")
}

func Test_Accept_EmailMismatch_ReturnsError(t *testing.T) {
	repository := &mockInvitationRepository{}
	memberAdder := &mockWorkspaceMemberAdder{}
	service := invitation.NewInvitationService(repository, memberAdder)

	rawToken := "raw-token"
	tokenHash := sha256Hex(rawToken)
	stored := &invitation.Invitation{
		ID:          uuid.New(),
		WorkspaceID: uuid.New(),
		Email:       "invitee@example.com",
		TokenHash:   tokenHash,
		Status:      invitation.StatusPending,
		ExpiresAt:   time.Now().Add(time.Hour),
	}
	command := invitation.AcceptCommand{
		Token:              rawToken,
		AcceptingUserID:    uuid.New(),
		AcceptingUserEmail: "someone-else@example.com",
	}

	repository.On("GetByTokenHash", mock.Anything, tokenHash).Return(stored, nil)

	err := service.Accept(context.Background(), command)

	require.Error(t, err)
	assert.ErrorIs(t, err, invitation.ErrEmailMismatch)
	repository.AssertExpectations(t)
	memberAdder.AssertNotCalled(t, "AddMember")
}

func Test_Accept_MemberAdderError_ReturnsWrappedError(t *testing.T) {
	repository := &mockInvitationRepository{}
	memberAdder := &mockWorkspaceMemberAdder{}
	service := invitation.NewInvitationService(repository, memberAdder)

	rawToken := "raw-token"
	tokenHash := sha256Hex(rawToken)
	workspaceID := uuid.New()
	acceptingUserID := uuid.New()
	stored := &invitation.Invitation{
		ID:          uuid.New(),
		WorkspaceID: workspaceID,
		Email:       "invitee@example.com",
		TokenHash:   tokenHash,
		Status:      invitation.StatusPending,
		ExpiresAt:   time.Now().Add(time.Hour),
	}
	command := invitation.AcceptCommand{
		Token:              rawToken,
		AcceptingUserID:    acceptingUserID,
		AcceptingUserEmail: "invitee@example.com",
	}
	memberAdderErr := errors.New("db down")

	repository.On("GetByTokenHash", mock.Anything, tokenHash).Return(stored, nil)
	memberAdder.On("AddMember", mock.Anything, workspaceID, acceptingUserID).Return(memberAdderErr)

	err := service.Accept(context.Background(), command)

	require.Error(t, err)
	assert.ErrorIs(t, err, memberAdderErr)
	repository.AssertExpectations(t)
	memberAdder.AssertExpectations(t)
	repository.AssertNotCalled(t, "MarkAccepted", mock.Anything, mock.Anything, mock.Anything)
}
