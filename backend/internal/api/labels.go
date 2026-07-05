package api

import (
	"context"
	"fmt"

	"github.com/wlindb/issue-tracker/internal/domain/tracker/label"

	"github.com/wlindb/issue-tracker/internal/api/model"
)

// LabelService is what the handler needs from the domain.
type LabelService interface {
	Create(ctx context.Context, name string) (label.Label, error)
	Search(ctx context.Context, name string) ([]label.Label, error)
}

// LabelHandler handles label-related HTTP requests.
type LabelHandler struct {
	service LabelService
}

// NewLabelHandler returns a new LabelHandler.
func NewLabelHandler(service LabelService) LabelHandler {
	return LabelHandler{service: service}
}

func (h *LabelHandler) CreateLabel(ctx context.Context, req model.CreateLabelRequestObject) (model.CreateLabelResponseObject, error) {
	if req.Body == nil || req.Body.Name == "" {
		return model.CreateLabel400JSONResponse{
			BadRequestJSONResponse: newBadRequest("invalid_input", "name is required"),
		}, nil
	}
	result, err := h.service.Create(ctx, req.Body.Name)
	if err != nil {
		return nil, fmt.Errorf("create label: %w", err)
	}
	return model.CreateLabel201JSONResponse(labelFromDomain(result)), nil
}

func (h *LabelHandler) ListLabels(ctx context.Context, req model.ListLabelsRequestObject) (model.ListLabelsResponseObject, error) {
	var query string
	if req.Params.Search != nil {
		query = *req.Params.Search
	}
	results, err := h.service.Search(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list labels: %w", err)
	}
	return model.ListLabels200JSONResponse{Items: labelsFromDomain(results)}, nil
}
