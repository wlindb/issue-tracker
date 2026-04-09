package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"

	"github.com/wlindb/issue-tracker/internal/api"
)

// WorkspaceMemberChecker validates that a user is a member of a workspace.
type WorkspaceMemberChecker interface {
	IsMember(ctx context.Context, workspaceID uuid.UUID, userID uuid.UUID) (bool, error)
}

// WorkspaceMembershipMiddleware reads the workspaceId path parameter, validates
// membership, and stores the workspace ID in context. Routes without a
// workspaceId path parameter (e.g. GET /workspaces) are passed through
// unchanged. Must be applied after UserIDMiddleware.
func WorkspaceMembershipMiddleware(checker WorkspaceMemberChecker) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			param := c.Param("workspaceId")
			if param == "" {
				return next(c)
			}
			workspaceID, err := uuid.Parse(param)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "workspaceId must be a valid UUID")
			}
			userID, err := api.UserIDFromContext(c.Request().Context())
			if err != nil {
				return echo.ErrUnauthorized
			}
			member, err := checker.IsMember(c.Request().Context(), workspaceID, userID)
			if err != nil {
				return fmt.Errorf("check workspace membership: %w", err)
			}
			if !member {
				return echo.ErrForbidden
			}
			ctx := api.WithWorkspaceID(c.Request().Context(), workspaceID)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

// UserIDMiddleware extracts the `sub` claim from the validated JWT (set by
// JwtMiddleware) and injects it as a uuid.UUID into the request context.
// Must be applied after JwtMiddleware.
func UserIDMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, ok := c.Get("user").(*jwt.Token)
			if !ok || token == nil {
				return echo.ErrUnauthorized
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return echo.ErrUnauthorized
			}
			sub, ok := claims["sub"].(string)
			if !ok || sub == "" {
				return echo.ErrUnauthorized
			}
			id, err := uuid.Parse(sub)
			if err != nil {
				return echo.ErrUnauthorized
			}
			ctx := api.WithUserID(c.Request().Context(), id)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

// JwtMiddleware returns an Echo middleware that validates Bearer JWTs against
// the JWKS served at jwksURL.
func JwtMiddleware(jwksURL string) echo.MiddlewareFunc {
	jwks, err := keyfunc.NewDefault([]string{jwksURL})
	if err != nil {
		log.Fatalf("failed to create JWK set from %s: %s", jwksURL, err)
	}
	return echojwt.WithConfig(echojwt.Config{
		KeyFunc: jwks.Keyfunc,
	})
}
