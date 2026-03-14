package api

import (
	"context"

	"github.com/wlindb/issue-tracker/internal/api/model"
)

func (h *Handler) ListIssues(_ context.Context, _ model.ListIssuesRequestObject) (model.ListIssuesResponseObject, error) {
	return model.ListIssues500JSONResponse{InternalServerErrorJSONResponse: model.InternalServerErrorJSONResponse(notImplemented())}, nil
}

func (h *Handler) CreateIssue(_ context.Context, _ model.CreateIssueRequestObject) (model.CreateIssueResponseObject, error) {
	return model.CreateIssue500JSONResponse{InternalServerErrorJSONResponse: model.InternalServerErrorJSONResponse(notImplemented())}, nil
}

func (h *Handler) GetIssue(_ context.Context, _ model.GetIssueRequestObject) (model.GetIssueResponseObject, error) {
	return model.GetIssue500JSONResponse{InternalServerErrorJSONResponse: model.InternalServerErrorJSONResponse(notImplemented())}, nil
}

func (h *Handler) UpdateIssue(_ context.Context, _ model.UpdateIssueRequestObject) (model.UpdateIssueResponseObject, error) {
	return model.UpdateIssue500JSONResponse{InternalServerErrorJSONResponse: model.InternalServerErrorJSONResponse(notImplemented())}, nil
}

func (h *Handler) DeleteIssue(_ context.Context, _ model.DeleteIssueRequestObject) (model.DeleteIssueResponseObject, error) {
	return model.DeleteIssue500JSONResponse{InternalServerErrorJSONResponse: model.InternalServerErrorJSONResponse(notImplemented())}, nil
}
