package api

import (
	"context"

	"github.com/wlindb/issue-tracker/internal/api/generated"
)

func (h *Handler) ListComments(ctx context.Context, request generated.ListCommentsRequestObject) (generated.ListCommentsResponseObject, error) {
	return generated.ListComments500JSONResponse{InternalServerErrorJSONResponse: generated.InternalServerErrorJSONResponse(notImplemented())}, nil
}

func (h *Handler) CreateComment(ctx context.Context, request generated.CreateCommentRequestObject) (generated.CreateCommentResponseObject, error) {
	return generated.CreateComment500JSONResponse{InternalServerErrorJSONResponse: generated.InternalServerErrorJSONResponse(notImplemented())}, nil
}

func (h *Handler) UpdateComment(ctx context.Context, request generated.UpdateCommentRequestObject) (generated.UpdateCommentResponseObject, error) {
	return generated.UpdateComment500JSONResponse{InternalServerErrorJSONResponse: generated.InternalServerErrorJSONResponse(notImplemented())}, nil
}

func (h *Handler) DeleteComment(ctx context.Context, request generated.DeleteCommentRequestObject) (generated.DeleteCommentResponseObject, error) {
	return generated.DeleteComment500JSONResponse{InternalServerErrorJSONResponse: generated.InternalServerErrorJSONResponse(notImplemented())}, nil
}
