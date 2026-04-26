package embeddednats

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	natsgo "github.com/nats-io/nats.go"

	"github.com/wlindb/issue-tracker/internal/pkg/domain/event"
)

type NATSEventSubscriber[T any] struct {
	connection *natsgo.Conn
	subject    string
}

func NewNATSEventSubscriber[T any](connection *natsgo.Conn, subject string) NATSEventSubscriber[T] {
	return NATSEventSubscriber[T]{connection: connection, subject: subject}
}

func (s NATSEventSubscriber[T]) Subscribe(handler event.Subscriber[T]) error {
	_, err := s.connection.Subscribe(s.subject, func(message *natsgo.Msg) {
		var event T
		if err := json.Unmarshal(message.Data, &event); err != nil {
			slog.Error("unmarshal issue created event", "error", err)
			return
		}
		// TODO: Fix context
		ctx := context.Background()
		handler(ctx, event)
	})
	if err != nil {
		return fmt.Errorf("subscribe subject %s: %w", s.subject, err)
	}

	return nil
}
