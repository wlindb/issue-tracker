//go:build !integration

package issuedocument_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/wlindb/issue-tracker/internal/domain/search/issuedocument"
)

type mockIssueDocumentRepository struct {
	mock.Mock
}

func (m *mockIssueDocumentRepository) Create(ctx context.Context, document issuedocument.IssueDocument) (issuedocument.IssueDocument, error) {
	args := m.Called(ctx, document)
	result, _ := args.Get(0).(issuedocument.IssueDocument)
	return result, args.Error(1)
}

func (m *mockIssueDocumentRepository) Update(ctx context.Context, document issuedocument.IssueDocument) (issuedocument.IssueDocument, error) {
	args := m.Called(ctx, document)
	result, _ := args.Get(0).(issuedocument.IssueDocument)
	return result, args.Error(1)
}

func (m *mockIssueDocumentRepository) Find(ctx context.Context, description string, embedding []float32) (issuedocument.IssueDocuments, error) {
	args := m.Called(ctx, description, embedding)
	result, _ := args.Get(0).(issuedocument.IssueDocuments)
	return result, args.Error(1)
}

func (m *mockIssueDocumentRepository) Get(ctx context.Context, id uuid.UUID) (issuedocument.IssueDocument, error) {
	args := m.Called(ctx, id)
	result, _ := args.Get(0).(issuedocument.IssueDocument)
	return result, args.Error(1)
}

type mockEmbeddingGenerator struct {
	mock.Mock
}

func (m *mockEmbeddingGenerator) GenerateEmbedding(ctx context.Context, title, description string) ([]float32, error) {
	args := m.Called(ctx, title, description)
	result, _ := args.Get(0).([]float32)
	return result, args.Error(1)
}

func (m *mockEmbeddingGenerator) GenerateQueryEmbedding(ctx context.Context, query string) ([]float32, error) {
	args := m.Called(ctx, query)
	result, _ := args.Get(0).([]float32)
	return result, args.Error(1)
}

func Test_Create_ValidInput_ReturnsNil(t *testing.T) {
	repository := &mockIssueDocumentRepository{}
	embeddingGenerator := &mockEmbeddingGenerator{}
	service := issuedocument.NewIssueDocumentService(repository, embeddingGenerator)

	id := uuid.New()
	updatedAt := time.Now()
	embedding := []float32{0.1, 0.2, 0.3}
	expected := issuedocument.NewIssueDocument(id, "test issue", "test description", updatedAt, embedding)

	embeddingGenerator.On("GenerateEmbedding", mock.Anything, "test issue", "test description").Return(embedding, nil)
	repository.On("Create", mock.Anything, expected).Return(expected, nil)

	err := service.Create(context.Background(), id, "test issue", "test description", updatedAt)
	require.NoError(t, err)
	repository.AssertExpectations(t)
	embeddingGenerator.AssertExpectations(t)
}

func Test_Create_RepositoryError_ReturnsError(t *testing.T) {
	repository := &mockIssueDocumentRepository{}
	embeddingGenerator := &mockEmbeddingGenerator{}
	service := issuedocument.NewIssueDocumentService(repository, embeddingGenerator)

	id := uuid.New()
	updatedAt := time.Now()
	embedding := []float32{0.1, 0.2, 0.3}
	repoErr := errors.New("create failed")

	embeddingGenerator.On("GenerateEmbedding", mock.Anything, "test issue", "test description").Return(embedding, nil)
	repository.On("Create", mock.Anything, issuedocument.NewIssueDocument(id, "test issue", "test description", updatedAt, embedding)).
		Return(issuedocument.IssueDocument{}, repoErr)

	err := service.Create(context.Background(), id, "test issue", "test description", updatedAt)
	require.Error(t, err)
	assert.ErrorContains(t, err, "create issue document")
	assert.ErrorIs(t, err, repoErr)
	repository.AssertExpectations(t)
	embeddingGenerator.AssertExpectations(t)
}

func Test_Create_EmbeddingError_ReturnsError(t *testing.T) {
	repository := &mockIssueDocumentRepository{}
	embeddingGenerator := &mockEmbeddingGenerator{}
	service := issuedocument.NewIssueDocumentService(repository, embeddingGenerator)

	id := uuid.New()
	updatedAt := time.Now()
	embeddingErr := errors.New("embedding failed")

	embeddingGenerator.On("GenerateEmbedding", mock.Anything, "test issue", "test description").Return(nil, embeddingErr)

	err := service.Create(context.Background(), id, "test issue", "test description", updatedAt)
	require.Error(t, err)
	assert.ErrorContains(t, err, "generate embedding")
	assert.ErrorIs(t, err, embeddingErr)
	repository.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	embeddingGenerator.AssertExpectations(t)
}

