package main

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/wlindb/issue-tracker/internal/api"
	"github.com/wlindb/issue-tracker/internal/api/generated"
	"github.com/wlindb/issue-tracker/internal/auth"
	"github.com/wlindb/issue-tracker/internal/config"
	"github.com/wlindb/issue-tracker/internal/infrastructure/postgres"
)

//go:embed static
var staticFiles embed.FS

func newServer(cfg *config.Config, pool *pgxpool.Pool) (*echo.Echo, error) {
	userRepo := postgres.NewUserRepo(pool)

	authSvc, err := auth.New(userRepo, cfg.JWTPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("auth service: %w", err)
	}

	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.GET("/openapi.json", func(c echo.Context) error {
		swagger, err := generated.GetSwagger()
		if err != nil {
			return fmt.Errorf("loading openapi spec: %w", err)
		}
		return c.JSON(http.StatusOK, swagger)
	})

	e.FileFS("/docs", "docs.html", echo.MustSubFS(staticFiles, "static"))

	h := &api.Handler{Auth: authSvc}
	strict := generated.NewStrictHandler(h, nil)
	generated.RegisterHandlersWithBaseURL(e, strict, "/api/v1")

	return e, nil
}
