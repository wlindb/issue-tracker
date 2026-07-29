package embeddednats_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	natsgo "github.com/nats-io/nats.go"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	embeddednats "github.com/wlindb/issue-tracker/internal/pkg/nats"
)

func Test_WorkspaceSubject_Subject_ValidWorkspaceID_ReturnsWorkspaceScopedSubject(t *testing.T) {
	workspaceID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	expected := fmt.Sprintf("workspaces.%s.issues.created", workspaceID)

	actual := embeddednats.IssueCreatedSubject.Subject(workspaceID)

	assert.Equal(t, expected, actual)
}

func Test_WorkspaceSubject_WorkspaceID_ValidSubject_ReturnsWorkspaceID(t *testing.T) {
	expected := uuid.MustParse("00000000-0000-0000-0000-000000000003")
	subject := embeddednats.IssueCreatedSubject.Subject(expected)

	actual, err := embeddednats.IssueCreatedSubject.WorkspaceID(subject)

	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func Test_WorkspaceSubject_WorkspaceID_MismatchedSubject_ReturnsError(t *testing.T) {
	_, err := embeddednats.IssueCreatedSubject.WorkspaceID("workspaces.not-a-uuid.issues.created.extra")

	require.Error(t, err)
}

func Test_WorkspaceSubject_All_ReturnsWildcardSubject(t *testing.T) {
	expected := "workspaces.*.issues.created"

	actual := embeddednats.IssueCreatedSubject.All()

	assert.Equal(t, expected, actual)
}

func Test_IssueCommentSubject_WorkspaceID_ValidSubject_ReturnsWorkspaceID(t *testing.T) {
	expected := uuid.MustParse("00000000-0000-0000-0000-000000000004")
	issueID := uuid.MustParse("00000000-0000-0000-0000-000000000005")
	subject := embeddednats.CommentCreatedSubject.Subject(expected, issueID)

	actual, err := embeddednats.CommentCreatedSubject.WorkspaceID(subject)

	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func Test_IssueCommentSubject_WorkspaceID_MismatchedSubject_ReturnsError(t *testing.T) {
	_, err := embeddednats.CommentCreatedSubject.WorkspaceID("workspaces.not-a-uuid.issues.created")

	require.Error(t, err)
}

func Test_NATSEventPublisher_Publish_WithSubjectAndPayload_PublishesMessage(t *testing.T) {
	server, err := embeddednats.StartEmbeddedServer()
	require.NoError(t, err)
	t.Cleanup(server.Shutdown)

	connection, err := embeddednats.Connect(server)
	require.NoError(t, err)
	t.Cleanup(connection.Close)

	workspaceID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	expectedSubject := embeddednats.IssueCreatedSubject.Subject(workspaceID)

	received := make(chan string, 1)
	_, err = connection.Subscribe(expectedSubject, func(message *natsgo.Msg) {
		received <- message.Subject
	})
	require.NoError(t, err)

	payload, err := json.Marshal(testEvent{Name: "test"})
	require.NoError(t, err)

	publisher := embeddednats.NewNATSEventPublisher(connection)

	err = publisher.Publish(context.Background(), expectedSubject, payload)
	require.NoError(t, err)

	select {
	case actualSubject := <-received:
		assert.Equal(t, expectedSubject, actualSubject)
	case <-time.After(time.Second):
		t.Fatal("did not receive NATS message within timeout")
	}
}

type testEvent struct {
	Name string `json:"name"`
}