func Test_UpdateTitle_ChangedTitle_CallsRepositoryUpdate(t *testing.T) {
	repository := &mockIssueDocumentRepository{}
	service := issuedocument.NewIssueDocumentService(repository, &mockEmbeddingGenerator{})

	id := uuid.New()
	updatedAt := time.Now()
	current := issuedocument.NewIssueDocument(id, "old title", "description", updatedAt, nil)
	expected := current
	expected.Title = "new title"

	repository.On("Get", mock.Anything, id).Return(current, nil)
	repository.On("Update", mock.Anything, expected).Return(expected, nil)

	err := service.UpdateTitle(context.Background(), id, "new title")
	require.NoError(t, err)
	repository.AssertExpectations(t)
}

func Test_UpdateTitle_SameTitle_NoOpDoesNotCallUpdate(t *testing.T) {
	repository := &mockIssueDocumentRepository{}
	service := issuedocument.NewIssueDocumentService(repository, &mockEmbeddingGenerator{})

	id := uuid.New()
	current := issuedocument.NewIssueDocument(id, "same title", "description", time.Now(), nil)
	repository.On("Get", mock.Anything, id).Return(current, nil)

	err := service.UpdateTitle(context.Background(), id, "same title")
	require.NoError(t, err)
	repository.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	repository.AssertExpectations(t)
}

