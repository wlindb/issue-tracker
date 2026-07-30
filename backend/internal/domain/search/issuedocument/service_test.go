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

func (m *mockIssueDocumentRepository) Find(ctx context.Context, description string) (issuedocument.IssueDocuments, error) {
	args := m.Called(ctx, description)
	result, _ := args.Get(0).(issuedocument.IssueDocuments)
	return result, args.Error(1)
}

func Test_Create_ValidInput_ReturnsNil(t *testing.T) {
	repository := &mockIssueDocumentRepository{}
	service := issuedocument.NewIssueDocumentService(repository)

	id := uuid.New()
	updatedAt := time.Now()
	expected := issuedocument.NewIssueDocument(id, "test issue", "test description", updatedAt)
	repository.On("Create", mock.Anything, expected).Return(expected, nil)

	err := service.Create(context.Background(), id, "test issue", "test description", updatedAt)
	require.NoError(t, err)
	repository.AssertExpectations(t)
}

func Test_Create_RepositoryError_ReturnsError(t *testing.T) {
	repository := &mockIssueDocumentRepository{}
	service := issuedocument.NewIssueDocumentService(repository)

	id := uuid.New()
	updatedAt := time.Now()
	repoErr := errors.New("create failed")

	repository.On("Create", mock.Anything, issuedocument.NewIssueDocument(id, "test issue", "test description", updatedAt)).
		Return(issuedocument.IssueDocument{}, repoErr)

	err := service.Create(context.Background(), id, "test issue", "test description", updatedAt)
	require.Error(t, err)
	assert.ErrorContains(t, err, "create issue document")
	assert.ErrorIs(t, err, repoErr)
	repository.AssertExpectations(t)
}
