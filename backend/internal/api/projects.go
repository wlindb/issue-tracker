package api

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/wlindb/issue-tracker/internal/api/model"
	trackerdomain "github.com/wlindb/issue-tracker/internal/domain/tracker/project"
)

// ProjectService is what the handler needs from the domain.
type ProjectService interface {
	Create(ctx context.Context, project trackerdomain.Project) (trackerdomain.Project, error)
	List(ctx context.Context, query trackerdomain.ListProjectQuery) (trackerdomain.Projects, error)
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
	project, err := trackerdomain.New(uuid.New(), slugFromName(req.Body.Name), req.Body.Name, req.Body.Description, userID)
	if errors.Is(err, trackerdomain.ErrInvalidProject) {
		return model.CreateProject400JSONResponse{
			BadRequestJSONResponse: newBadRequest("invalid_input", "name is required"),
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}
	result, err := h.service.Create(ctx, project)
	if err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}
	return model.CreateProject201JSONResponse(projectFromDomain(result)), nil
}

func slugFromName(name string) string {
	return strings.ReplaceAll(strings.ToLower(name), " ", "-")
}

func (h *Handler) GetProject(_ context.Context, _ model.GetProjectRequestObject) (model.GetProjectResponseObject, error) {
	return model.GetProject500JSONResponse{InternalServerErrorJSONResponse: model.InternalServerErrorJSONResponse(notImplemented())}, nil
}

func (h *Handler) UpdateProject(_ context.Context, _ model.UpdateProjectRequestObject) (model.UpdateProjectResponseObject, error) {
	return model.UpdateProject500JSONResponse{InternalServerErrorJSONResponse: model.InternalServerErrorJSONResponse(notImplemented())}, nil
}

func (h *Handler) DeleteProject(_ context.Context, _ model.DeleteProjectRequestObject) (model.DeleteProjectResponseObject, error) {
	return model.DeleteProject500JSONResponse{InternalServerErrorJSONResponse: model.InternalServerErrorJSONResponse(notImplemented())}, nil
}
