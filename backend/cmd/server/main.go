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

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/wlindb/issue-tracker/internal/api"
	"github.com/wlindb/issue-tracker/internal/config"
	issuedomain "github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
	trackerdomain "github.com/wlindb/issue-tracker/internal/domain/tracker/project"
	workspacedomain "github.com/wlindb/issue-tracker/internal/domain/tracker/workspace"
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

	// Run migrations as the superuser before creating the restricted app pool.
	// Migrations may create roles (e.g. appuser) that the app pool depends on.
	migrationPool, err := db.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("database (migration pool): %w", err)
	}
	if err := trackerinfra.Migrate(ctx, migrationPool); err != nil {
		migrationPool.Close()
		return fmt.Errorf("tracker migrate: %w", err)
	}
	migrationPool.Close()
	log.Println("tracker migrations applied")

	pool, err := db.New(ctx, cfg.DatabaseURL,
		db.WithAppSessionVars(
			api.WorkspaceIDFromContext,
			api.UserIDFromContext,
		),
		db.WithAppRole("appuser"),
	)
	if err != nil {
		return fmt.Errorf("database: %w", err)
	}
	defer pool.Close()
	log.Println("database connected")

	tracer := otel.Tracer(cfg.OTELServiceName)

	workspaceService := workspacedomain.NewWorkspaceService(
		trackerinfra.NewTracingWorkspaceRepository(trackerinfra.NewWorkspaceRepository(pool), tracer),
	)
	h := newHandler(pool, tracer, workspaceService)

	e, err := newServer(h, cfg, workspaceService)
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

func newHandler(pool *pgxpool.Pool, tracer trace.Tracer, workspaceService *workspacedomain.WorkspaceService) *api.Handler {
	projectRepository := trackerinfra.NewTracingProjectRepository(
		trackerinfra.NewProjectRepository(pool),
		tracer,
	)
	issueRepository := trackerinfra.NewTracingIssueRepository(
		trackerinfra.NewIssueRepository(pool),
		tracer,
	)
	return &api.Handler{
		WorkspaceHandler: api.NewWorkspaceHandler(workspaceService),
		ProjectHandler: api.NewProjectHandler(
			trackerinfra.NewTracingProjectService(
				trackerdomain.NewProjectService(projectRepository),
				tracer,
			),
		),
		IssueHandler: api.NewIssueHandler(
			trackerinfra.NewTracingIssueService(
				issuedomain.NewIssueService(issueRepository),
				tracer,
			),
		),
	}
}
