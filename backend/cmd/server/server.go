package main

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/wlindb/issue-tracker/internal/api"
	"github.com/wlindb/issue-tracker/internal/api/generated"
)

//go:embed static
var staticFiles embed.FS

func newServer(h *api.Handler) (*echo.Echo, error) {
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

	strict := generated.NewStrictHandler(h, nil)
	generated.RegisterHandlersWithBaseURL(e, strict, "/api/v1")

	return e, nil
}
