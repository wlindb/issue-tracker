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
	issuedomain "github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
	trackerdomain "github.com/wlindb/issue-tracker/internal/domain/tracker/project"
	"github.com/wlindb/issue-tracker/internal/infrastructure/db"
	trackerinfra "github.com/wlindb/issue-tracker/internal/infrastructure/tracker"
	"github.com/wlindb/issue-tracker/internal/pkg/telemetry"
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

	otelShutdown, err := telemetry.Setup(ctx, telemetry.Config{
		ServiceName: cfg.OTELServiceName,
	})
	if err != nil {
		return fmt.Errorf("telemetry: %w", err)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := otelShutdown(shutdownCtx); err != nil {
			log.Printf("telemetry shutdown: %v", err)
		}
	}()
	log.Println("telemetry initialised")

	pool, err := db.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("database: %w", err)
	}
	defer pool.Close()
	log.Println("database connected")

	if err := trackerinfra.Migrate(ctx, pool); err != nil {
		return fmt.Errorf("tracker migrate: %w", err)
	}
	log.Println("tracker migrations applied")

	projectRepository := trackerinfra.NewProjectRepository(pool)
	issueRepository := trackerinfra.NewIssueRepository(pool)
	h := &api.Handler{
		ProjectHandler: api.NewProjectHandler(trackerdomain.NewProjectService(projectRepository)),
		IssueHandler:   api.NewIssueHandler(issuedomain.NewIssueService(issueRepository)),
	}

	e, err := newServer(h, cfg)
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
