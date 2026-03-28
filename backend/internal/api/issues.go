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

func (h *Handler) UpdateIssueTitle(_ context.Context, _ model.UpdateIssueTitleRequestObject) (model.UpdateIssueTitleResponseObject, error) {
	return model.UpdateIssueTitle500JSONResponse{InternalServerErrorJSONResponse: model.InternalServerErrorJSONResponse(notImplemented())}, nil
}

func (h *Handler) UpdateIssueDescription(_ context.Context, _ model.UpdateIssueDescriptionRequestObject) (model.UpdateIssueDescriptionResponseObject, error) {
	return model.UpdateIssueDescription500JSONResponse{InternalServerErrorJSONResponse: model.InternalServerErrorJSONResponse(notImplemented())}, nil
}

func (h *Handler) UpdateIssueStatus(_ context.Context, _ model.UpdateIssueStatusRequestObject) (model.UpdateIssueStatusResponseObject, error) {
	return model.UpdateIssueStatus500JSONResponse{InternalServerErrorJSONResponse: model.InternalServerErrorJSONResponse(notImplemented())}, nil
}

func (h *Handler) UpdateIssuePriority(_ context.Context, _ model.UpdateIssuePriorityRequestObject) (model.UpdateIssuePriorityResponseObject, error) {
	return model.UpdateIssuePriority500JSONResponse{InternalServerErrorJSONResponse: model.InternalServerErrorJSONResponse(notImplemented())}, nil
}

func (h *Handler) UpdateIssueAssignee(_ context.Context, _ model.UpdateIssueAssigneeRequestObject) (model.UpdateIssueAssigneeResponseObject, error) {
	return model.UpdateIssueAssignee500JSONResponse{InternalServerErrorJSONResponse: model.InternalServerErrorJSONResponse(notImplemented())}, nil
}

func (h *Handler) DeleteIssue(_ context.Context, _ model.DeleteIssueRequestObject) (model.DeleteIssueResponseObject, error) {
	return model.DeleteIssue500JSONResponse{InternalServerErrorJSONResponse: model.InternalServerErrorJSONResponse(notImplemented())}, nil
}
