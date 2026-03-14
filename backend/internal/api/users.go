package api

import (
	"context"

	"github.com/wlindb/issue-tracker/internal/api/model"
)

func (h *Handler) GetMe(_ context.Context, _ model.GetMeRequestObject) (model.GetMeResponseObject, error) {
	return model.GetMe500JSONResponse{InternalServerErrorJSONResponse: model.InternalServerErrorJSONResponse(notImplemented())}, nil
}
