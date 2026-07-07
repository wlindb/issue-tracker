package embeddednats

import (
	"context"
	"fmt"

	natsgo "github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
)

type NATSEventPublisher struct {
	connection *natsgo.Conn
}

func NewNATSEventPublisher(connection *natsgo.Conn) *NATSEventPublisher {
	return &NATSEventPublisher{connection: connection}
}

func (publisher NATSEventPublisher) Publish(ctx context.Context, subject string, payload []byte) error {
	msg := &natsgo.Msg{Subject: subject, Data: payload, Header: natsgo.Header{}}
	otel.GetTextMapPropagator().Inject(ctx, natsHeaderCarrier(msg.Header))
	if err := publisher.connection.PublishMsg(msg); err != nil {
		return fmt.Errorf("publish to %s publish: %w", subject, err)
	}

	return nil
}
