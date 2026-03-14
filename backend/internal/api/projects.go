package api

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	openapiTypes "github.com/oapi-codegen/runtime/types"

	"github.com/wlindb/issue-tracker/internal/api/generated"
	trackerdomain "github.com/wlindb/issue-tracker/internal/domain/tracker/project"
)

// ProjectService is what the handler needs from the domain.
type ProjectService interface {
	Create(ctx context.Context, ownerID uuid.UUID, name string, description *string) (*trackerdomain.Project, error)
}

type ProjectHandler struct {
	service ProjectService
}

func NewProjectHandler(service ProjectService) ProjectHandler {
	return ProjectHandler{service: service}
}

func (h *Handler) ListProjects(_ context.Context, _ generated.ListProjectsRequestObject) (generated.ListProjectsResponseObject, error) {
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
	userID := userIDFromContext(ctx)
	if userID == uuid.Nil {
		return nil, fmt.Errorf("missing user ID in context")
	}
	project, err := h.ProjectHandler.service.Create(ctx, userID, req.Body.Name, req.Body.Description)
	if err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}
	return generated.CreateProject201JSONResponse{
		Id:          openapiTypes.UUID(project.ID),
		Name:        project.Name,
		Description: project.Description,
		OwnerId:     openapiTypes.UUID(project.OwnerID),
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
	}, nil
}

func (h *Handler) GetProject(_ context.Context, _ generated.GetProjectRequestObject) (generated.GetProjectResponseObject, error) {
	return generated.GetProject500JSONResponse{InternalServerErrorJSONResponse: generated.InternalServerErrorJSONResponse(notImplemented())}, nil
}

func (h *Handler) UpdateProject(_ context.Context, _ generated.UpdateProjectRequestObject) (generated.UpdateProjectResponseObject, error) {
	return generated.UpdateProject500JSONResponse{InternalServerErrorJSONResponse: generated.InternalServerErrorJSONResponse(notImplemented())}, nil
}

func (h *Handler) DeleteProject(_ context.Context, _ generated.DeleteProjectRequestObject) (generated.DeleteProjectResponseObject, error) {
	return generated.DeleteProject500JSONResponse{InternalServerErrorJSONResponse: generated.InternalServerErrorJSONResponse(notImplemented())}, nil
}
