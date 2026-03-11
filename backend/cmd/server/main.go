// Package main is the composition root for the issue-tracker server.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wlindb/issue-tracker/internal/api"
	"github.com/wlindb/issue-tracker/internal/config"
	"github.com/wlindb/issue-tracker/internal/db"
	authdomain "github.com/wlindb/issue-tracker/internal/domain/auth"
	authinfra "github.com/wlindb/issue-tracker/internal/infrastructure/auth"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	pool, err := db.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("database: %w", err)
	}
	defer pool.Close()
	log.Println("database connected")

	if err = authinfra.Migrate(ctx, pool); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	log.Println("migrations applied")

	userRepository := authinfra.NewUserRepository(pool)

	authService, err := authdomain.NewAutService(userRepository, cfg.JWTPrivateKey)
	if err != nil {
		return fmt.Errorf("auth service: %w", err)
	}

	h := &api.Handler{
		AuthHandler:    api.NewAuthHandler(authService),
		ProjectHandler: api.NewProjectHandler(nil), // TODO: wire real service
	}

	e, err := newServer(h)
	if err != nil {
		return fmt.Errorf("server: %w", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := e.Shutdown(shutdownCtx); err != nil {
			e.Logger.Fatal(err)
		}
	}()

	if err := e.Start(cfg.ServerAddr); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server: %w", err)
	}
	return nil
}
