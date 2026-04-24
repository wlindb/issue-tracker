package event

import (
	"context"
	"sync"
)

type Publisher[T any] func(context.Context, T)

type Event[T any] interface {
	Emit(ctx context.Context, payload T)
	AddSubscriber(handler Publisher[T]) int
}

type BaseEvent[T any] struct {
	mu         sync.RWMutex
	publishers []Publisher[T]
}

func NewBaseEvent[T any]() *BaseEvent[T] {
	return &BaseEvent[T]{
		publishers: make([]Publisher[T], 0, 1),
	}
}

func (event *BaseEvent[T]) AddPublisher(publisher Publisher[T]) int {
	event.mu.Lock()
	defer event.mu.Unlock()

	if publisher == nil {
		panic("consumer cannot be nil")
	}

	event.publishers = append(event.publishers, publisher)

	return len(event.publishers)
}

func (event *BaseEvent[T]) Emit(ctx context.Context, payload T) {
	go func() {
		done := make(chan bool)
		go event.notify(ctx, payload, done)

		select {
		case <-done:
		case <-ctx.Done():
		}
	}()
}

func (event *BaseEvent[T]) notify(ctx context.Context, payload T, done chan<- bool) {
	var wg sync.WaitGroup
	for _, publish := range event.Subscribers() {
		wg.Go(func() {
			publish(ctx, payload)
		})
	}
	wg.Wait()
	done <- true
}

func (event *BaseEvent[T]) Subscribers() []Publisher[T] {
	event.mu.RLock()
	defer event.mu.RUnlock()
	return append([]Publisher[T]{}, event.publishers...)
}
