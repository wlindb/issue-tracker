package invitation

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Status is the lifecycle state of an Invitation.
type Status string

const (
	StatusPending  Status = "pending"
	StatusAccepted Status = "accepted"
)

const (
	expiry          = 7 * 24 * time.Hour
	tokenByteLength = 32
)

var (
	ErrInvalidInvitation         = errors.New("invalid invitation")
	ErrInvitationNotFound        = errors.New("invitation not found")
	ErrInvitationExpired         = errors.New("invitation expired")
	ErrInvitationAlreadyAccepted = errors.New("invitation already accepted")
	ErrEmailMismatch             = errors.New("accepting user's email does not match invitation")
)

// Invitation is the domain entity for a workspace invitation.
type Invitation struct {
	ID              uuid.UUID
	WorkspaceID     uuid.UUID
	Email           string
	InvitedByUserID uuid.UUID
	TokenHash       string
	Token           string // transient only, never persisted; carried through to InvitationCreatedEvent
	Status          Status
	ExpiresAt       time.Time
	CreatedAt       time.Time
	AcceptedAt      *time.Time
}

// InviteCommand holds all inputs needed to invite an email address to a workspace.
type InviteCommand struct {
	WorkspaceID     uuid.UUID
	InvitedByUserID uuid.UUID
	Email           string
}

// ToInvitation validates the command and builds a pending Invitation.
// Returns ErrInvalidInvitation if the command contains invalid data.
func (c InviteCommand) ToInvitation(id uuid.UUID, tokenHash string, now time.Time) (Invitation, error) {
	if c.WorkspaceID == uuid.Nil {
		return Invitation{}, ErrInvalidInvitation
	}
	if c.InvitedByUserID == uuid.Nil {
		return Invitation{}, ErrInvalidInvitation
	}
	if c.Email == "" {
		return Invitation{}, ErrInvalidInvitation
	}
	return Invitation{
		ID:              id,
		WorkspaceID:     c.WorkspaceID,
		Email:           c.Email,
		InvitedByUserID: c.InvitedByUserID,
		TokenHash:       tokenHash,
		Status:          StatusPending,
		ExpiresAt:       now.Add(expiry),
		CreatedAt:       now,
	}, nil
}

// InvitationRepository defines the persistence operations for invitations.
type InvitationRepository interface {
	Create(ctx context.Context, invitation Invitation) (Invitation, error)
	GetByTokenHash(ctx context.Context, tokenHash string) (*Invitation, error)
	MarkAccepted(ctx context.Context, id uuid.UUID, acceptedAt time.Time) error
}

// WorkspaceMemberAdder is the minimal workspace dependency this package needs,
// satisfied by workspace.WorkspaceService.
type WorkspaceMemberAdder interface {
	AddMember(ctx context.Context, workspaceID uuid.UUID, userID uuid.UUID) error
}

// generateToken returns a random raw token and its sha256 hex hash.
func generateToken() (rawToken string, tokenHash string, err error) {
	raw := make([]byte, tokenByteLength)
	if _, err := rand.Read(raw); err != nil {
		return "", "", fmt.Errorf("generate token: %w", err)
	}
	rawToken = hex.EncodeToString(raw)
	return rawToken, hashToken(rawToken), nil
}

// hashToken hashes a raw token the same way generateToken does, for accept-time lookup.
func hashToken(rawToken string) string {
	sum := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(sum[:])
}
