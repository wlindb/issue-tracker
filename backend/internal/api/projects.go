package api

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/wlindb/issue-tracker/internal/api/model"
	trackerdomain "github.com/wlindb/issue-tracker/internal/domain/tracker/project"
)

// ProjectService is what the handler needs from the domain.
type ProjectService interface {
	Create(ctx context.Context, id uuid.UUID, ownerID uuid.UUID, name string, description *string) (*trackerdomain.Project, error)
	List(ctx context.Context, ownerID uuid.UUID, cursor *string, limit *int) ([]trackerdomain.Project, *string, error)
}

type ProjectHandler struct {
	service ProjectService
}

func NewProjectHandler(service ProjectService) ProjectHandler {
	return ProjectHandler{service: service}
}

func (h *Handler) ListProjects(ctx context.Context, req model.ListProjectsRequestObject) (model.ListProjectsResponseObject, error) {
	userID := userIDFromContext(ctx)
	if userID == uuid.Nil {
		return model.ListProjects401JSONResponse{
			UnauthorizedJSONResponse: model.UnauthorizedJSONResponse(model.Error{
				Code:    "unauthorized",
				Message: "authentication required",
			}),
		}, nil
	}
	projects, nextCursor, err := h.service.List(ctx, userID, req.Params.Cursor, req.Params.Limit)
	if err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}
	items := make([]model.Project, len(projects))
	for i, p := range projects {
		items[i] = model.Project{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			OwnerId:     p.OwnerID,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		}
	}
	return model.ListProjects200JSONResponse{
		Items:      items,
		NextCursor: nextCursor,
	}, nil
}

func (h *Handler) CreateProject(ctx context.Context, req model.CreateProjectRequestObject) (model.CreateProjectResponseObject, error) {
	if req.Body.Name == "" {
		return model.CreateProject400JSONResponse{
			BadRequestJSONResponse: model.BadRequestJSONResponse(model.Error{
				Code:    "invalid_input",
				Message: "name is required",
			}),
		}, nil
	}
	userID := userIDFromContext(ctx)
	if userID == uuid.Nil {
		return nil, fmt.Errorf("missing user ID in context")
	}
	id := uuid.New()
	project, err := h.service.Create(ctx, id, userID, req.Body.Name, req.Body.Description)
	if err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}
	return model.CreateProject201JSONResponse{
		Id:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		OwnerId:     project.OwnerID,
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
	}, nil
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
