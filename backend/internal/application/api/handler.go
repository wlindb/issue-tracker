// Package api is the composition root for the HTTP API: it embeds the
// per-module handlers (tracker, search, ...) into a single Handler that
// implements model.StrictServerInterface, generated from api/openapi.yaml.
package api

import (
	"github.com/wlindb/issue-tracker/internal/application/api/model"
	searchapi "github.com/wlindb/issue-tracker/internal/application/search"
	trackerapi "github.com/wlindb/issue-tracker/internal/application/tracker/api"
)

// Handler composes every module's handler into the full StrictServerInterface.
type Handler struct {
	trackerapi.Handler
	searchapi.SearchHandler
}

var _ model.StrictServerInterface = (*Handler)(nil)
