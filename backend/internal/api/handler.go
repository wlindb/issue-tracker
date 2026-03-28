package api

import (
	"github.com/wlindb/issue-tracker/internal/api/model"
)

// Handler composes per-domain-area handlers and implements StrictServerInterface.
// Stub methods for unimplemented areas remain directly on Handler.
type Handler struct {
	ProjectHandler
	CommentHandler
	IssueHandler
}

var _ model.StrictServerInterface = (*Handler)(nil)

func notImplemented() model.Error {
	return model.Error{Code: "not_implemented", Message: "not implemented"}
}
