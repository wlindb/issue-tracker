package api_test

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/wlindb/issue-tracker/internal/api"
)

// injectUser returns an Echo middleware that injects a fixed user UUID into the
// request context, simulating what a future JWT middleware would do.
func injectUser(id uuid.UUID) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := api.WithUserID(c.Request().Context(), id)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}
