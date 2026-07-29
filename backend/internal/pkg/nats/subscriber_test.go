//go:build !integration

package embeddednats_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	key "github.com/wlindb/issue-tracker/internal/pkg/context"
	embeddednats "github.com/wlindb/issue-tracker/internal/pkg/nats"
)

func Test_NATSEventSubscriber_Subscribe_ValidMessage_InjectsWorkspaceIDIntoContext(t *testing.T) {
	server, err := embeddednats.StartEmbeddedServer()
	require.NoError(t, err)
	t.Cleanup(server.Shutdown)

	connection, err := embeddednats.Connect(server)
	require.NoError(t, err)
	t.Cleanup(connection.Close)

	expectedWorkspaceID := uuid.MustParse("00000000-0000-0000-0000-000000000006")
	subject := embeddednats.IssueCreatedSubject.Subject(expectedWorkspaceID)

	subscriber := embeddednats.NewNATSEventSubscriber[testEvent](connection, embeddednats.IssueCreatedSubject)

	received := make(chan context.Context, 1)
	err = subscriber.Subscribe(func(ctx context.Context, _ testEvent) error {
		received <- ctx
		return nil
	})
	require.NoError(t, err)

	payload, err := json.Marshal(testEvent{Name: "test"})
	require.NoError(t, err)

	publisher := embeddednats.NewNATSEventPublisher(connection)
	err = publisher.Publish(context.Background(), subject, payload)
	require.NoError(t, err)

	select {
	case ctx := <-received:
		actual, ok := ctx.Value(key.WorkspaceID).(uuid.UUID)
		require.True(t, ok)
		require.Equal(t, expectedWorkspaceID, actual)
	case <-time.After(time.Second):
		t.Fatal("did not receive NATS message within timeout")
	}
}

func Test_NATSEventSubscriber_Subscribe_UnparsableSubject_DoesNotInvokeHandler(t *testing.T) {
	server, err := embeddednats.StartEmbeddedServer()
	require.NoError(t, err)
	t.Cleanup(server.Shutdown)

	connection, err := embeddednats.Connect(server)
	require.NoError(t, err)
	t.Cleanup(connection.Close)

	subject := "workspaces.not-a-uuid.issues.created"

	subscriber := embeddednats.NewNATSEventSubscriber[testEvent](connection, embeddednats.IssueCreatedSubject)

	called := make(chan struct{}, 1)
	err = subscriber.Subscribe(func(_ context.Context, _ testEvent) error {
		called <- struct{}{}
		return nil
	})
	require.NoError(t, err)

	payload, err := json.Marshal(testEvent{Name: "test"})
	require.NoError(t, err)

	publisher := embeddednats.NewNATSEventPublisher(connection)
	err = publisher.Publish(context.Background(), subject, payload)
	require.NoError(t, err)

	select {
	case <-called:
		t.Fatal("handler was invoked despite unparsable workspace ID in subject")
	case <-time.After(200 * time.Millisecond):
	}
}
