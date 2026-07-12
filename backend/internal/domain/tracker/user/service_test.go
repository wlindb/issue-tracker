//go:build !integration

package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/wlindb/issue-tracker/internal/domain/tracker/user"
)

type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) Upsert(ctx context.Context, u user.User) (user.User, error) {
	args := m.Called(ctx, u)
	return args.Get(0).(user.User), args.Error(1)
}

func Test_Upsert_ValidCommand_ReturnsUser(t *testing.T) {
	repository := &mockUserRepository{}
	service := user.NewUserService(repository)

	id := uuid.New()
	command := user.UpsertUserCommand{ID: id, Email: "jane@example.com", Name: "Jane Doe"}
	expected := user.User{ID: id, Email: "jane@example.com", Name: "Jane Doe"}

	repository.On("Upsert", mock.Anything, mock.MatchedBy(func(u user.User) bool {
		return u.ID == id && u.Email == command.Email && u.Name == command.Name
	})).Return(expected, nil)

	actual, err := service.Upsert(context.Background(), command)
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
	repository.AssertExpectations(t)
}

func Test_Upsert_InvalidCommand_ReturnsError(t *testing.T) {
	repository := &mockUserRepository{}
	service := user.NewUserService(repository)

	command := user.UpsertUserCommand{ID: uuid.New(), Email: ""}

	_, err := service.Upsert(context.Background(), command)

	require.Error(t, err)
	assert.ErrorIs(t, err, user.ErrInvalidUser)
	repository.AssertNotCalled(t, "Upsert")
}

func Test_Upsert_RepositoryError_ReturnsError(t *testing.T) {
	repository := &mockUserRepository{}
	service := user.NewUserService(repository)

	command := user.UpsertUserCommand{ID: uuid.New(), Email: "jane@example.com", Name: "Jane Doe"}
	repoErr := errors.New("db error")

	repository.On("Upsert", mock.Anything, mock.MatchedBy(func(u user.User) bool {
		return u.Email == command.Email
	})).Return(user.User{}, repoErr)

	_, err := service.Upsert(context.Background(), command)
	require.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
	repository.AssertExpectations(t)
}
