package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	DatabaseURL     string
	ServerAddr      string
	JWKSUrl         string
	OTELServiceName string
	OTELEndpoint    string
	OTELHeaders     map[string]string
}

func Load() (*Config, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	serverAddr := os.Getenv("SERVER_ADDR")
	if serverAddr == "" {
		serverAddr = ":8080"
	}

	jwksURL := os.Getenv("JWKS_URL")
	if jwksURL == "" {
		return nil, fmt.Errorf("JWKS_URL is required")
	}

	otelEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if otelEndpoint == "" {
		return nil, fmt.Errorf("OTEL_EXPORTER_OTLP_ENDPOINT is required")
	}

	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	if serviceName == "" {
		serviceName = "issue-tracker"
	}

	return &Config{
		DatabaseURL:     databaseURL,
		ServerAddr:      serverAddr,
		JWKSUrl:         jwksURL,
		OTELServiceName: serviceName,
		OTELEndpoint:    otelEndpoint,
		OTELHeaders:     parseOTELHeaders(os.Getenv("OTEL_EXPORTER_OTLP_HEADERS")),
	}, nil
}

func parseOTELHeaders(raw string) map[string]string {
	if raw == "" {
		return nil
	}
	headers := make(map[string]string)
	for _, pair := range strings.Split(raw, ",") {
		k, v, ok := strings.Cut(pair, "=")
		if ok {
			headers[strings.TrimSpace(k)] = strings.TrimSpace(v)
		}
	}
	return headers
}
