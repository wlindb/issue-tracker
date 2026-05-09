package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL     string
	ServerAddr      string
	JWKSUrl         string
	OTELServiceName string
	// NATSPort is the port the embedded NATS server listens on for external clients.
	// 0 means loopback-only on a random port (internal use only).
	NATSPort int
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

	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	if serviceName == "" {
		serviceName = "issue-tracker"
	}

	natsPort := 0
	if raw := os.Getenv("NATS_PORT"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			return nil, fmt.Errorf("NATS_PORT must be a valid integer: %w", err)
		}
		natsPort = parsed
	}

	return &Config{
		DatabaseURL:     databaseURL,
		ServerAddr:      serverAddr,
		JWKSUrl:         jwksURL,
		OTELServiceName: serviceName,
		NATSPort:        natsPort,
	}, nil
}
