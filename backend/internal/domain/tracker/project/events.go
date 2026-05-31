package project

import (
	"time"

	"github.com/wlindb/issue-tracker/internal/pkg/domain/event"
)

type ProjectCreatedEvent struct {
	OccurredAt time.Time `json:"occurred_at"`
	Payload    Project   `json:"payload"`
}

var Created = event.New[ProjectCreatedEvent]()
