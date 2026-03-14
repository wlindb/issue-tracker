package api

import (
	"context"

	"github.com/wlindb/issue-tracker/internal/api/generated"
)

func (h *Handler) ListIssues(_ context.Context, _ generated.ListIssuesRequestObject) (generated.ListIssuesResponseObject, error) {
	return generated.ListIssues500JSONResponse{InternalServerErrorJSONResponse: generated.InternalServerErrorJSONResponse(notImplemented())}, nil
}

func (h *Handler) CreateIssue(_ context.Context, _ generated.CreateIssueRequestObject) (generated.CreateIssueResponseObject, error) {
	return generated.CreateIssue500JSONResponse{InternalServerErrorJSONResponse: generated.InternalServerErrorJSONResponse(notImplemented())}, nil
}

func (h *Handler) GetIssue(_ context.Context, _ generated.GetIssueRequestObject) (generated.GetIssueResponseObject, error) {
	return generated.GetIssue500JSONResponse{InternalServerErrorJSONResponse: generated.InternalServerErrorJSONResponse(notImplemented())}, nil
}

func (h *Handler) UpdateIssue(_ context.Context, _ generated.UpdateIssueRequestObject) (generated.UpdateIssueResponseObject, error) {
	return generated.UpdateIssue500JSONResponse{InternalServerErrorJSONResponse: generated.InternalServerErrorJSONResponse(notImplemented())}, nil
}

func (h *Handler) DeleteIssue(_ context.Context, _ generated.DeleteIssueRequestObject) (generated.DeleteIssueResponseObject, error) {
	return generated.DeleteIssue500JSONResponse{InternalServerErrorJSONResponse: generated.InternalServerErrorJSONResponse(notImplemented())}, nil
}
