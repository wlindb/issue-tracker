package embeddednats

import (
	"context"
	"encoding/json"

	natsgo "github.com/nats-io/nats.go"

	"github.com/wlindb/issue-tracker/internal/pkg/domain/event"
)

type NATSEventPublisher[T any] struct {
	connection *natsgo.Conn
	Publisher  event.Publisher[T]
}

func NewNATSEventPublisher[T any](connection *natsgo.Conn, subject string) *NATSEventPublisher[T] {
	return &NATSEventPublisher[T]{connection: connection, Publisher: publish[T](connection, subject)}
}

func publish[T any](connection *natsgo.Conn, subject string) event.Publisher[T] {
	return func(_ context.Context, event T) {
		payload, err := json.Marshal(event)
		if err != nil {
			return
		}
		if err := connection.Publish(subject, payload); err != nil {
			return
		}
	}
}
