package api

import (
	"github.com/wlindb/issue-tracker/internal/api/generated"
)

// Handler composes per-domain-area handlers and implements StrictServerInterface.
// Stub methods for unimplemented areas remain directly on Handler.
type Handler struct {
	AuthHandler
	ProjectHandler
}

var _ generated.StrictServerInterface = (*Handler)(nil)

func notImplemented() generated.Error {
	return generated.Error{Code: "not_implemented", Message: "not implemented"}
}
