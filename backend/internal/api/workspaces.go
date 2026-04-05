package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/wlindb/issue-tracker/internal/api/model"
	workspacedomain "github.com/wlindb/issue-tracker/internal/domain/tracker/workspace"
)

// WorkspaceService is what the handler needs from the domain.
type WorkspaceService interface {
	Create(ctx context.Context, ownerID uuid.UUID, name string) (*workspacedomain.Workspace, error)
	Get(ctx context.Context, id uuid.UUID) (*workspacedomain.Workspace, error)
	List(ctx context.Context, userID uuid.UUID) ([]workspacedomain.Workspace, error)
}

type WorkspaceHandler struct {
	service WorkspaceService
}

func NewWorkspaceHandler(service WorkspaceService) WorkspaceHandler {
	return WorkspaceHandler{service: service}
}

func (h *WorkspaceHandler) ListWorkspaces(ctx context.Context, _ model.ListWorkspacesRequestObject) (model.ListWorkspacesResponseObject, error) {
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return model.ListWorkspaces401JSONResponse{
			UnauthorizedJSONResponse: newUnauthorized("unauthorized", "authentication required"),
		}, nil
	}
	workspaces, err := h.service.List(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list workspaces: %w", err)
	}
	return model.ListWorkspaces200JSONResponse{
		Items:      workspacesFromDomain(workspaces),
		NextCursor: nil,
	}, nil
}

func (h *WorkspaceHandler) CreateWorkspace(ctx context.Context, req model.CreateWorkspaceRequestObject) (model.CreateWorkspaceResponseObject, error) {
	if req.Body.Name == "" {
		return model.CreateWorkspace400JSONResponse{
			BadRequestJSONResponse: newBadRequest("invalid_input", "name is required"),
		}, nil
	}
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("create workspace: %w", err)
	}
	workspace, err := h.service.Create(ctx, userID, req.Body.Name)
	if err != nil {
		return nil, fmt.Errorf("create workspace: %w", err)
	}
	return model.CreateWorkspace201JSONResponse(workspaceFromDomain(*workspace)), nil
}

func (h *WorkspaceHandler) GetWorkspace(ctx context.Context, req model.GetWorkspaceRequestObject) (model.GetWorkspaceResponseObject, error) {
	workspace, err := h.service.Get(ctx, req.WorkspaceId)
	if err != nil {
		if errors.Is(err, workspacedomain.ErrWorkspaceNotFound) {
			return model.GetWorkspace404JSONResponse{
				NotFoundJSONResponse: newNotFound("not_found", "workspace not found"),
			}, nil
		}
		return nil, fmt.Errorf("get workspace: %w", err)
	}
	return model.GetWorkspace200JSONResponse(workspaceFromDomain(*workspace)), nil
}
