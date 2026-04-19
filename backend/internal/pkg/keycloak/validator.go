package keycloak

import (
	"context"
	"fmt"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TokenValidator validates a raw JWT string and returns the user ID extracted from the sub claim.
type TokenValidator interface {
	ValidateToken(ctx context.Context, rawToken string) (uuid.UUID, error)
}

// KeycloakTokenValidator validates JWTs against a Keycloak JWKS endpoint.
type KeycloakTokenValidator struct {
	keyfunc jwt.Keyfunc
}

// NewKeycloakTokenValidator creates a validator that fetches signing keys from jwksURL.
func NewKeycloakTokenValidator(jwksURL string) (*KeycloakTokenValidator, error) {
	jwks, err := keyfunc.NewDefault([]string{jwksURL})
	if err != nil {
		return nil, fmt.Errorf("create jwks keyfunc: %w", err)
	}
	return &KeycloakTokenValidator{keyfunc: jwks.Keyfunc}, nil
}

// ValidateToken validates rawToken, extracts the sub claim, and returns it as a uuid.UUID.
func (v *KeycloakTokenValidator) ValidateToken(_ context.Context, rawToken string) (uuid.UUID, error) {
	var zero uuid.UUID
	token, err := jwt.Parse(rawToken, v.keyfunc, jwt.WithExpirationRequired())
	if err != nil {
		return zero, fmt.Errorf("validate token: %w", err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return zero, fmt.Errorf("unexpected claims type")
	}
	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		return zero, fmt.Errorf("missing sub claim")
	}
	userID, err := uuid.Parse(sub)
	if err != nil {
		return zero, fmt.Errorf("sub claim is not a valid UUID: %w", err)
	}
	return userID, nil
}
