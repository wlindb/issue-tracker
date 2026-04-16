package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/wlindb/issue-tracker/internal/api/model"
	trackerdomain "github.com/wlindb/issue-tracker/internal/domain/tracker/project"
)

// ProjectService is what the handler needs from the domain.
type ProjectService interface {
	Create(ctx context.Context, command trackerdomain.CreateProjectCommand) (trackerdomain.Project, error)
	List(ctx context.Context, query trackerdomain.ListProjectQuery) (trackerdomain.Projects, error)
	Get(ctx context.Context, id uuid.UUID) (trackerdomain.Project, error)
}

type ProjectHandler struct {
	service ProjectService
}

func NewProjectHandler(service ProjectService) ProjectHandler {
	return ProjectHandler{service: service}
}

func (h *ProjectHandler) ListProjects(ctx context.Context, req model.ListProjectsRequestObject) (model.ListProjectsResponseObject, error) {
	if _, err := userIDFromContext(ctx); err != nil {
		return model.ListProjects401JSONResponse{
			UnauthorizedJSONResponse: newUnauthorized("unauthorized", "authentication required"),
		}, nil
	}
	projects, err := h.service.List(ctx, listProjectQueryFromRequest(req.Params))
	if err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}
	return model.ListProjects200JSONResponse{
		Items:      projectsFromDomain(projects.Items),
		NextCursor: nil,
	}, nil
}

func (h *ProjectHandler) CreateProject(ctx context.Context, req model.CreateProjectRequestObject) (model.CreateProjectResponseObject, error) {
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}
	if req.Body.Name == "" {
		return model.CreateProject400JSONResponse{
			BadRequestJSONResponse: newBadRequest("invalid_input", "name is required"),
		}, nil
	}
	command := trackerdomain.CreateProjectCommand{
		Name:        req.Body.Name,
		Description: req.Body.Description,
		OwnerID:     userID,
	}
	result, err := h.service.Create(ctx, command)
	if err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}
	return model.CreateProject201JSONResponse(projectFromDomain(result)), nil
}

func (h *ProjectHandler) GetProject(ctx context.Context, req model.GetProjectRequestObject) (model.GetProjectResponseObject, error) {
	if _, err := userIDFromContext(ctx); err != nil {
		return model.GetProject401JSONResponse{
			UnauthorizedJSONResponse: newUnauthorized("unauthorized", "authentication required"),
		}, nil
	}
	project, err := h.service.Get(ctx, req.ProjectId)
	if errors.Is(err, trackerdomain.ErrProjectNotFound) {
		return model.GetProject404JSONResponse{
			NotFoundJSONResponse: newNotFound("not_found", "project not found"),
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get project: %w", err)
	}
	return model.GetProject200JSONResponse(projectFromDomain(project)), nil
}

func (h *Handler) UpdateProject(_ context.Context, _ model.UpdateProjectRequestObject) (model.UpdateProjectResponseObject, error) {
	return model.UpdateProject500JSONResponse{InternalServerErrorJSONResponse: model.InternalServerErrorJSONResponse(notImplemented())}, nil
}

func (h *Handler) DeleteProject(_ context.Context, _ model.DeleteProjectRequestObject) (model.DeleteProjectResponseObject, error) {
	return model.DeleteProject500JSONResponse{InternalServerErrorJSONResponse: model.InternalServerErrorJSONResponse(notImplemented())}, nil
}
