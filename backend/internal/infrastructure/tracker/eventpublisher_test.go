//go:build !integration

package tracker

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	issuedomain "github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
)

func Test_IssueCreatedSubjectResolver_Resolve_MissingWorkspaceID_ReturnsError(t *testing.T) {
	resolver := IssueCreatedSubjectResolver{}
	event := issuedomain.IssueCreatedEvent{OccurredAt: time.Now().UTC()}

	_, err := resolver.Resolve(context.Background(), event)

	assert.ErrorContains(t, err, "workspace ID missing from context")
}
