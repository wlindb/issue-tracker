package auth_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	natsjwt "github.com/nats-io/jwt/v2"
	natsserver "github.com/nats-io/nats-server/v2/server"
	natsgo "github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/wlindb/issue-tracker/internal/application/auth"
)

// mockTokenValidator implements natsauth.TokenValidator for testing.
type mockTokenValidator struct {
	mock.Mock
}

func (m *mockTokenValidator) ValidateToken(ctx context.Context, rawToken string) (uuid.UUID, error) {
	args := m.Called(ctx, rawToken)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

// mockWorkspaceMemberChecker implements natsauth.WorkspaceMemberChecker for testing.
type mockWorkspaceMemberChecker struct {
	mock.Mock
}

func (m *mockWorkspaceMemberChecker) IsMember(ctx context.Context, workspaceID uuid.UUID, userID uuid.UUID) (bool, error) {
	args := m.Called(ctx, workspaceID, userID)
	return args.Bool(0), args.Error(1)
}

// startPlainNATSServer starts a plain embedded NATS server with no auth configuration.
func startPlainNATSServer(t *testing.T) *natsserver.Server {
	t.Helper()
	options := &natsserver.Options{
		Host:   "127.0.0.1",
		Port:   natsserver.RANDOM_PORT,
		NoLog:  true,
		NoSigs: true,
	}
	server, err := natsserver.NewServer(options)
	require.NoError(t, err)
	server.Start()
	require.True(t, server.ReadyForConnections(5*time.Second), "plain NATS server did not become ready")
	t.Cleanup(server.Shutdown)
	return server
}

func connectToServer(t *testing.T, server *natsserver.Server) *natsgo.Conn {
	t.Helper()
	connection, err := natsgo.Connect(server.ClientURL())
	require.NoError(t, err)
	t.Cleanup(connection.Close)
	return connection
}

// buildTestAuthRequest creates an AuthorizationRequestClaims JWT for testing.
func buildTestAuthRequest(t *testing.T, username, password string) (requestJWT string, userNKey string) {
	t.Helper()

	serverKeyPair, err := nkeys.CreateServer()
	require.NoError(t, err)
	serverPublicKey, err := serverKeyPair.PublicKey()
	require.NoError(t, err)

	userKeyPair, err := nkeys.CreateUser()
	require.NoError(t, err)
	userPublicKey, err := userKeyPair.PublicKey()
	require.NoError(t, err)

	claims := natsjwt.NewAuthorizationRequestClaims(serverPublicKey)
	claims.UserNkey = userPublicKey
	claims.ConnectOptions = natsjwt.ConnectOptions{
		Username: username,
		Password: password,
	}
	claims.Server = natsjwt.ServerID{ID: "test-server-id"}

	token, err := claims.Encode(serverKeyPair)
	require.NoError(t, err)

	return token, userPublicKey
}

// sendAuthRequest publishes a raw auth request to subject and waits for the response.
func sendAuthRequest(t *testing.T, connection *natsgo.Conn, subject string, requestJWT string) *natsgo.Msg {
	t.Helper()
	response, err := connection.Request(subject, []byte(requestJWT), 5*time.Second)
	require.NoError(t, err)
	return response
}

func newTestHandler(
	tokenValidator auth.TokenValidator,
	memberChecker auth.WorkspaceMemberChecker,
) (auth.AuthCalloutHandler, nkeys.KeyPair) {
	issuerKeyPair, err := nkeys.CreateAccount()
	if err != nil {
		panic(fmt.Sprintf("create issuer key pair: %v", err))
	}
	return auth.NewAuthCalloutHandler(tokenValidator, memberChecker, issuerKeyPair), issuerKeyPair
}

func Test_AuthCalloutHandler_Handle_ValidCredentials_RespondsWithSignedJWT(t *testing.T) {
	workspaceID := uuid.New()
	userID := uuid.New()
	keycloakToken := "valid.keycloak.token"

	tokenValidator := &mockTokenValidator{}
	tokenValidator.On("ValidateToken", mock.Anything, keycloakToken).Return(userID, nil)

	memberChecker := &mockWorkspaceMemberChecker{}
	memberChecker.On("IsMember", mock.Anything, workspaceID, userID).Return(true, nil)

	handler, _ := newTestHandler(tokenValidator, memberChecker)

	server := startPlainNATSServer(t)
	handlerConn := connectToServer(t, server)
	senderConn := connectToServer(t, server)

	const testSubject = "test.auth.callout"
	subscription, err := handlerConn.Subscribe(testSubject, handler.Handle)
	require.NoError(t, err)
	t.Cleanup(func() { _ = subscription.Unsubscribe() })

	requestJWT, userNKey := buildTestAuthRequest(t, workspaceID.String(), keycloakToken)
	response := sendAuthRequest(t, senderConn, testSubject, requestJWT)

	responseClaims, err := natsjwt.DecodeAuthorizationResponseClaims(string(response.Data))
	require.NoError(t, err)
	require.Empty(t, responseClaims.Error, "expected no error in response")
	require.NotEmpty(t, responseClaims.Jwt, "expected user JWT in response")
	require.Equal(t, userNKey, responseClaims.Subject)

	userClaims, err := natsjwt.DecodeUserClaims(responseClaims.Jwt)
	require.NoError(t, err)

	expectedSubject := fmt.Sprintf("workspaces.%s.>", workspaceID)
	require.Equal(t, natsjwt.StringList{expectedSubject}, userClaims.Pub.Allow)
	require.Equal(t, natsjwt.StringList{expectedSubject}, userClaims.Sub.Allow)
}

func Test_AuthCalloutHandler_Handle_InvalidWorkspaceID_RespondsWithError(t *testing.T) {
	handler, _ := newTestHandler(&mockTokenValidator{}, &mockWorkspaceMemberChecker{})

	server := startPlainNATSServer(t)
	handlerConn := connectToServer(t, server)
	senderConn := connectToServer(t, server)

	const testSubject = "test.auth.callout"
	subscription, err := handlerConn.Subscribe(testSubject, handler.Handle)
	require.NoError(t, err)
	t.Cleanup(func() { _ = subscription.Unsubscribe() })

	requestJWT, _ := buildTestAuthRequest(t, "not-a-uuid", "some.token")
	response := sendAuthRequest(t, senderConn, testSubject, requestJWT)

	responseClaims, err := natsjwt.DecodeAuthorizationResponseClaims(string(response.Data))
	require.NoError(t, err)
	require.NotEmpty(t, responseClaims.Error)
	require.Empty(t, responseClaims.Jwt)
}

func Test_AuthCalloutHandler_Handle_EmptyPassword_RespondsWithError(t *testing.T) {
	handler, _ := newTestHandler(&mockTokenValidator{}, &mockWorkspaceMemberChecker{})

	server := startPlainNATSServer(t)
	handlerConn := connectToServer(t, server)
	senderConn := connectToServer(t, server)

	const testSubject = "test.auth.callout"
	subscription, err := handlerConn.Subscribe(testSubject, handler.Handle)
	require.NoError(t, err)
	t.Cleanup(func() { _ = subscription.Unsubscribe() })

	requestJWT, _ := buildTestAuthRequest(t, uuid.New().String(), "")
	response := sendAuthRequest(t, senderConn, testSubject, requestJWT)

	responseClaims, err := natsjwt.DecodeAuthorizationResponseClaims(string(response.Data))
	require.NoError(t, err)
	require.NotEmpty(t, responseClaims.Error)
	require.Empty(t, responseClaims.Jwt)
}

func Test_AuthCalloutHandler_Handle_InvalidJWT_RespondsWithError(t *testing.T) {
	tokenValidator := &mockTokenValidator{}
	tokenValidator.On("ValidateToken", mock.Anything, "bad.token").Return(uuid.Nil, errors.New("invalid"))

	handler, _ := newTestHandler(tokenValidator, &mockWorkspaceMemberChecker{})

	server := startPlainNATSServer(t)
	handlerConn := connectToServer(t, server)
	senderConn := connectToServer(t, server)

	const testSubject = "test.auth.callout"
	subscription, err := handlerConn.Subscribe(testSubject, handler.Handle)
	require.NoError(t, err)
	t.Cleanup(func() { _ = subscription.Unsubscribe() })

	requestJWT, _ := buildTestAuthRequest(t, uuid.New().String(), "bad.token")
	response := sendAuthRequest(t, senderConn, testSubject, requestJWT)

	responseClaims, err := natsjwt.DecodeAuthorizationResponseClaims(string(response.Data))
	require.NoError(t, err)
	require.NotEmpty(t, responseClaims.Error)
	require.Empty(t, responseClaims.Jwt)
}

func Test_AuthCalloutHandler_Handle_NotMember_RespondsWithError(t *testing.T) {
	workspaceID := uuid.New()
	userID := uuid.New()
	keycloakToken := "valid.token"

	tokenValidator := &mockTokenValidator{}
	tokenValidator.On("ValidateToken", mock.Anything, keycloakToken).Return(userID, nil)

	memberChecker := &mockWorkspaceMemberChecker{}
	memberChecker.On("IsMember", mock.Anything, workspaceID, userID).Return(false, nil)

	handler, _ := newTestHandler(tokenValidator, memberChecker)

	server := startPlainNATSServer(t)
	handlerConn := connectToServer(t, server)
	senderConn := connectToServer(t, server)

	const testSubject = "test.auth.callout"
	subscription, err := handlerConn.Subscribe(testSubject, handler.Handle)
	require.NoError(t, err)
	t.Cleanup(func() { _ = subscription.Unsubscribe() })

	requestJWT, _ := buildTestAuthRequest(t, workspaceID.String(), keycloakToken)
	response := sendAuthRequest(t, senderConn, testSubject, requestJWT)

	responseClaims, err := natsjwt.DecodeAuthorizationResponseClaims(string(response.Data))
	require.NoError(t, err)
	require.NotEmpty(t, responseClaims.Error)
	require.Empty(t, responseClaims.Jwt)
}

func Test_AuthCalloutHandler_Handle_MembershipCheckError_RespondsWithError(t *testing.T) {
	workspaceID := uuid.New()
	userID := uuid.New()
	keycloakToken := "valid.token"

	tokenValidator := &mockTokenValidator{}
	tokenValidator.On("ValidateToken", mock.Anything, keycloakToken).Return(userID, nil)

	memberChecker := &mockWorkspaceMemberChecker{}
	memberChecker.On("IsMember", mock.Anything, workspaceID, userID).Return(false, errors.New("db error"))

	handler, _ := newTestHandler(tokenValidator, memberChecker)

	server := startPlainNATSServer(t)
	handlerConn := connectToServer(t, server)
	senderConn := connectToServer(t, server)

	const testSubject = "test.auth.callout"
	subscription, err := handlerConn.Subscribe(testSubject, handler.Handle)
	require.NoError(t, err)
	t.Cleanup(func() { _ = subscription.Unsubscribe() })

	requestJWT, _ := buildTestAuthRequest(t, workspaceID.String(), keycloakToken)
	response := sendAuthRequest(t, senderConn, testSubject, requestJWT)

	responseClaims, err := natsjwt.DecodeAuthorizationResponseClaims(string(response.Data))
	require.NoError(t, err)
	require.NotEmpty(t, responseClaims.Error)
	require.Empty(t, responseClaims.Jwt)
}

func Test_AuthCalloutHandler_Handle_MalformedAuthRequest_NoResponse(t *testing.T) {
	handler, _ := newTestHandler(&mockTokenValidator{}, &mockWorkspaceMemberChecker{})

	server := startPlainNATSServer(t)
	handlerConn := connectToServer(t, server)
	senderConn := connectToServer(t, server)

	const testSubject = "test.auth.callout"
	subscription, err := handlerConn.Subscribe(testSubject, handler.Handle)
	require.NoError(t, err)
	t.Cleanup(func() { _ = subscription.Unsubscribe() })

	// Publish bad data (not a valid JWT) — handler must not panic; no response expected.
	err = senderConn.Publish(testSubject, []byte("not-a-jwt"))
	require.NoError(t, err)

	// Give the handler time to process and ensure no panic.
	time.Sleep(100 * time.Millisecond)
}
