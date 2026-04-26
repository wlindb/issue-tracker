package event

import (
	"context"
	"sync"
)

type (
	Publisher[T any]  func(context.Context, T)
	Subscriber[T any] func(context.Context, T)
)

// type Event[T any] interface {
// 	Emit(ctx context.Context, payload T)
// 	AddPublisher(handler Publisher[T]) int
// }

type Event[T any] struct {
	mu         sync.RWMutex
	publishers []Publisher[T]
}

func New[T any]() *Event[T] {
	return &Event[T]{
		publishers: make([]Publisher[T], 0, 1),
	}
}

func (event *Event[T]) AddPublisher(publisher Publisher[T]) int {
	event.mu.Lock()
	defer event.mu.Unlock()

	if publisher == nil {
		panic("consumer cannot be nil")
	}

	event.publishers = append(event.publishers, publisher)

	return len(event.publishers)
}

func (event *Event[T]) Publish(ctx context.Context, payload T) {
	go func() {
		done := make(chan bool)
		go event.notify(ctx, payload, done)

		select {
		case <-done:
		case <-ctx.Done():
		}
	}()
}

func (event *Event[T]) notify(ctx context.Context, payload T, done chan<- bool) {
	var wg sync.WaitGroup
	for _, publish := range event.Subscribers() {
		wg.Go(func() {
			publish(ctx, payload)
		})
	}
	wg.Wait()
	done <- true
}

func (event *Event[T]) Subscribers() []Publisher[T] {
	event.mu.RLock()
	defer event.mu.RUnlock()
	return append([]Publisher[T]{}, event.publishers...)
}
