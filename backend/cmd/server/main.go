// Package main is the composition root for the issue-tracker server.
package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/wlindb/issue-tracker/internal/api"
	"github.com/wlindb/issue-tracker/internal/api/generated"
)

//go:embed static
var staticFiles embed.FS

func main() {
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

	log.Fatal(e.Start(":8080"))
}
