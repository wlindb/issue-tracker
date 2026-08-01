//go:build !integration

package search_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/wlindb/issue-tracker/internal/application/api/model"
	"github.com/wlindb/issue-tracker/internal/application/search"
	"github.com/wlindb/issue-tracker/internal/domain/search/issuedocument"
)

type mockIssueDocumentFinder struct {
	mock.Mock
}

func (m *mockIssueDocumentFinder) Find(ctx context.Context, description string) (issuedocument.IssueDocuments, error) {
	args := m.Called(ctx, description)
	result, _ := args.Get(0).(issuedocument.IssueDocuments)
	return result, args.Error(1)
}

type mockIssueRepository struct {
	mock.Mock
}

func (m *mockIssueRepository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]model.Issue, error) {
	args := m.Called(ctx, ids)
	result, _ := args.Get(0).([]model.Issue)
	return result, args.Error(1)
}

func Test_SearchIssues_UnorderedRepositoryResult_ReturnsIssuesInFindRankOrder(t *testing.T) {
	firstID, secondID, thirdID := uuid.New(), uuid.New(), uuid.New()
	documents := issuedocument.NewIssueDocuments([]issuedocument.IssueDocument{
		{ID: secondID},
		{ID: firstID},
		{ID: thirdID},
	})

	documentFinder := &mockIssueDocumentFinder{}
	issueRepository := &mockIssueRepository{}
	documentFinder.On("Find", mock.Anything, "query").Return(documents, nil)
	issueRepository.On("GetByIDs", mock.Anything, documents.IDs()).Return([]model.Issue{
		{Id: firstID, Title: "First"},
		{Id: thirdID, Title: "Third"},
		{Id: secondID, Title: "Second"},
	}, nil)

	service := search.NewSearcher(documentFinder, issueRepository)

	actual, err := service.SearchIssues(context.Background(), "query")

	require.NoError(t, err)
	require.Len(t, actual, 3)
	assert.Equal(t, secondID, actual[0].Id)
	assert.Equal(t, firstID, actual[1].Id)
	assert.Equal(t, thirdID, actual[2].Id)
}

func Test_SearchIssues_FindRepositoryError_ReturnsWrappedError(t *testing.T) {
	findErr := errors.New("find failed")
	documentFinder := &mockIssueDocumentFinder{}
	issueRepository := &mockIssueRepository{}
	documentFinder.On("Find", mock.Anything, "query").Return(issuedocument.IssueDocuments{}, findErr)

	service := search.NewSearcher(documentFinder, issueRepository)

	_, err := service.SearchIssues(context.Background(), "query")

	require.Error(t, err)
	assert.ErrorIs(t, err, findErr)
	assert.ErrorContains(t, err, "find")
	issueRepository.AssertNotCalled(t, "GetByIDs", mock.Anything, mock.Anything)
}

func Test_SearchIssues_GetByIDsRepositoryError_ReturnsWrappedError(t *testing.T) {
	getErr := errors.New("get by ids failed")
	documents := issuedocument.NewIssueDocuments([]issuedocument.IssueDocument{{ID: uuid.New()}})
	documentFinder := &mockIssueDocumentFinder{}
	issueRepository := &mockIssueRepository{}
	documentFinder.On("Find", mock.Anything, "query").Return(documents, nil)
	issueRepository.On("GetByIDs", mock.Anything, documents.IDs()).Return(nil, getErr)

	service := search.NewSearcher(documentFinder, issueRepository)

	_, err := service.SearchIssues(context.Background(), "query")

	require.Error(t, err)
	assert.ErrorIs(t, err, getErr)
	assert.ErrorContains(t, err, "GetByIDs")
}
