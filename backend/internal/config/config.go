package config

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL string
	ServerAddr  string
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

	return &Config{
		DatabaseURL: databaseURL,
		ServerAddr:  serverAddr,
	}, nil
}
