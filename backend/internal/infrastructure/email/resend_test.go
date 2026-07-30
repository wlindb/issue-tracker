package email

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewResendSender_ValidAPIKeyAndFrom_StoresFromAddress(t *testing.T) {
	sender := NewResendSender("fake-api-key", "Acme <invites@acme.com>")

	require.NotNil(t, sender)
	assert.Equal(t, "Acme <invites@acme.com>", sender.from)
}
