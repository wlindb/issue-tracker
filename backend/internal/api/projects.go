package api

import (
	"context"

	"github.com/wlindb/issue-tracker/internal/api/generated"
)

func (h *Handler) ListProjects(ctx context.Context, request generated.ListProjectsRequestObject) (generated.ListProjectsResponseObject, error) {
	return generated.ListProjects500JSONResponse{InternalServerErrorJSONResponse: generated.InternalServerErrorJSONResponse(notImplemented())}, nil
}

func (h *Handler) CreateProject(ctx context.Context, request generated.CreateProjectRequestObject) (generated.CreateProjectResponseObject, error) {
	return generated.CreateProject500JSONResponse{InternalServerErrorJSONResponse: generated.InternalServerErrorJSONResponse(notImplemented())}, nil
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
