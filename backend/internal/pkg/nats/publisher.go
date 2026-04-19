package embeddednats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	natsgo "github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	key "github.com/wlindb/issue-tracker/internal/pkg/context"
	"github.com/wlindb/issue-tracker/internal/pkg/domain/event"
)

type NATSEventPublisher[T any] struct {
	connection *natsgo.Conn
	Publisher  event.Publisher[T]
}

func NewNATSEventPublisher[T any](connection *natsgo.Conn, subject WorkspaceSubject) *NATSEventPublisher[T] {
	return &NATSEventPublisher[T]{connection: connection, Publisher: publish[T](connection, subject)}
}

func publish[T any](connection *natsgo.Conn, workspaceSubject WorkspaceSubject) event.Publisher[T] {
	return func(ctx context.Context, evt T) error {
		tracer := otel.Tracer("nats-publisher")
		ctx, span := tracer.Start(ctx, workspaceSubject.subject, trace.WithSpanKind(trace.SpanKindProducer))
		defer span.End()

		workspaceID, ok := ctx.Value(key.WorkspaceID).(uuid.UUID)
		if !ok {
			err := fmt.Errorf("publish: workspace ID missing from context")
			span.RecordError(err)
			span.SetStatus(codes.Error, "workspace ID missing from context")
			return err
		}
		subject := workspaceSubject.Subject(workspaceID)

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
}
