package embeddednats_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	natsgo "github.com/nats-io/nats.go"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	key "github.com/wlindb/issue-tracker/internal/pkg/context"
	embeddednats "github.com/wlindb/issue-tracker/internal/pkg/nats"
)

func Test_WorkspaceSubject_Subject_ValidWorkspaceID_ReturnsWorkspaceScopedSubject(t *testing.T) {
	workspaceID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	expected := fmt.Sprintf("workspaces.%s.issues.created", workspaceID)

	actual := embeddednats.IssueCreatedSubject.Subject(workspaceID)

	assert.Equal(t, expected, actual)
}

func Test_NATSEventPublisher_Publish_WorkspaceIDInContext_PublishesToWorkspaceScopedSubject(t *testing.T) {
	server, err := embeddednats.StartEmbeddedServer()
	require.NoError(t, err)
	t.Cleanup(server.Shutdown)

	connection, err := embeddednats.Connect(server)
	require.NoError(t, err)
	t.Cleanup(connection.Close)

	workspaceID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	expectedSubject := fmt.Sprintf("workspaces.%s.issues.created", workspaceID)

	received := make(chan string, 1)
	_, err = connection.Subscribe(expectedSubject, func(message *natsgo.Msg) {
		received <- message.Subject
	})
	require.NoError(t, err)

	publisher := embeddednats.NewNATSEventPublisher(connection, workspaceReturningSubjectResolver{})
	ctx := context.WithValue(context.Background(), key.WorkspaceID, workspaceID)

	err = publisher.Publisher(ctx, testEvent{Name: "test"})
	require.NoError(t, err)

	select {
	case actualSubject := <-received:
		assert.Equal(t, expectedSubject, actualSubject)
	case <-time.After(time.Second):
		t.Fatal("did not receive NATS message within timeout")
	}
}

func Test_NATSEventPublisher_Publish_MissingWorkspaceID_ReturnsError(t *testing.T) {
	server, err := embeddednats.StartEmbeddedServer()
	require.NoError(t, err)
	t.Cleanup(server.Shutdown)

	connection, err := embeddednats.Connect(server)
	require.NoError(t, err)
	t.Cleanup(connection.Close)

	publisher := embeddednats.NewNATSEventPublisher(connection, workspaceMissingSubjectResolver{})

	err = publisher.Publisher(context.Background(), testEvent{Name: "test"})

	assert.ErrorContains(t, err, "workspace ID missing from context")
}

type testEvent struct {
	Name string `json:"name"`
}

func workspaceID(ctx context.Context) (uuid.UUID, error) {
	workspaceID, ok := ctx.Value(key.WorkspaceID).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("workspace ID missing from context")
	}
	return workspaceID, nil
}

type workspaceReturningSubjectResolver struct{}

func (workspaceReturningSubjectResolver) Resolve(ctx context.Context, _ testEvent) (string, error) {
	workspaceID, err := workspaceID(ctx)
	return embeddednats.IssueCreatedSubject.Subject(workspaceID), err
}

type workspaceMissingSubjectResolver struct{}

func (workspaceMissingSubjectResolver) Resolve(_ context.Context, _ testEvent) (string, error) {
	return "", errors.New("mocked workspace ID missing from context")
}
