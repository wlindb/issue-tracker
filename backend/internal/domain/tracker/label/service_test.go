//go:build !integration

package label_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/wlindb/issue-tracker/internal/domain/tracker/label"
)

type mockLabelRepository struct {
	mock.Mock
}

func (m *mockLabelRepository) GetOrCreate(ctx context.Context, name string) (label.Label, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(label.Label), args.Error(1)
}

func (m *mockLabelRepository) ListByIDs(ctx context.Context, ids []uuid.UUID) ([]label.Label, error) {
	args := m.Called(ctx, ids)
	result, _ := args.Get(0).([]label.Label)
	return result, args.Error(1)
}

func (m *mockLabelRepository) SearchByName(ctx context.Context, name string) ([]label.Label, error) {
	args := m.Called(ctx, name)
	result, _ := args.Get(0).([]label.Label)
	return result, args.Error(1)
}

func Test_Create_ValidName_ReturnsLabel(t *testing.T) {
	repository := &mockLabelRepository{}
	service := label.NewLabelService(repository)

	name := "bug"
	expected := label.Label{ID: uuid.New(), Name: name}

	repository.On("GetOrCreate", mock.Anything, name).Return(expected, nil)

	actual, err := service.Create(context.Background(), name)
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
	repository.AssertExpectations(t)
}

func Test_Create_RepositoryError_ReturnsWrappedError(t *testing.T) {
	repository := &mockLabelRepository{}
	service := label.NewLabelService(repository)

	repoErr := errors.New("db error")
	repository.On("GetOrCreate", mock.Anything, "bug").Return(label.Label{}, repoErr)

	_, err := service.Create(context.Background(), "bug")
	require.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
	repository.AssertExpectations(t)
}

func Test_Search_ValidName_ReturnsLabels(t *testing.T) {
	repository := &mockLabelRepository{}
	service := label.NewLabelService(repository)

	name := "bug"
	expected := []label.Label{
		{ID: uuid.New(), Name: "bug"},
		{ID: uuid.New(), Name: "bugfix"},
	}

	repository.On("SearchByName", mock.Anything, name).Return(expected, nil)

	actual, err := service.Search(context.Background(), name)
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
	repository.AssertExpectations(t)
}

func Test_Search_RepositoryError_ReturnsWrappedError(t *testing.T) {
	repository := &mockLabelRepository{}
	service := label.NewLabelService(repository)

	repoErr := errors.New("db error")
	repository.On("SearchByName", mock.Anything, "bug").Return([]label.Label(nil), repoErr)

	_, err := service.Search(context.Background(), "bug")
	require.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
	repository.AssertExpectations(t)
}
