package api

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/wlindb/issue-tracker/internal/api/model"
)

//go:embed static
var staticFiles embed.FS

func RegisterPublicRoutes(e *echo.Echo) {
	e.GET("/openapi.json", func(c echo.Context) error {
		swagger, err := model.GetSwagger()
		if err != nil {
			return fmt.Errorf("loading openapi spec: %w", err)
		}
		return c.JSON(http.StatusOK, swagger)
	})

	e.FileFS("/public/docs", "docs.html", echo.MustSubFS(staticFiles, "static"))
}
