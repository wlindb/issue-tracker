// Package main is the composition root for the issue-tracker server.
package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/wlindb/issue-tracker/internal/api"
	"github.com/wlindb/issue-tracker/internal/api/generated"
	"github.com/wlindb/issue-tracker/internal/config"
	"github.com/wlindb/issue-tracker/internal/db"
)

//go:embed static
var staticFiles embed.FS

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	pool, err := db.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer pool.Close()
	log.Println("database connected")

	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	// Serve the embedded OpenAPI spec as JSON.
	e.GET("/openapi.json", func(c echo.Context) error {
		swagger, err := generated.GetSwagger()
		if err != nil {
			return fmt.Errorf("loading openapi spec: %w", err)
		}
		return c.JSON(http.StatusOK, swagger)
	})

	// Serve Swagger UI at /docs.
	e.FileFS("/docs", "docs.html", echo.MustSubFS(staticFiles, "static"))

	// Register all API handlers under /api/v1.
	h := &api.Handler{}
	strict := generated.NewStrictHandler(h, nil)
	generated.RegisterHandlersWithBaseURL(e, strict, "/api/v1")

	log.Fatal(e.Start(cfg.ServerAddr))
}
