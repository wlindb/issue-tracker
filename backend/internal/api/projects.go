package api

import (
	"context"
	"fmt"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/wlindb/issue-tracker/internal/api/generated"
	projectsdomain "github.com/wlindb/issue-tracker/internal/domain/projects"

	"github.com/google/uuid"
)

// ProjectServicer is what the handler needs from the domain.
type ProjectServicer interface {
	Create(ctx context.Context, ownerID uuid.UUID, name string, description *string) (*projectsdomain.Project, error)
}

type ProjectHandler struct {
	service ProjectServicer
}

func NewProjectHandler(service ProjectServicer) ProjectHandler {
	return ProjectHandler{service: service}
}

func (h *Handler) ListProjects(ctx context.Context, request generated.ListProjectsRequestObject) (generated.ListProjectsResponseObject, error) {
	return generated.ListProjects500JSONResponse{InternalServerErrorJSONResponse: generated.InternalServerErrorJSONResponse(notImplemented())}, nil
}

func (h *Handler) CreateProject(ctx context.Context, req generated.CreateProjectRequestObject) (generated.CreateProjectResponseObject, error) {
	if req.Body.Name == "" {
		return generated.CreateProject400JSONResponse{
			BadRequestJSONResponse: generated.BadRequestJSONResponse(generated.Error{
				Code:    "invalid_input",
				Message: "name is required",
			}),
		}, nil
	}
	callerID := callerIDFromContext(ctx)
	project, err := h.ProjectHandler.service.Create(ctx, callerID, req.Body.Name, req.Body.Description)
	if err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}
	return generated.CreateProject201JSONResponse{
		Id:          openapi_types.UUID(project.ID),
		Name:        project.Name,
		Description: project.Description,
		OwnerId:     openapi_types.UUID(project.OwnerID),
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
	}, nil
}

func (h *Handler) GetProject(ctx context.Context, request generated.GetProjectRequestObject) (generated.GetProjectResponseObject, error) {
	return generated.GetProject500JSONResponse{InternalServerErrorJSONResponse: generated.InternalServerErrorJSONResponse(notImplemented())}, nil
}

func (h *Handler) UpdateProject(ctx context.Context, request generated.UpdateProjectRequestObject) (generated.UpdateProjectResponseObject, error) {
	return generated.UpdateProject500JSONResponse{InternalServerErrorJSONResponse: generated.InternalServerErrorJSONResponse(notImplemented())}, nil
}

func (h *Handler) DeleteProject(ctx context.Context, request generated.DeleteProjectRequestObject) (generated.DeleteProjectResponseObject, error) {
	return generated.DeleteProject500JSONResponse{InternalServerErrorJSONResponse: generated.InternalServerErrorJSONResponse(notImplemented())}, nil
}
