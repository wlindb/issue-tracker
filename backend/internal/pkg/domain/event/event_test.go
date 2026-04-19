package event_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wlindb/issue-tracker/internal/pkg/domain/event"
)

func Test_Publish_ContextOverride_CallsOverridePublisher(t *testing.T) {
	t.Parallel()

	e := event.New[string]()
	var received string
	ctx := event.WithPublisher[string](context.Background(), func(_ context.Context, payload string) error {
		received = payload
		return nil
	})

	err := e.Publish(ctx, "hello")

	require.NoError(t, err)
	assert.Equal(t, "hello", received)
}

func Test_Publish_ContextOverride_IgnoresSingleton(t *testing.T) {
	t.Parallel()

	e := event.New[string]()
	singletonCalled := false
	require.NoError(t, e.AddPublisher(func(_ context.Context, _ string) error {
		singletonCalled = true
		return nil
	}))

	ctx := event.WithPublisher[string](context.Background(), func(_ context.Context, _ string) error {
		return nil
	})

	err := e.Publish(ctx, "hello")

	require.NoError(t, err)
	assert.False(t, singletonCalled)
}

func Test_Publish_NoContextOverride_CallsSingleton(t *testing.T) {
	t.Parallel()

	e := event.New[string]()
	var received string
	require.NoError(t, e.AddPublisher(func(_ context.Context, payload string) error {
		received = payload
		return nil
	}))

	err := e.Publish(context.Background(), "world")

	require.NoError(t, err)
	assert.Equal(t, "world", received)
}

func Test_Publish_NeitherPublisherNorContextOverride_ReturnsError(t *testing.T) {
	t.Parallel()

	e := event.New[string]()

	err := e.Publish(context.Background(), "orphan")

	require.Error(t, err)
}

func Test_Publish_ContextOverride_IsolatedAcrossParallelTests(t *testing.T) {
	t.Parallel()

	e := event.New[string]()

	t.Run("first", func(t *testing.T) {
		t.Parallel()
		var received string
		ctx := event.WithPublisher[string](context.Background(), func(_ context.Context, payload string) error {
			received = payload
			return nil
		})
		require.NoError(t, e.Publish(ctx, "first"))
		assert.Equal(t, "first", received)
	})

	t.Run("second", func(t *testing.T) {
		t.Parallel()
		var received string
		ctx := event.WithPublisher[string](context.Background(), func(_ context.Context, payload string) error {
			received = payload
			return nil
		})
		require.NoError(t, e.Publish(ctx, "second"))
		assert.Equal(t, "second", received)
	})
}

func Test_Publish_ContextOverride_PropagatesPublisherError(t *testing.T) {
	t.Parallel()

	e := event.New[string]()
	publisherErr := errors.New("downstream failure")
	ctx := event.WithPublisher[string](context.Background(), func(_ context.Context, _ string) error {
		return publisherErr
	})

	err := e.Publish(ctx, "fail")

	require.Error(t, err)
	assert.ErrorIs(t, err, publisherErr)
}
