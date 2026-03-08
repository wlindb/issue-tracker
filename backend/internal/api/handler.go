package api

import (
	"github.com/wlindb/issue-tracker/internal/api/generated"
	"github.com/wlindb/issue-tracker/internal/auth"
)

// Handler implements generated.StrictServerInterface.
type Handler struct {
	Auth *auth.Service
}

var _ generated.StrictServerInterface = (*Handler)(nil)

func notImplemented() generated.Error {
	return generated.Error{Code: "not_implemented", Message: "not implemented"}
}
