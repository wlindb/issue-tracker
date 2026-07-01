package api

import (
	"context"

	"github.com/wlindb/issue-tracker/internal/api/model"
)

// LabelHandler handles label-related HTTP requests.
type LabelHandler struct{}

// NewLabelHandler returns a new LabelHandler.
func NewLabelHandler() LabelHandler {
	return LabelHandler{}
}

func (h *LabelHandler) CreateLabel(_ context.Context, req model.CreateLabelRequestObject) (model.CreateLabelResponseObject, error) {
	if req.Body == nil || req.Body.Name == "" {
		return model.CreateLabel400JSONResponse{
			BadRequestJSONResponse: newBadRequest("invalid_input", "name is required"),
		}, nil
	}
	return model.CreateLabel201JSONResponse(model.Label{}), nil
}

func (h *LabelHandler) ListLabels(_ context.Context, _ model.ListLabelsRequestObject) (model.ListLabelsResponseObject, error) {
	return model.ListLabels200JSONResponse{Items: []model.Label{}}, nil
}
