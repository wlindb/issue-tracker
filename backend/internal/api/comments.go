package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	commentdomain "github.com/wlindb/issue-tracker/internal/domain/tracker/comment"

	"github.com/wlindb/issue-tracker/internal/api/model"
)

// CommentService is what the handler needs from the domain.
type CommentService interface {
	List(ctx context.Context, issueID uuid.UUID, query commentdomain.ListCommentQuery) (commentdomain.Comments, error)
	Create(ctx context.Context, id uuid.UUID, issueID uuid.UUID, authorID uuid.UUID, body string) (*commentdomain.Comment, error)
}

type CommentHandler struct {
	service CommentService
}

func NewCommentHandler(service CommentService) CommentHandler {
	return CommentHandler{service: service}
}

func (h *CommentHandler) ListComments(ctx context.Context, req model.ListCommentsRequestObject) (model.ListCommentsResponseObject, error) {
	if _, err := userIDFromContext(ctx); err != nil {
		return model.ListComments401JSONResponse{
			UnauthorizedJSONResponse: newUnauthorized("unauthorized", "authentication required"),
		}, nil
	}
	comments, err := h.service.List(ctx, req.IssueId, listCommentQueryFromRequest(req.Params))
	if errors.Is(err, commentdomain.ErrIssueNotFound) {
		return model.ListComments404JSONResponse{
			NotFoundJSONResponse: newNotFound("issue_not_found", "issue not found"),
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("list comments: %w", err)
	}
	return model.ListComments200JSONResponse{
		Items:      commentsFromDomain(comments.Items),
		NextCursor: nil,
	}, nil
}

func (h *CommentHandler) CreateComment(ctx context.Context, req model.CreateCommentRequestObject) (model.CreateCommentResponseObject, error) {
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return model.CreateComment401JSONResponse{
			UnauthorizedJSONResponse: newUnauthorized("unauthorized", "authentication required"),
		}, nil
	}
	if req.Body == nil || req.Body.Body == "" {
		return model.CreateComment400JSONResponse{
			BadRequestJSONResponse: newBadRequest("invalid_input", "body is required"),
		}, nil
	}
	comment, err := h.service.Create(ctx, uuid.New(), req.IssueId, userID, req.Body.Body)
	if errors.Is(err, commentdomain.ErrIssueNotFound) {
		return model.CreateComment404JSONResponse{
			NotFoundJSONResponse: newNotFound("issue_not_found", "issue not found"),
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("create comment: %w", err)
	}
	return model.CreateComment201JSONResponse(commentFromDomain(*comment)), nil
}

func (h *CommentHandler) DeleteComment(_ context.Context, _ model.DeleteCommentRequestObject) (model.DeleteCommentResponseObject, error) {
	return model.DeleteComment500JSONResponse{InternalServerErrorJSONResponse: model.InternalServerErrorJSONResponse(notImplemented())}, nil
}
