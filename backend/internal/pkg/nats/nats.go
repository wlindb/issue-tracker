package embeddednats

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	natsserver "github.com/nats-io/nats-server/v2/server"
	natsgo "github.com/nats-io/nats.go"
)

const readyTimeout = 5 * time.Second

// AuthCalloutSubject is the NATS system subject the server publishes auth callout requests to.
const AuthCalloutSubject = "$SYS.REQ.USER.AUTH"

// WorkspaceSubject builds workspace-scoped NATS subjects from a format pattern.
type WorkspaceSubject struct {
	subject string
}

// Subject formats the NATS subject string for the given workspace ID.
func (s WorkspaceSubject) Subject(workspaceID uuid.UUID) string {
	return fmt.Sprintf(s.subject, workspaceID)
}

// IssueCreatedSubject is the workspace-scoped subject pattern for issue created events.
var IssueCreatedSubject = WorkspaceSubject{subject: "workspaces.%s.issues.created"}

// IssueCreatedSubjectAll is the wildcard subject for internal consumers.
const IssueCreatedSubjectAll = "workspaces.*.issues.created"

// IssueStatusUpdatedSubject is the workspace-scoped subject pattern for issue status updated events.
var IssueStatusUpdatedSubject = WorkspaceSubject{subject: "workspaces.%s.issues.status_updated"}

// IssueStatusUpdatedSubjectAll is the wildcard subject for internal consumers.
const IssueStatusUpdatedSubjectAll = "workspaces.*.issues.status_updated"

// IssuePriorityUpdatedSubject is the workspace-scoped subject pattern for issue priority updated events.
var IssuePriorityUpdatedSubject = WorkspaceSubject{subject: "workspaces.%s.issues.priority_updated"}

// IssuePriorityUpdatedSubjectAll is the wildcard subject for internal consumers.
const IssuePriorityUpdatedSubjectAll = "workspaces.*.issues.priority_updated"

// IssueTitleUpdatedSubject is the workspace-scoped subject pattern for issue title updated events.
var IssueTitleUpdatedSubject = WorkspaceSubject{subject: "workspaces.%s.issues.title_updated"}

// IssueTitleUpdatedSubjectAll is the wildcard subject for internal consumers.
const IssueTitleUpdatedSubjectAll = "workspaces.*.issues.title_updated"

// IssueAssigneeUpdatedSubject is the workspace-scoped subject pattern for issue assignee updated events.
var IssueAssigneeUpdatedSubject = WorkspaceSubject{subject: "workspaces.%s.issues.assignee_updated"}

// IssueAssigneeUpdatedSubjectAll is the wildcard subject for internal consumers.
const IssueAssigneeUpdatedSubjectAll = "workspaces.*.issues.assignee_updated"

// IssueDescriptionUpdatedSubject is the workspace-scoped subject pattern for issue description updated events.
var IssueDescriptionUpdatedSubject = WorkspaceSubject{subject: "workspaces.%s.issues.description_updated"}

// IssueDescriptionUpdatedSubjectAll is the wildcard subject for internal consumers.
const IssueDescriptionUpdatedSubjectAll = "workspaces.*.issues.description_updated"

// CommentCreatedSubject is the workspace-and-issue-scoped subject pattern for comment created events.
var CommentCreatedSubject = IssueCommentSubject{subject: "workspaces.%s.issues.%s.comments.created"}

// CommentCreatedSubjectAll is the wildcard subject for internal consumers.
const CommentCreatedSubjectAll = "workspaces.*.issues.*.comments.created"

// IssueCommentSubject builds workspace-and-issue-scoped NATS subjects from a format pattern.
type IssueCommentSubject struct {
	subject string
}

// Subject formats the NATS subject string for the given workspace and issue IDs.
func (s IssueCommentSubject) Subject(workspaceID, issueID uuid.UUID) string {
	return fmt.Sprintf(s.subject, workspaceID, issueID)
}

