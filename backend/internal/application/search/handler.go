// Package search contains Echo HTTP handlers for the search module,
// implementing a subset of model.StrictServerInterface generated from
// api/openapi.yaml via oapi-codegen.
package search

import (
	"context"

	"github.com/wlindb/issue-tracker/internal/application/api/model"
)

// SearchHandler implements the search module's StrictServerInterface methods.
type SearchHandler struct{}

// NewSearchHandler creates a SearchHandler.
func NewSearchHandler() SearchHandler {
	return SearchHandler{}
}

func (h SearchHandler) SearchIssues(_ context.Context, req model.SearchIssuesRequestObject) (model.SearchIssuesResponseObject, error) {
	if req.Body == nil || req.Body.Query == "" {
		return model.SearchIssues400JSONResponse{
			BadRequestJSONResponse: model.BadRequestJSONResponse(model.Error{Code: "invalid_input", Message: "query is required"}),
		}, nil
	}
	return model.SearchIssues501JSONResponse{
		NotImplementedJSONResponse: model.NotImplementedJSONResponse(model.Error{Code: "not_implemented", Message: "not implemented"}),
	}, nil
}
