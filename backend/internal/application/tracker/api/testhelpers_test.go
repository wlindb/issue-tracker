package api_test

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/wlindb/issue-tracker/internal/application/tracker/api"
)

// testWorkspaceID is a fixed workspace UUID used across contract tests.
// It is used only for routing (no workspace middleware runs in unit tests).
var testWorkspaceID = uuid.MustParse("00000000-0000-0000-0000-000000000099")

// wsPath returns the full API path prefixed with the test workspace segment.
func wsPath(path string) string {
	return "/api/v1/workspaces/" + testWorkspaceID.String() + path
}

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

// injectUserClaims returns an Echo middleware that injects fixed JWT profile
// claims into the request context, simulating what UserIDMiddleware would do.
func injectUserClaims(claims api.UserClaims) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := api.WithUserClaims(c.Request().Context(), claims)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}
