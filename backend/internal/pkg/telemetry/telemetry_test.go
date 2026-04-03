//go:build !integration

package telemetry_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wlindb/issue-tracker/internal/pkg/telemetry"
)

func Test_Setup_EmptyServiceName_ReturnsError(t *testing.T) {
	shutdown, err := telemetry.Setup(context.Background(), telemetry.Config{})
	require.Error(t, err)
	assert.Nil(t, shutdown)
	assert.Contains(t, err.Error(), "service name is required")
}

func Test_Setup_ValidConfig_ReturnsShutdownFunc(t *testing.T) {
	cfg := telemetry.Config{
		ServiceName: "test-service",
	}

	shutdown, err := telemetry.Setup(context.Background(), cfg)
	require.NoError(t, err)
	require.NotNil(t, shutdown)

	// Shutdown may return an error when the collector is unreachable
	// (batched data cannot be flushed). We only verify that Setup succeeds.
	_ = shutdown(context.Background())
}
