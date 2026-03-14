package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/wlindb/issue-tracker/internal/api/model"
	authdomain "github.com/wlindb/issue-tracker/internal/domain/auth"
)

type AuthHandler struct {
	service *authdomain.AuthService
}

func NewAuthHandler(service *authdomain.AuthService) AuthHandler {
	return AuthHandler{service: service}
}

func (h AuthHandler) Login(ctx context.Context, request model.LoginRequestObject) (model.LoginResponseObject, error) {
	user, token, err := h.service.Login(ctx, request.Body.Email, request.Body.Password)
	if err != nil {
		if errors.Is(err, authdomain.ErrInvalidCredentials) {
			return model.Login401JSONResponse{
				UnauthorizedJSONResponse: model.UnauthorizedJSONResponse{
					Code:    "invalid_credentials",
					Message: "invalid email or password",
				},
			}, nil
		}
		return nil, fmt.Errorf("login: %w", err)
	}

	return model.Login200JSONResponse{
		Token: token,
		User: model.User{
			Id:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
}

func (h AuthHandler) Register(ctx context.Context, request model.RegisterRequestObject) (model.RegisterResponseObject, error) {
	user, token, err := h.service.Register(ctx, request.Body.Email, request.Body.Name, request.Body.Password)
	if err != nil {
		if errors.Is(err, authdomain.ErrEmailTaken) {
			return model.Register422JSONResponse{
				UnprocessableEntityJSONResponse: model.UnprocessableEntityJSONResponse{
					Code:    "email_taken",
					Message: "an account with this email already exists",
				},
			}, nil
		}
		return nil, fmt.Errorf("register: %w", err)
	}

	return model.Register201JSONResponse{
		Token: token,
		User: model.User{
			Id:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
}
