package api

import (
	"context"

	"github.com/wlindb/issue-tracker/internal/api/generated"
)

func (h *Handler) GetMe(ctx context.Context, request generated.GetMeRequestObject) (generated.GetMeResponseObject, error) {
	return generated.GetMe500JSONResponse{InternalServerErrorJSONResponse: generated.InternalServerErrorJSONResponse(notImplemented())}, nil
}
