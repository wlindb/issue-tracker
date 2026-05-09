package auth

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	natsjwt "github.com/nats-io/jwt/v2"
	natsgo "github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
)

// TokenValidator validates a raw JWT and extracts the user ID.
type TokenValidator interface {
	ValidateToken(ctx context.Context, rawToken string) (uuid.UUID, error)
}

// WorkspaceMemberChecker verifies that a user belongs to a workspace.
type WorkspaceMemberChecker interface {
	IsMember(ctx context.Context, workspaceID uuid.UUID, userID uuid.UUID) (bool, error)
}

// AuthCalloutHandler handles NATS auth callout requests on $SYS.REQ.USER.AUTH.
// Pass Handle directly to nats.Conn.Subscribe.
type AuthCalloutHandler struct {
	tokenValidator TokenValidator
	memberChecker  WorkspaceMemberChecker
	issuerKeyPair  nkeys.KeyPair
}

// NewAuthCalloutHandler creates an AuthCalloutHandler that signs responses with issuerKeyPair.
func NewAuthCalloutHandler(
	tokenValidator TokenValidator,
	memberChecker WorkspaceMemberChecker,
	issuerKeyPair nkeys.KeyPair,
) AuthCalloutHandler {
	return AuthCalloutHandler{
		tokenValidator: tokenValidator,
		memberChecker:  memberChecker,
		issuerKeyPair:  issuerKeyPair,
	}
}

// Handle processes a single NATS auth callout message.
func (h AuthCalloutHandler) Handle(message *natsgo.Msg) {
	requestClaims, err := natsjwt.DecodeAuthorizationRequestClaims(string(message.Data))
	if err != nil {
		slog.Error("decode auth request claims", "error", err)
		return
	}

	userNKey := requestClaims.UserNkey
	serverID := requestClaims.Server.ID
	workspaceIDString := requestClaims.ConnectOptions.Username
	rawToken := requestClaims.ConnectOptions.Password

	workspaceID, err := uuid.Parse(workspaceIDString)
	if err != nil {
		slog.Error("parse workspaceID", "error", err)
		h.respondWithError(message, userNKey, serverID, "invalid workspace_id")
		return
	}

	if rawToken == "" {
		slog.Error("rawToken empty")
		h.respondWithError(message, userNKey, serverID, "missing token")
		return
	}

	userID, err := h.tokenValidator.ValidateToken(context.Background(), rawToken)
	if err != nil {
		slog.Error("validateToken", "error", err)
		h.respondWithError(message, userNKey, serverID, "unauthorized")
		return
	}

	isMember, err := h.memberChecker.IsMember(context.Background(), workspaceID, userID)
	if err != nil {
		slog.Error("membership check failed", "workspace_id", workspaceID, "error", err)
		h.respondWithError(message, userNKey, serverID, "authorization service unavailable")
		return
	}
	if !isMember {
		slog.Error("not member")
		h.respondWithError(message, userNKey, serverID, "not a member of workspace")
		return
	}

	responseJWT, err := h.makeResponseJWT(workspaceID.String(), userNKey, serverID)
	if err != nil {
		slog.Error("encode auth response claims", "error", err)
		h.respondWithError(message, userNKey, serverID, "internal error")
		return
	}

	if err := message.Respond([]byte(responseJWT)); err != nil {
		slog.Error("respond to auth callout", "error", err)
	}
}

func (h AuthCalloutHandler) makeResponseJWT(workspaceID, userNKey, serverID string) (string, error) {
	scopedSubject := fmt.Sprintf("workspaces.%s.>", workspaceID)
	userClaims := natsjwt.NewUserClaims(userNKey)
	userClaims.Audience = "$G"
	userClaims.Pub.Allow = natsjwt.StringList{scopedSubject}
	userClaims.Sub.Allow = natsjwt.StringList{scopedSubject}

	userJWT, err := userClaims.Encode(h.issuerKeyPair)
	if err != nil {
		slog.Error("encode user claims", "error", err)
		return "", fmt.Errorf("encode user claims: %w", err)
	}

	responseClaims := natsjwt.NewAuthorizationResponseClaims(userNKey)
	responseClaims.Audience = serverID
	responseClaims.Jwt = userJWT

	responseJWT, err := responseClaims.Encode(h.issuerKeyPair)
	if err != nil {
		return "", fmt.Errorf("encode response claims: %w", err)
	}

	return responseJWT, nil
}

func (h AuthCalloutHandler) respondWithError(message *natsgo.Msg, userNKey, serverID, reason string) {
	responseClaims := natsjwt.NewAuthorizationResponseClaims(userNKey)
	responseClaims.Audience = serverID
	responseClaims.Error = reason

	responseJWT, err := responseClaims.Encode(h.issuerKeyPair)
	if err != nil {
		slog.Error("encode error auth response claims", "error", err)
		return
	}

	if err := message.Respond([]byte(responseJWT)); err != nil {
		slog.Error("respond with error to auth callout", "error", err)
	}
}
