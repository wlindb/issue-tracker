package api

import (
	"context"

	"github.com/wlindb/issue-tracker/internal/api/model"
)

func (h *Handler) ListComments(_ context.Context, _ model.ListCommentsRequestObject) (model.ListCommentsResponseObject, error) {
	return model.ListComments500JSONResponse{InternalServerErrorJSONResponse: model.InternalServerErrorJSONResponse(notImplemented())}, nil
}

func (h *Handler) CreateComment(_ context.Context, _ model.CreateCommentRequestObject) (model.CreateCommentResponseObject, error) {
	return model.CreateComment500JSONResponse{InternalServerErrorJSONResponse: model.InternalServerErrorJSONResponse(notImplemented())}, nil
}

func (h *Handler) DeleteComment(_ context.Context, _ model.DeleteCommentRequestObject) (model.DeleteCommentResponseObject, error) {
	return model.DeleteComment500JSONResponse{InternalServerErrorJSONResponse: model.InternalServerErrorJSONResponse(notImplemented())}, nil
}
