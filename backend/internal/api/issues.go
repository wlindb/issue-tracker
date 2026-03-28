package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/wlindb/issue-tracker/internal/api/model"
	issuedomain "github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
)

// Sentinel errors returned by IssueService implementations.
var (
	ErrIssueProjectNotFound = errors.New("issue: project not found")
	ErrIssueUnprocessable   = errors.New("issue: unprocessable entity")
)

// IssueService is what the handler needs from the domain.
type IssueService interface {
	ListIssues(ctx context.Context, projectID uuid.UUID, query issuedomain.ListIssueQuery) (issuedomain.IssuePage, error)
	CreateIssue(ctx context.Context, command issuedomain.CreateIssueCommand) (*issuedomain.Issue, error)
}

// IssueHandler holds the issue service dependency.
type IssueHandler struct {
	service IssueService
}

// NewIssueHandler creates an IssueHandler wired to the given service.
func NewIssueHandler(service IssueService) IssueHandler {
	return IssueHandler{service: service}
}

func (h *Handler) ListIssues(ctx context.Context, req model.ListIssuesRequestObject) (model.ListIssuesResponseObject, error) {
	if _, err := userIDFromContext(ctx); err != nil {
		return model.ListIssues401JSONResponse{
			UnauthorizedJSONResponse: newUnauthorized("unauthorized", "authentication required"),
		}, nil
	}
	query := listIssueQueryFromRequest(req.Params)
	page, err := h.IssueHandler.service.ListIssues(ctx, req.Params.ProjectId, query)
	if errors.Is(err, ErrIssueProjectNotFound) {
		return model.ListIssues404JSONResponse{
			NotFoundJSONResponse: newNotFound("not_found", "project not found"),
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("list issues: %w", err)
	}
	return model.ListIssues200JSONResponse{
		Items:      issuesFromDomain(page.Items),
		NextCursor: page.NextCursor,
	}, nil
}

func (h *Handler) CreateIssue(ctx context.Context, req model.CreateIssueRequestObject) (model.CreateIssueResponseObject, error) {
	if req.Body.ProjectId == uuid.Nil {
		return model.CreateIssue400JSONResponse{
			BadRequestJSONResponse: newBadRequest("invalid_input", "projectId is required"),
		}, nil
	}
	if req.Body.Title == "" {
		return model.CreateIssue400JSONResponse{
			BadRequestJSONResponse: newBadRequest("invalid_input", "title is required"),
		}, nil
	}
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return model.CreateIssue401JSONResponse{
			UnauthorizedJSONResponse: newUnauthorized("unauthorized", "authentication required"),
		}, nil
	}
	issue, err := h.IssueHandler.service.CreateIssue(ctx, createIssueCommandFromModel(req.Body.ProjectId, userID, *req.Body))
	if errors.Is(err, ErrIssueUnprocessable) {
		return model.CreateIssue422JSONResponse{
			UnprocessableEntityJSONResponse: newUnprocessable("unprocessable_entity", "request could not be processed"),
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("create issue: %w", err)
	}
	return model.CreateIssue201JSONResponse(issueFromDomain(*issue)), nil
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
