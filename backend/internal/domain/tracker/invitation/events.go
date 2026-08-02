package invitation

import (
	"context"
	"fmt"
	"time"

	"github.com/wlindb/issue-tracker/internal/pkg/domain/event"
)

type InvitationCreatedEvent struct {
	OccurredAt time.Time
	Payload    Invitation
}

var Created = event.New[InvitationCreatedEvent]()

// EmitCreated publishes an InvitationCreatedEvent for this invitation.
func (i Invitation) EmitCreated(ctx context.Context) error {
	event := InvitationCreatedEvent{OccurredAt: time.Now().UTC(), Payload: i}
	if err := Created.Publish(ctx, event); err != nil {
		return fmt.Errorf("invitation emit created: %w", err)
	}
	return nil
}
