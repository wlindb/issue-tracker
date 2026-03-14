package api

import (
	"context"

	"github.com/wlindb/issue-tracker/internal/api/generated"
)

func (h *Handler) ListComments(_ context.Context, _ generated.ListCommentsRequestObject) (generated.ListCommentsResponseObject, error) {
	return generated.ListComments500JSONResponse{InternalServerErrorJSONResponse: generated.InternalServerErrorJSONResponse(notImplemented())}, nil
}

func (h *Handler) CreateComment(_ context.Context, _ generated.CreateCommentRequestObject) (generated.CreateCommentResponseObject, error) {
	return generated.CreateComment500JSONResponse{InternalServerErrorJSONResponse: generated.InternalServerErrorJSONResponse(notImplemented())}, nil
}

func (h *Handler) UpdateComment(_ context.Context, _ generated.UpdateCommentRequestObject) (generated.UpdateCommentResponseObject, error) {
	return generated.UpdateComment500JSONResponse{InternalServerErrorJSONResponse: generated.InternalServerErrorJSONResponse(notImplemented())}, nil
}

func (h *Handler) DeleteComment(_ context.Context, _ generated.DeleteCommentRequestObject) (generated.DeleteCommentResponseObject, error) {
	return generated.DeleteComment500JSONResponse{InternalServerErrorJSONResponse: generated.InternalServerErrorJSONResponse(notImplemented())}, nil
}
