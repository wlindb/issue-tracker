package comment

import (
	"time"

	"github.com/wlindb/issue-tracker/internal/pkg/domain/event"
)

type CommentCreatedEvent struct {
	OccurredAt time.Time `json:"occurred_at"`
	Payload    Comment
}

var Created = event.New[CommentCreatedEvent]()
