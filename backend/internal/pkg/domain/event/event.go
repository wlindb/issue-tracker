package event

import (
	"context"
	"sync"
)

type Subscriber[T any] func(context.Context, T)

type Event[T any] interface {
	Emit(ctx context.Context, payload T)
	AddSubscriber(handler Subscriber[T]) int
}

type BaseEvent[T any] struct {
	mu          sync.RWMutex
	subscribers []Subscriber[T]
}

func NewBaseEvent[T any]() *BaseEvent[T] {
	return &BaseEvent[T]{
		subscribers: make([]Subscriber[T], 0, 1),
	}
}

func (event *BaseEvent[T]) AddSubscriber(subscriber Subscriber[T]) int {
	event.mu.Lock()
	defer event.mu.Unlock()

	if subscriber == nil {
		panic("consumer cannot be nil")
	}

	event.subscribers = append(event.subscribers, subscriber)

	return len(event.subscribers)
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
	for _, subscriber := range event.Subscribers() {
		wg.Go(func() {
			subscriber(ctx, payload)
		})
	}
	wg.Wait()
	done <- true
}

func (event *BaseEvent[T]) Subscribers() []Subscriber[T] {
	event.mu.RLock()
	defer event.mu.RUnlock()
	return append([]Subscriber[T]{}, event.subscribers...)
}
