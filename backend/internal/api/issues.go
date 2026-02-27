package api

import (
	"context"

	"github.com/wlindb/issue-tracker/internal/api/generated"
)

func (h *Handler) ListIssues(ctx context.Context, request generated.ListIssuesRequestObject) (generated.ListIssuesResponseObject, error) {
	return generated.ListIssues500JSONResponse{InternalServerErrorJSONResponse: generated.InternalServerErrorJSONResponse(notImplemented())}, nil
}

func (h *Handler) CreateIssue(ctx context.Context, request generated.CreateIssueRequestObject) (generated.CreateIssueResponseObject, error) {
	return generated.CreateIssue500JSONResponse{InternalServerErrorJSONResponse: generated.InternalServerErrorJSONResponse(notImplemented())}, nil
}

func (h *Handler) GetIssue(ctx context.Context, request generated.GetIssueRequestObject) (generated.GetIssueResponseObject, error) {
	return generated.GetIssue500JSONResponse{InternalServerErrorJSONResponse: generated.InternalServerErrorJSONResponse(notImplemented())}, nil
}

func (h *Handler) UpdateIssue(ctx context.Context, request generated.UpdateIssueRequestObject) (generated.UpdateIssueResponseObject, error) {
	return generated.UpdateIssue500JSONResponse{InternalServerErrorJSONResponse: generated.InternalServerErrorJSONResponse(notImplemented())}, nil
}

func (h *Handler) DeleteIssue(ctx context.Context, request generated.DeleteIssueRequestObject) (generated.DeleteIssueResponseObject, error) {
	return generated.DeleteIssue500JSONResponse{InternalServerErrorJSONResponse: generated.InternalServerErrorJSONResponse(notImplemented())}, nil
}
