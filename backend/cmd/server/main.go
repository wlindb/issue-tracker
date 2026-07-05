// Package main is the composition root for the issue-tracker server.
package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
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

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"

	"github.com/wlindb/issue-tracker/internal/api"
	"github.com/wlindb/issue-tracker/internal/application/auth"
	"github.com/wlindb/issue-tracker/internal/application/embedding"
	"github.com/wlindb/issue-tracker/internal/config"
	commentdomain "github.com/wlindb/issue-tracker/internal/domain/tracker/comment"
	"github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
	"github.com/wlindb/issue-tracker/internal/domain/tracker/label"
	trackerdomain "github.com/wlindb/issue-tracker/internal/domain/tracker/project"
	workspacedomain "github.com/wlindb/issue-tracker/internal/domain/tracker/workspace"
	"github.com/wlindb/issue-tracker/internal/infrastructure/db"
	trackerinfra "github.com/wlindb/issue-tracker/internal/infrastructure/tracker"
	keycloakpkg "github.com/wlindb/issue-tracker/internal/pkg/keycloak"
	embeddednats "github.com/wlindb/issue-tracker/internal/pkg/nats"
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

	otelCloser, err := newOtel(ctx, *cfg)
	if err != nil {
		return fmt.Errorf("create otel: %w", err)
	}
	defer otelCloser()

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

	tracer := otel.Tracer(cfg.OTELServiceName)

	setup, err := newNATSConnection(cfg.NATSPort, cfg.NATSWebSocketPort)
	if err != nil {
		return fmt.Errorf("create nats connection: %w", err)
	}
	defer setup.closer()

	if err := newEventHandlers(setup.connection); err != nil {
		return err
	}

	workspaceService := workspacedomain.NewWorkspaceService(
		trackerinfra.NewTracingWorkspaceRepository(trackerinfra.NewWorkspaceRepository(pool), tracer),
	)

	if err := newNATSAuthCallout(setup, workspaceService, cfg.JWKSUrl); err != nil {
		return fmt.Errorf("nats auth callout: %w", err)
	}

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

func newOtel(ctx context.Context, cfg config.Config) (func(), error) {
	var zero func()

	otelShutdown, err := telemetry.Setup(ctx, telemetry.Config{
		ServiceName: cfg.OTELServiceName,
	})
	if err != nil {
		return zero, fmt.Errorf("telemetry: %w", err)
	}
	closer := func() {
		shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := otelShutdown(shutdownCtx); err != nil {
			log.Printf("telemetry shutdown: %v", err)
		}
	}
	return closer, nil
}

type natsSetup struct {
	connection    *nats.Conn
	issuerKeyPair nkeys.KeyPair
	closer        func()
}

func newNATSConnection(natsPort int, natsWebSocketPort int) (natsSetup, error) {
	var zero natsSetup

	issuerKeyPair, err := nkeys.CreateAccount()
	if err != nil {
		return zero, fmt.Errorf("create nats issuer key pair: %w", err)
	}
	issuerPublicKey, err := issuerKeyPair.PublicKey()
	if err != nil {
		return zero, fmt.Errorf("get nats issuer public key: %w", err)
	}

	internalUser := uuid.New().String()
	internalPassword := mustGenerateRandomHex(32)

	serverOptions := []embeddednats.ServerOption{
		embeddednats.WithInternalUser(internalUser, internalPassword),
		embeddednats.WithAuthCallout(embeddednats.AuthCalloutConfig{
			IssuerPublicKey: issuerPublicKey,
			AuthUsers:       []string{internalUser},
		}),
	}
	if natsPort > 0 {
		serverOptions = append(serverOptions, embeddednats.WithExternalPort(natsPort))
	}
	if natsWebSocketPort > 0 {
		serverOptions = append(serverOptions, embeddednats.WithWebSocketPort(natsWebSocketPort))
	}

	natsServer, err := embeddednats.StartEmbeddedServer(serverOptions...)
	if err != nil {
		return zero, fmt.Errorf("embedded nats: %w", err)
	}
	if natsWebSocketPort > 0 {
		log.Printf("NATS WebSocket ready at %s", natsServer.WebsocketURL())
	}

	natsConnection, err := embeddednats.Connect(natsServer, nats.UserInfo(internalUser, internalPassword))
	if err != nil {
		natsServer.Shutdown()
		return zero, fmt.Errorf("nats connect: %w", err)
	}

	return natsSetup{
		connection:    natsConnection,
		issuerKeyPair: issuerKeyPair,
		closer: func() {
			natsServer.Shutdown()
			natsConnection.Close()
		},
	}, nil
}

func newNATSAuthCallout(setup natsSetup, workspaceService *workspacedomain.WorkspaceService, jwksURL string) error {
	tokenValidator, err := keycloakpkg.NewKeycloakTokenValidator(jwksURL)
	if err != nil {
		return fmt.Errorf("create keycloak token validator: %w", err)
	}
	handler := auth.NewAuthCalloutHandler(tokenValidator, workspaceService, setup.issuerKeyPair)
	if _, err := setup.connection.Subscribe(embeddednats.AuthCalloutSubject, handler.Handle); err != nil {
		return fmt.Errorf("subscribe to nats auth callout subject: %w", err)
	}
	return nil
}

func mustGenerateRandomHex(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		panic(fmt.Errorf("generate random hex: %w", err))
	}
	return hex.EncodeToString(bytes)
}

func newEventHandlers(connection *nats.Conn) error {
	if _, err := embedding.NewEmbeddingHandler(
		embedding.WithIssueCreated(
			embeddednats.NewNATSEventSubscriber[issue.IssueCreatedEvent](connection, embeddednats.IssueCreatedSubjectAll),
		),
	); err != nil {
		return fmt.Errorf("create embedding handler: %w", err)
	}
	if err := trackerinfra.NewEventPublisher(connection); err != nil {
		return fmt.Errorf("create event publisher: %w", err)
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
	labelRepository := trackerinfra.NewTracingLabelRepository(
		trackerinfra.NewLabelRepository(pool),
		tracer,
	)
	commentRepository := trackerinfra.NewTracingCommentRepository(
		trackerinfra.NewCommentRepository(pool),
		tracer,
	)

	tracingLabelRepository := trackerinfra.NewTracingLabelRepository(labelRepository, tracer)
	labelService := label.NewLabelService(tracingLabelRepository)
	tracingLabelService := trackerinfra.NewTracingLabelService(labelService, tracer)

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
				issue.NewIssueService(trackerinfra.NewUoW(pool), issueRepository, labelRepository),
				tracer,
			),
		),
		CommentHandler: api.NewCommentHandler(
			trackerinfra.NewTracingCommentService(
				commentdomain.NewCommentService(commentRepository),
				tracer,
			),
		),
		LabelHandler: api.NewLabelHandler(tracingLabelService),
	}
}
