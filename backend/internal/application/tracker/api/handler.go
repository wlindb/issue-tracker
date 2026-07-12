package api

import "github.com/wlindb/issue-tracker/internal/application/api/model"

// Handler composes per-domain-area handlers for the tracker module.
// Stub methods for unimplemented areas remain directly on Handler.
// It implements a subset of model.StrictServerInterface — the full interface
// is satisfied by the composition root in internal/application/api, which
// also embeds handlers from other modules (e.g. search).
type Handler struct {
	WorkspaceHandler
	ProjectHandler
	CommentHandler
	IssueHandler
	LabelHandler
	UserHandler
}

func notImplemented() model.Error {
	return model.Error{Code: "not_implemented", Message: "not implemented"}
}
