package api

import (
	"context"

	"github.com/wlindb/issue-tracker/internal/api/generated"
)

func (h *Handler) Login(ctx context.Context, request generated.LoginRequestObject) (generated.LoginResponseObject, error) {
	return generated.Login500JSONResponse{InternalServerErrorJSONResponse: generated.InternalServerErrorJSONResponse(notImplemented())}, nil
}

func (h *Handler) Register(ctx context.Context, request generated.RegisterRequestObject) (generated.RegisterResponseObject, error) {
	return generated.Register500JSONResponse{InternalServerErrorJSONResponse: generated.InternalServerErrorJSONResponse(notImplemented())}, nil
}