// ProjectCreatedSubject is the workspace-scoped subject pattern for project created events.
var ProjectCreatedSubject = WorkspaceSubject{subject: "workspaces.%s.projects.created"}

// ProjectCreatedSubjectAll is the wildcard subject for internal consumers.
const ProjectCreatedSubjectAll = "workspaces.*.projects.created"

// ServerOption is a functional option applied to the embedded NATS server configuration.
type ServerOption func(*natsserver.Options) error

// AuthCalloutConfig holds configuration for NATS Auth Callout.
type AuthCalloutConfig struct {
	// IssuerPublicKey is the account NKey public key used to sign authorization responses.
	IssuerPublicKey string
	// AuthUsers is the list of usernames that bypass the auth callout (internal service connections).
	AuthUsers []string
}

// WithAuthCallout configures the embedded server to use NATS Auth Callout.
// All connecting clients whose username is not in config.AuthUsers will be forwarded
// to the auth callout handler on AuthCalloutSubject.
func WithAuthCallout(config AuthCalloutConfig) ServerOption {
	return func(options *natsserver.Options) error {
		options.AuthCallout = &natsserver.AuthCallout{
			Issuer:    config.IssuerPublicKey,
			AuthUsers: config.AuthUsers,
		}
		return nil
	}
}

// WithExternalPort binds the server to all interfaces on the given port,
// making it reachable by external clients. Without this option the server
// binds to loopback only on a random port.
func WithExternalPort(port int) ServerOption {
	return func(options *natsserver.Options) error {
		options.Host = "0.0.0.0"
		options.Port = port
		return nil
	}
}

// WithWebSocketPort enables the WebSocket listener on all interfaces on the
// given port. NoTLS is set to true so the listener accepts plain ws:// connections
// (suitable for development or when TLS is terminated upstream).
func WithWebSocketPort(port int) ServerOption {
	return func(options *natsserver.Options) error {
		options.Websocket.Host = "0.0.0.0"
		options.Websocket.Port = port
		options.Websocket.NoTLS = true
		return nil
	}
}

// WithInternalUser adds a named user to the server that can connect without
// going through the auth callout. Use this together with WithAuthCallout by
// listing the same username in AuthCalloutConfig.AuthUsers.
func WithInternalUser(username, password string) ServerOption {
	return func(options *natsserver.Options) error {
		options.Users = append(options.Users, &natsserver.User{
			Username: username,
			Password: password,
		})
		return nil
	}
}

// StartEmbeddedServer starts an in-process NATS server and returns it once it is
// ready to accept connections. The caller is responsible for calling server.Shutdown().
func StartEmbeddedServer(opts ...ServerOption) (*natsserver.Server, error) {
	options := &natsserver.Options{
		Host:   "127.0.0.1",
		Port:   natsserver.RANDOM_PORT,
		NoLog:  true,
		NoSigs: true,
	}
	for _, option := range opts {
		if err := option(options); err != nil {
			return nil, fmt.Errorf("apply server option: %w", err)
		}
	}
	server, err := natsserver.NewServer(options)
	if err != nil {
		return nil, fmt.Errorf("new embedded nats server: %w", err)
	}
	server.Start()
	if !server.ReadyForConnections(readyTimeout) {
		return nil, fmt.Errorf("embedded nats server did not become ready within %s", readyTimeout)
	}
	return server, nil
}

// Connect returns a NATS client connection to the given embedded server.
// Additional nats.go client options (e.g. natsgo.UserInfo) may be passed via opts.
// The caller is responsible for calling connection.Close() on shutdown.
func Connect(server *natsserver.Server, opts ...natsgo.Option) (*natsgo.Conn, error) {
	connection, err := natsgo.Connect(server.ClientURL(), opts...)
	if err != nil {
		return nil, fmt.Errorf("connect to embedded nats: %w", err)
	}
	return connection, nil
}
