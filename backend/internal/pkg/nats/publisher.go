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

type NATSEventPublisher[T any] struct {
	connection *natsgo.Conn
	Publisher  event.Publisher[T]
}

type SubjectResolver[T any] interface {
	Resolve(ctx context.Context, event T) (string, error)
}

func NewNATSEventPublisher[T any](connection *natsgo.Conn, subjectResolver SubjectResolver[T]) *NATSEventPublisher[T] {
	return &NATSEventPublisher[T]{connection: connection, Publisher: publish(connection, subjectResolver)}
}

func publish[T any](connection *natsgo.Conn, subjectResolver SubjectResolver[T]) event.Publisher[T] {
	return func(ctx context.Context, event T) error {
		tracer := otel.Tracer("nats-publisher")
		ctx, span := tracer.Start(ctx, "nats-publisher", trace.WithSpanKind(trace.SpanKindProducer))
		defer span.End()

		subject, err := subjectResolver.Resolve(ctx, event)
		if err != nil {
			err := fmt.Errorf("publish: workspace ID missing from context")
			span.RecordError(err)
			span.SetStatus(codes.Error, "workspace ID missing from context")
			return err
		}

		return publishMsg(ctx, span, connection, subject, event)
	}
}

func publishMsg[T any](ctx context.Context, span trace.Span, connection *natsgo.Conn, subject string, evt T) error {
	payload, err := json.Marshal(evt)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to marshal event")
		return fmt.Errorf("publish to %s marshal: %w", subject, err)
	}

	msg := &natsgo.Msg{Subject: subject, Data: payload, Header: natsgo.Header{}}
	otel.GetTextMapPropagator().Inject(ctx, natsHeaderCarrier(msg.Header))
	if err := connection.PublishMsg(msg); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "publish failed")
		return fmt.Errorf("publish to %s publish: %w", subject, err)
	}

	return nil
}
