package invitation

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
)

// InvitationService implements the domain logic for managing workspace invitations.
type InvitationService struct {
	repository  InvitationRepository
	memberAdder WorkspaceMemberAdder
}

// NewInvitationService creates an InvitationService wired to the given repository and member adder.
func NewInvitationService(repository InvitationRepository, memberAdder WorkspaceMemberAdder) *InvitationService {
	return &InvitationService{repository: repository, memberAdder: memberAdder}
}

// AcceptCommand holds all inputs needed to accept an invitation.
type AcceptCommand struct {
	Token              string
	AcceptingUserID    uuid.UUID
	AcceptingUserEmail string
}

// Invite generates a token, persists a pending invitation, and emits InvitationCreatedEvent.
func (s *InvitationService) Invite(ctx context.Context, command InviteCommand) (Invitation, error) {
	rawToken, tokenHash, err := generateToken()
	if err != nil {
		return Invitation{}, fmt.Errorf("invite: %w", err)
	}

	invitation, err := command.ToInvitation(uuid.New(), tokenHash, time.Now())
	if err != nil {
		return Invitation{}, fmt.Errorf("invite: %w", err)
	}
	invitation.Token = rawToken

	result, err := s.repository.Create(ctx, invitation)
	if err != nil {
		return Invitation{}, fmt.Errorf("invite: %w", err)
	}

	result.Token = rawToken
	if err := result.EmitCreated(ctx); err != nil {
		slog.Error("publish invitation created event", "error", err)
	}
	return result, nil
}

// Accept validates the token (existence, status, expiry, email match), adds the
// accepting user to the workspace, and marks the invitation accepted.
func (s *InvitationService) Accept(ctx context.Context, command AcceptCommand) error {
	tokenHash := hashToken(command.Token)

	invite, err := s.repository.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		return fmt.Errorf("accept invitation: %w", err)
	}
	if invite == nil {
		return fmt.Errorf("accept invitation: %w", ErrInvitationNotFound)
	}
	if invite.Status != StatusPending {
		return fmt.Errorf("accept invitation: %w", ErrInvitationAlreadyAccepted)
	}
	if time.Now().After(invite.ExpiresAt) {
		return fmt.Errorf("accept invitation: %w", ErrInvitationExpired)
	}
	if !strings.EqualFold(invite.Email, command.AcceptingUserEmail) {
		return fmt.Errorf("accept invitation: %w", ErrEmailMismatch)
	}

	if err := s.memberAdder.AddMember(ctx, invite.WorkspaceID, command.AcceptingUserID); err != nil {
		return fmt.Errorf("accept invitation: %w", err)
	}

	if err := s.repository.MarkAccepted(ctx, invite.ID, time.Now()); err != nil {
		return fmt.Errorf("accept invitation: %w", err)
	}

	return nil
}
