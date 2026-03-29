package main

import (
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"

	"github.com/wlindb/issue-tracker/internal/api"
	apimiddleware "github.com/wlindb/issue-tracker/internal/api/middleware"
	"github.com/wlindb/issue-tracker/internal/api/model"
	"github.com/wlindb/issue-tracker/internal/config"
)

func newServer(h *api.Handler, cfg *config.Config) (*echo.Echo, error) {
	logger := slog.New(slog.NewMultiHandler(
		slog.NewJSONHandler(os.Stdout, nil),
		otelslog.NewHandler(cfg.OTELServiceName),
	))

	e := echo.New()
	e.HTTPErrorHandler = api.HTTPErrorHandler
	e.Use(otelecho.Middleware(cfg.OTELServiceName))
	e.Use(api.RequestLogger(logger))
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())

	api.RegisterPublicRoutes(e)

	protected := e.Group("", apimiddleware.JwtMiddleware(cfg.JWKSUrl), apimiddleware.UserIDMiddleware())
	strict := model.NewStrictHandler(h, nil)
	model.RegisterHandlersWithBaseURL(protected, strict, "/api/v1")

	return e, nil
}
