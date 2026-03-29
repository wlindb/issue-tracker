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
type Config struct {
	ServiceName  string
	OTLPEndpoint string
	OTLPHeaders  map[string]string
}

// Setup initialises the OTel SDK and registers global providers.
// It returns a shutdown function that flushes and closes all exporters.
// The caller should defer shutdown(ctx) during graceful server stop.
func Setup(ctx context.Context, cfg Config) (func(context.Context) error, error) {
	if cfg.ServiceName == "" {
		return nil, fmt.Errorf("telemetry: service name is required")
	}
	if cfg.OTLPEndpoint == "" {
		return nil, fmt.Errorf("telemetry: OTLP endpoint is required")
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(semconv.ServiceName(cfg.ServiceName)),
		resource.WithTelemetrySDK(),
	)
	if err != nil {
		return nil, fmt.Errorf("telemetry resource: %w", err)
	}

	shutdownFuncs := make([]func(context.Context) error, 0, 3)
	shutdownAll := func(ctx context.Context) error {
		var errs []error
		for _, fn := range shutdownFuncs {
			errs = append(errs, fn(ctx))
		}
		return errors.Join(errs...)
	}

	// --- Traces ---
	traceOpts := []otlptracehttp.Option{otlptracehttp.WithEndpointURL(cfg.OTLPEndpoint)}
	if len(cfg.OTLPHeaders) > 0 {
		traceOpts = append(traceOpts, otlptracehttp.WithHeaders(cfg.OTLPHeaders))
	}
	traceExp, err := otlptracehttp.New(ctx, traceOpts...)
	if err != nil {
		return shutdownAll, fmt.Errorf("trace exporter: %w", err)
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExp, sdktrace.WithBatchTimeout(5*time.Second)),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	shutdownFuncs = append(shutdownFuncs, tp.Shutdown)

	// --- Metrics ---
	metricOpts := []otlpmetrichttp.Option{otlpmetrichttp.WithEndpointURL(cfg.OTLPEndpoint)}
	if len(cfg.OTLPHeaders) > 0 {
		metricOpts = append(metricOpts, otlpmetrichttp.WithHeaders(cfg.OTLPHeaders))
	}
	metricExp, err := otlpmetrichttp.New(ctx, metricOpts...)
	if err != nil {
		return shutdownAll, fmt.Errorf("metric exporter: %w", err)
	}
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExp, sdkmetric.WithInterval(30*time.Second))),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(mp)
	shutdownFuncs = append(shutdownFuncs, mp.Shutdown)

	// --- Logs ---
	logOpts := []otlploghttp.Option{otlploghttp.WithEndpointURL(cfg.OTLPEndpoint)}
	if len(cfg.OTLPHeaders) > 0 {
		logOpts = append(logOpts, otlploghttp.WithHeaders(cfg.OTLPHeaders))
	}
	logExp, err := otlploghttp.New(ctx, logOpts...)
	if err != nil {
		return shutdownAll, fmt.Errorf("log exporter: %w", err)
	}
	lp := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExp, sdklog.WithExportTimeout(5*time.Second))),
		sdklog.WithResource(res),
	)
	otellog.SetLoggerProvider(lp)
	shutdownFuncs = append(shutdownFuncs, lp.Shutdown)

	return shutdownAll, nil
}
