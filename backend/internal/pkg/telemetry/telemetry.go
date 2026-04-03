// Package telemetry initialises OpenTelemetry (traces, metrics, logs) and
// exports via OTLP/HTTP—ready for Grafana Cloud or any OTLP-compatible backend.
package telemetry

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	otellog "go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// Config holds the values needed to configure OpenTelemetry exporters.
// Endpoint and headers are read directly from the standard OTel environment
// variables (OTEL_EXPORTER_OTLP_ENDPOINT, OTEL_EXPORTER_OTLP_HEADERS) by
// the SDK, so they are not included here.
type Config struct {
	ServiceName string
}

// Setup initialises the OTel SDK and registers global providers.
// It returns a shutdown function that flushes and closes all exporters.
// The caller should defer shutdown(ctx) during graceful server stop.
func Setup(ctx context.Context, cfg Config) (func(context.Context) error, error) {
	var zero func(context.Context) error

	if cfg.ServiceName == "" {
		return zero, fmt.Errorf("telemetry: service name is required")
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(semconv.ServiceName(cfg.ServiceName)),
		resource.WithTelemetrySDK(),
	)
	if err != nil {
		return zero, fmt.Errorf("telemetry resource: %w", err)
	}

	traceShutdown, err := setupTraces(ctx, res)
	if err != nil {
		return zero, err
	}

	metricShutdown, err := setupMetrics(ctx, res)
	if err != nil {
		return zero, err
	}

	logShutdown, err := setupLogs(ctx, res)
	if err != nil {
		return zero, err
	}

	shutdownFuncs := []func(context.Context) error{traceShutdown, metricShutdown, logShutdown}
	shutdownAll := func(ctx context.Context) error {
		var errs []error
		for _, fn := range shutdownFuncs {
			errs = append(errs, fn(ctx))
		}
		return errors.Join(errs...)
	}
	return shutdownAll, nil
}

func setupTraces(ctx context.Context, res *resource.Resource) (func(context.Context) error, error) {
	exporter, err := otlptracehttp.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("trace exporter: %w", err)
	}
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter, sdktrace.WithBatchTimeout(5*time.Second)),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	return provider.Shutdown, nil
}

func setupMetrics(ctx context.Context, res *resource.Resource) (func(context.Context) error, error) {
	exporter, err := otlpmetrichttp.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("metric exporter: %w", err)
	}
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(30*time.Second))),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(provider)
	return provider.Shutdown, nil
}

func setupLogs(ctx context.Context, res *resource.Resource) (func(context.Context) error, error) {
	exporter, err := otlploghttp.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("log exporter: %w", err)
	}
	provider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exporter, sdklog.WithExportTimeout(5*time.Second))),
		sdklog.WithResource(res),
	)
	otellog.SetLoggerProvider(provider)
	return provider.Shutdown, nil
}
