package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL          string
	SearchDatabaseURL    string
	ServerAddr           string
	JWKSUrl              string
	OTELServiceName      string
	GeminiAPIKey         string
	GeminiEmbeddingModel string
	ResendAPIKey         string
	ResendFromEmail      string
	// NATSPort is the port the embedded NATS server listens on for external clients.
	// 0 means loopback-only on a random port (internal use only).
	NATSPort int
	// NATSWebSocketPort is the port the embedded NATS server listens on for WebSocket clients.
	// 0 means WebSocket is disabled.
	NATSWebSocketPort int
}

func Load() (*Config, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	searchDatabaseURL := os.Getenv("SEARCH_DATABASE_URL")
	if searchDatabaseURL == "" {
		return nil, fmt.Errorf("SEARCH_DATABASE_URL is required")
	}

	serverAddr := os.Getenv("SERVER_ADDR")
	if serverAddr == "" {
		serverAddr = ":8080"
	}

	jwksURL := os.Getenv("JWKS_URL")
	if jwksURL == "" {
		return nil, fmt.Errorf("JWKS_URL is required")
	}

	geminiAPIKey := os.Getenv("GEMINI_API_KEY")
	if geminiAPIKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is required")
	}

	geminiEmbeddingModel := os.Getenv("GEMINI_EMBEDDING_MODEL")
	if geminiEmbeddingModel == "" {
		return nil, fmt.Errorf("GEMINI_EMBEDDING_MODEL is required")
	}

	resendAPIKey := os.Getenv("RESEND_API_KEY")
	if resendAPIKey == "" {
		return nil, fmt.Errorf("RESEND_API_KEY is required")
	}

	resendFromEmail := os.Getenv("RESEND_FROM_EMAIL")
	if resendFromEmail == "" {
		return nil, fmt.Errorf("RESEND_FROM_EMAIL is required")
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

	natsWebSocketPort := 0
	if raw := os.Getenv("NATS_WEBSOCKET_PORT"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			return nil, fmt.Errorf("NATS_WEBSOCKET_PORT must be a valid integer: %w", err)
		}
		natsWebSocketPort = parsed
	}

	return &Config{
		DatabaseURL:          databaseURL,
		SearchDatabaseURL:    searchDatabaseURL,
		ServerAddr:           serverAddr,
		JWKSUrl:              jwksURL,
		OTELServiceName:      serviceName,
		GeminiAPIKey:         geminiAPIKey,
		GeminiEmbeddingModel: geminiEmbeddingModel,
		ResendAPIKey:         resendAPIKey,
		ResendFromEmail:      resendFromEmail,
		NATSPort:             natsPort,
		NATSWebSocketPort:    natsWebSocketPort,
	}, nil
}
