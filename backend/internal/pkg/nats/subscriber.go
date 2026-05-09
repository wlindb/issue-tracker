package embeddednats

import (
	"context"
	"encoding/json"
	"fmt"

	natsgo "github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

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
	tracer := otel.Tracer("nats-subscriber")

	_, err := s.connection.Subscribe(s.subject, func(message *natsgo.Msg) {
		ctx := otel.GetTextMapPropagator().Extract(context.Background(), natsHeaderCarrier(message.Header))
		ctx, span := tracer.Start(ctx, s.subject, trace.WithSpanKind(trace.SpanKindConsumer))
		defer span.End()

		var event T
		if err := json.Unmarshal(message.Data, &event); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "failed to unmarshal")
			return
		}

		if err := handler(ctx, event); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "handler error")
			return
		}
	})
	if err != nil {
		return fmt.Errorf("subscribe subject %s: %w", s.subject, err)
	}

	return nil
}
