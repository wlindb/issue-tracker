package event

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type (
	Publisher[T any]    func(context.Context, T) error
	Subscriber[T any]   func(context.Context, T) error
	publisherKey[T any] struct{}
)

type SubscriberOf[T any] interface {
	Subscribe(handler Subscriber[T]) error
}

type Event[T any] struct {
	mu        sync.RWMutex
	publisher Publisher[T]
}

func New[T any]() *Event[T] {
	return &Event[T]{}
}

func (event *Event[T]) AddPublisher(publisher Publisher[T]) error {
	event.mu.Lock()
	defer event.mu.Unlock()

	if publisher == nil {
		return errors.New("publisher cannot be nil")
	}

	if event.publisher != nil {
		return errors.New("publisher can only be assigned once")
	}

	event.publisher = publisher
	return nil
}

// WithPublisher returns a context carrying a publisher override. When Publish
// is called with this context, the override is used instead of the singleton.
func WithPublisher[T any](ctx context.Context, publisher Publisher[T]) context.Context {
	return context.WithValue(ctx, publisherKey[T]{}, publisher)
}

func (event *Event[T]) Publish(ctx context.Context, payload T) error {
	if override, ok := ctx.Value(publisherKey[T]{}).(Publisher[T]); ok {
		return override(ctx, payload)
	}

	event.mu.RLock()
	defer event.mu.RUnlock()
	if event.publisher == nil {
		return errors.New("publisher is missing")
	}

	if err := event.publisher(ctx, payload); err != nil {
		return fmt.Errorf("publish event: %w", err)
	}

	return nil
}
