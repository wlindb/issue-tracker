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
	cfg := telemetry.Config{
		OTLPEndpoint: "http://127.0.0.1:4318",
	}

	shutdown, err := telemetry.Setup(context.Background(), cfg)
	require.Error(t, err)
	assert.Nil(t, shutdown)
	assert.Contains(t, err.Error(), "service name is required")
}

func Test_Setup_EmptyEndpoint_ReturnsError(t *testing.T) {
	cfg := telemetry.Config{
		ServiceName: "test-service",
	}

	shutdown, err := telemetry.Setup(context.Background(), cfg)
	require.Error(t, err)
	assert.Nil(t, shutdown)
	assert.Contains(t, err.Error(), "OTLP endpoint is required")
}

func Test_Setup_EmptyConfig_ReturnsError(t *testing.T) {
	shutdown, err := telemetry.Setup(context.Background(), telemetry.Config{})
	require.Error(t, err)
	assert.Nil(t, shutdown)
}

func Test_Setup_ValidConfig_ReturnsShutdownFunc(t *testing.T) {
	cfg := telemetry.Config{
		ServiceName:  "test-service",
		OTLPEndpoint: "http://127.0.0.1:4318",
	}

	shutdown, err := telemetry.Setup(context.Background(), cfg)
	require.NoError(t, err)
	require.NotNil(t, shutdown)

	// Shutdown may return an error when the collector is unreachable
	// (batched data cannot be flushed). We only verify that Setup succeeds.
	_ = shutdown(context.Background())
}

func Test_Setup_WithHeaders_ReturnsShutdownFunc(t *testing.T) {
	cfg := telemetry.Config{
		ServiceName:  "test-service",
		OTLPEndpoint: "http://127.0.0.1:4318",
		OTLPHeaders:  map[string]string{"Authorization": "Basic dGVzdDp0ZXN0"},
	}

	shutdown, err := telemetry.Setup(context.Background(), cfg)
	require.NoError(t, err)
	require.NotNil(t, shutdown)

	_ = shutdown(context.Background())
}
