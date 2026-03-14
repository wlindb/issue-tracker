package api

import (
	"context"

	"github.com/wlindb/issue-tracker/internal/api/generated"
)

func (h *Handler) GetMe(_ context.Context, _ generated.GetMeRequestObject) (generated.GetMeResponseObject, error) {
	return generated.GetMe500JSONResponse{InternalServerErrorJSONResponse: generated.InternalServerErrorJSONResponse(notImplemented())}, nil
}
