package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/wlindb/issue-tracker/internal/api/generated"
)

// HTTPErrorHandler is a custom Echo error handler that formats all unhandled
// errors as the API's Error schema, keeping responses consistent with the
// OpenAPI spec. Expected errors (documented responses) are handled in
// individual handlers by returning typed response objects; only unexpected
// errors reach this handler.
func HTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	code := http.StatusInternalServerError
	body := generated.Error{
		Code:    "internal_error",
		Message: "an unexpected error occurred",
	}

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		body = generated.Error{
			Code:    http.StatusText(code),
			Message: fmt.Sprintf("%v", he.Message),
		}
	}

	_ = c.JSON(code, body)
}
