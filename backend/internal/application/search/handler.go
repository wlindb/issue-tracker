// Package search contains Echo HTTP handlers for the search module,
package search

import (
	"context"
	"fmt"

	"github.com/wlindb/issue-tracker/internal/application/api/model"
)

type SearchService interface {
	SearchIssues(ctx context.Context, description string) ([]model.Issue, error)
}

// SearchHandler implements the search module's StrictServerInterface methods.
type SearchHandler struct {
	service SearchService
}

// NewSearchHandler creates a SearchHandler.
func NewSearchHandler(service SearchService) SearchHandler {
	return SearchHandler{service: service}
}

func (h SearchHandler) SearchIssues(ctx context.Context, req model.SearchIssuesRequestObject) (model.SearchIssuesResponseObject, error) {
	if req.Body == nil || req.Body.Query == "" {
		return model.SearchIssues400JSONResponse{
			BadRequestJSONResponse: model.BadRequestJSONResponse(model.Error{Code: "invalid_input", Message: "query is required"}),
		}, nil
	}

	issues, err := h.service.SearchIssues(ctx, req.Body.Query)
	if err != nil {
		return nil, fmt.Errorf("search issues: %w", err)
	}

	return model.SearchIssues200JSONResponse{
		Items: issues,
	}, nil
}
