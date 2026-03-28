package api

import "github.com/wlindb/issue-tracker/internal/api/model"

// newUnauthorized constructs an UnauthorizedJSONResponse with the given error
// code and message.  It is embedded in per-endpoint 401 response types.
func newUnauthorized(code, message string) model.UnauthorizedJSONResponse {
	return model.UnauthorizedJSONResponse(model.Error{Code: code, Message: message})
}

// newBadRequest constructs a BadRequestJSONResponse with the given error code
// and message.  It is embedded in per-endpoint 400 response types.
func newBadRequest(code, message string) model.BadRequestJSONResponse {
	return model.BadRequestJSONResponse(model.Error{Code: code, Message: message})
}

// newNotFound constructs a NotFoundJSONResponse with the given error code and
// message.  It is embedded in per-endpoint 404 response types.
func newNotFound(code, message string) model.NotFoundJSONResponse {
	return model.NotFoundJSONResponse(model.Error{Code: code, Message: message})
}
