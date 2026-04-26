package event

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type (
	Publisher[T any]  func(context.Context, T) error
	Subscriber[T any] func(context.Context, T)
)

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

func (event *Event[T]) Publish(ctx context.Context, payload T) error {
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
