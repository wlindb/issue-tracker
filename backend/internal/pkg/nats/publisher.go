package embeddednats

import (
	"context"
	"encoding/json"
	"fmt"

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
	return func(_ context.Context, event T) error {
		payload, err := json.Marshal(event)
		if err != nil {
			return fmt.Errorf("publish to %s marshal: %w", subject, err)
		}
		if err := connection.Publish(subject, payload); err != nil {
			return fmt.Errorf("publish to %s publish: %w", subject, err)
		}
		return nil
	}
}
