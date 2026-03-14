package main

import (
	"embed"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/wlindb/issue-tracker/internal/api"
	apimiddleware "github.com/wlindb/issue-tracker/internal/api/middleware"
	"github.com/wlindb/issue-tracker/internal/api/model"
	"github.com/wlindb/issue-tracker/internal/config"
)

//go:embed static
var staticFiles embed.FS

func newServer(h *api.Handler, cfg *config.Config) (*echo.Echo, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	e := echo.New()
	e.HTTPErrorHandler = api.HTTPErrorHandler
	e.Use(api.RequestLogger(logger))
	e.Use(echomiddleware.Recover())

	e.GET("/openapi.json", func(c echo.Context) error {
		swagger, err := model.GetSwagger()
		if err != nil {
			return fmt.Errorf("loading openapi spec: %w", err)
		}
		return c.JSON(http.StatusOK, swagger)
	})

	e.FileFS("/docs", "docs.html", echo.MustSubFS(staticFiles, "static"))

	e.Use(apimiddleware.JwtMiddleware(cfg.JWKSUrl))
	e.Use(apimiddleware.UserIDMiddleware())

	strict := model.NewStrictHandler(h, nil)
	model.RegisterHandlersWithBaseURL(e, strict, "/api/v1")

	return e, nil
}
