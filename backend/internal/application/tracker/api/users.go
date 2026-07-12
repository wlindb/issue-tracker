package api

import (
	"context"
	"fmt"

	"github.com/wlindb/issue-tracker/internal/application/tracker/api/model"
	userdomain "github.com/wlindb/issue-tracker/internal/domain/tracker/user"
)

// UserService is what the handler needs from the domain.
type UserService interface {
	Upsert(ctx context.Context, command userdomain.UpsertUserCommand) (userdomain.User, error)
}

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) UserHandler {
	return UserHandler{service: service}
}

func (h *UserHandler) UpsertCurrentUser(ctx context.Context, _ model.UpsertCurrentUserRequestObject) (model.UpsertCurrentUserResponseObject, error) {
	userID, err := UserIDFromContext(ctx)
	if err != nil {
		return model.UpsertCurrentUser401JSONResponse{
			UnauthorizedJSONResponse: newUnauthorized("unauthorized", "authentication required"),
		}, nil
	}
	claims, err := UserClaimsFromContext(ctx)
	if err != nil {
		return model.UpsertCurrentUser401JSONResponse{
			UnauthorizedJSONResponse: newUnauthorized("unauthorized", "authentication required"),
		}, nil
	}
	command := userdomain.UpsertUserCommand{
		ID:    userID,
		Email: claims.Email,
		Name:  claims.Name,
	}
	user, err := h.service.Upsert(ctx, command)
	if err != nil {
		return nil, fmt.Errorf("upsert current user: %w", err)
	}
	return model.UpsertCurrentUser200JSONResponse(userFromDomain(user)), nil
}