func Test_UpdateTitle_GetRepositoryError_ReturnsWrappedError(t *testing.T) {
	repository := &mockIssueDocumentRepository{}
	service := issuedocument.NewIssueDocumentService(repository, &mockEmbeddingGenerator{})

	id := uuid.New()
	repository.On("Get", mock.Anything, id).Return(issuedocument.IssueDocument{}, issuedocument.ErrIssueDocumentNotFound)

	err := service.UpdateTitle(context.Background(), id, "new title")
	require.Error(t, err)
	assert.ErrorIs(t, err, issuedocument.ErrIssueDocumentNotFound)
	assert.ErrorContains(t, err, "update issue document title")
	repository.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func Test_UpdateTitle_UpdateRepositoryError_ReturnsWrappedError(t *testing.T) {
	repository := &mockIssueDocumentRepository{}
	service := issuedocument.NewIssueDocumentService(repository, &mockEmbeddingGenerator{})

	id := uuid.New()
	current := issuedocument.NewIssueDocument(id, "old title", "description", time.Now(), nil)
	repository.On("Get", mock.Anything, id).Return(current, nil)
	repository.On("Update", mock.Anything, mock.Anything).Return(issuedocument.IssueDocument{}, issuedocument.ErrUpdateConflict)

	err := service.UpdateTitle(context.Background(), id, "new title")
	require.Error(t, err)
	assert.ErrorIs(t, err, issuedocument.ErrUpdateConflict)
}

func Test_UpdateDescription_ChangedDescription_CallsRepositoryUpdate(t *testing.T) {
	repository := &mockIssueDocumentRepository{}
	service := issuedocument.NewIssueDocumentService(repository, &mockEmbeddingGenerator{})

	id := uuid.New()
	updatedAt := time.Now()
	current := issuedocument.NewIssueDocument(id, "title", "old description", updatedAt, nil)
	expected := current
	expected.Description = "new description"

	repository.On("Get", mock.Anything, id).Return(current, nil)
	repository.On("Update", mock.Anything, expected).Return(expected, nil)

	err := service.UpdateDescription(context.Background(), id, "new description")
	require.NoError(t, err)
	repository.AssertExpectations(t)
}

func Test_UpdateDescription_SameDescription_NoOpDoesNotCallUpdate(t *testing.T) {
	repository := &mockIssueDocumentRepository{}
	service := issuedocument.NewIssueDocumentService(repository, &mockEmbeddingGenerator{})

	id := uuid.New()
	current := issuedocument.NewIssueDocument(id, "title", "same description", time.Now(), nil)
	repository.On("Get", mock.Anything, id).Return(current, nil)

	err := service.UpdateDescription(context.Background(), id, "same description")
	require.NoError(t, err)
	repository.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	repository.AssertExpectations(t)
}

func Test_UpdateDescription_GetRepositoryError_ReturnsWrappedError(t *testing.T) {
	repository := &mockIssueDocumentRepository{}
	service := issuedocument.NewIssueDocumentService(repository, &mockEmbeddingGenerator{})

	id := uuid.New()
	repository.On("Get", mock.Anything, id).Return(issuedocument.IssueDocument{}, issuedocument.ErrIssueDocumentNotFound)

	err := service.UpdateDescription(context.Background(), id, "new description")
	require.Error(t, err)
	assert.ErrorIs(t, err, issuedocument.ErrIssueDocumentNotFound)
	assert.ErrorContains(t, err, "update issue document description")
	repository.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func Test_UpdateDescription_UpdateRepositoryError_ReturnsWrappedError(t *testing.T) {
	repository := &mockIssueDocumentRepository{}
	service := issuedocument.NewIssueDocumentService(repository, &mockEmbeddingGenerator{})

	id := uuid.New()
	current := issuedocument.NewIssueDocument(id, "title", "old description", time.Now(), nil)
	repository.On("Get", mock.Anything, id).Return(current, nil)
	repository.On("Update", mock.Anything, mock.Anything).Return(issuedocument.IssueDocument{}, issuedocument.ErrUpdateConflict)

	err := service.UpdateDescription(context.Background(), id, "new description")
	require.Error(t, err)
	assert.ErrorIs(t, err, issuedocument.ErrUpdateConflict)
}

func Test_Find_ValidDescription_ReturnsIssueDocuments(t *testing.T) {
	repository := &mockIssueDocumentRepository{}
	embeddingGenerator := &mockEmbeddingGenerator{}
	service := issuedocument.NewIssueDocumentService(repository, embeddingGenerator)

	description := "database timeout"
	embedding := []float32{0.1, 0.2, 0.3}
	expected := issuedocument.NewIssueDocuments([]issuedocument.IssueDocument{{ID: uuid.New()}, {ID: uuid.New()}})

	embeddingGenerator.On("GenerateQueryEmbedding", mock.Anything, description).Return(embedding, nil)
	repository.On("Find", mock.Anything, description, embedding).Return(expected, nil)

	actual, err := service.Find(context.Background(), description)

	require.NoError(t, err)
	assert.Equal(t, expected, actual)
	embeddingGenerator.AssertExpectations(t)
	repository.AssertExpectations(t)
}

func Test_Find_QueryEmbeddingGenerationError_ReturnsError(t *testing.T) {
	repository := &mockIssueDocumentRepository{}
	embeddingGenerator := &mockEmbeddingGenerator{}
	service := issuedocument.NewIssueDocumentService(repository, embeddingGenerator)

	description := "database timeout"
	embeddingErr := errors.New("query embedding failed")
	embeddingGenerator.On("GenerateQueryEmbedding", mock.Anything, description).Return(nil, embeddingErr)

	_, err := service.Find(context.Background(), description)

	require.Error(t, err)
	assert.ErrorIs(t, err, embeddingErr)
	assert.ErrorContains(t, err, "generate query embedding")
	repository.AssertNotCalled(t, "Find", mock.Anything, mock.Anything, mock.Anything)
	embeddingGenerator.AssertExpectations(t)
}

func Test_Find_RepositoryError_ReturnsError(t *testing.T) {
	repository := &mockIssueDocumentRepository{}
	embeddingGenerator := &mockEmbeddingGenerator{}
	service := issuedocument.NewIssueDocumentService(repository, embeddingGenerator)

	description := "database timeout"
	embedding := []float32{0.1, 0.2, 0.3}
	repositoryErr := errors.New("find failed")

	embeddingGenerator.On("GenerateQueryEmbedding", mock.Anything, description).Return(embedding, nil)
	repository.On("Find", mock.Anything, description, embedding).Return(issuedocument.IssueDocuments{}, repositoryErr)

	_, err := service.Find(context.Background(), description)

	require.Error(t, err)
	assert.ErrorIs(t, err, repositoryErr)
	assert.ErrorContains(t, err, "find issue documents")
	embeddingGenerator.AssertExpectations(t)
	repository.AssertExpectations(t)
}
