package config

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL     string
	ServerAddr      string
	JWKSUrl         string
	OTELServiceName string
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

	return &Config{
		DatabaseURL:     databaseURL,
		ServerAddr:      serverAddr,
		JWKSUrl:         jwksURL,
		OTELServiceName: serviceName,
	}, nil
}
