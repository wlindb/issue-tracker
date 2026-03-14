package api

import (
	"log"

	"github.com/MicahParks/keyfunc/v3"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

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
