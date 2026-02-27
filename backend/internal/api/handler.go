package api

import "github.com/wlindb/issue-tracker/internal/api/generated"

// Handler implements generated.StrictServerInterface.
// All methods return 501 Not Implemented until real logic is wired in.
type Handler struct{}

var _ generated.StrictServerInterface = (*Handler)(nil)

func notImplemented() generated.Error {
	return generated.Error{Code: "not_implemented", Message: "not implemented"}
}
