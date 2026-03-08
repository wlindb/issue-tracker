package api

import (
	"context"
	"errors"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/wlindb/issue-tracker/internal/api/generated"
	"github.com/wlindb/issue-tracker/internal/domain"
)

type AuthHandler struct {
	svc *domain.Service
}

func NewAuthHandler(svc *domain.Service) AuthHandler {
	return AuthHandler{svc: svc}
}

func (h AuthHandler) Login(ctx context.Context, request generated.LoginRequestObject) (generated.LoginResponseObject, error) {
	user, token, err := h.svc.Login(ctx, string(request.Body.Email), request.Body.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			return generated.Login401JSONResponse{
				UnauthorizedJSONResponse: generated.UnauthorizedJSONResponse{
					Code:    "invalid_credentials",
					Message: "invalid email or password",
				},
			}, nil
		}
		return generated.Login500JSONResponse{
			InternalServerErrorJSONResponse: generated.InternalServerErrorJSONResponse{
				Code:    "internal_error",
				Message: "an unexpected error occurred",
			},
		}, nil
	}

	return generated.Login200JSONResponse{
		Token: token,
		User: generated.User{
			Id:        openapi_types.UUID(user.ID),
			Email:     openapi_types.Email(user.Email),
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
}

func (h AuthHandler) Register(ctx context.Context, request generated.RegisterRequestObject) (generated.RegisterResponseObject, error) {
	user, token, err := h.svc.Register(ctx, string(request.Body.Email), request.Body.Name, request.Body.Password)
	if err != nil {
		if errors.Is(err, domain.ErrEmailTaken) {
			return generated.Register422JSONResponse{
				UnprocessableEntityJSONResponse: generated.UnprocessableEntityJSONResponse{
					Code:    "email_taken",
					Message: "an account with this email already exists",
				},
			}, nil
		}
		return generated.Register500JSONResponse{
			InternalServerErrorJSONResponse: generated.InternalServerErrorJSONResponse{
				Code:    "internal_error",
				Message: "an unexpected error occurred",
			},
		}, nil
	}

	return generated.Register201JSONResponse{
		Token: token,
		User: generated.User{
			Id:        openapi_types.UUID(user.ID),
			Email:     openapi_types.Email(user.Email),
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
}
