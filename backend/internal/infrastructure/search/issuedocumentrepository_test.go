//go:build !integration

package search

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	issuedocumentdomain "github.com/wlindb/issue-tracker/internal/domain/search/issuedocument"
	searchdb "github.com/wlindb/issue-tracker/internal/infrastructure/search/generated"
)

type mockDBTX struct{ mock.Mock }

func (m *mockDBTX) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	a := m.Called(ctx, sql, args)
	return a.Get(0).(pgx.Row)
}

func (m *mockDBTX) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	a := m.Called(ctx, sql, args)
	if rows, ok := a.Get(0).(pgx.Rows); ok {
		return rows, a.Error(1)
	}
	return nil, a.Error(1)
}

func (m *mockDBTX) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	a := m.Called(ctx, sql, args)
	return a.Get(0).(pgconn.CommandTag), a.Error(1)
}

type mockIssueDocumentRow struct {
	document searchdb.IssueDocument
	err      error
}

func (r *mockIssueDocumentRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}
	*(dest[0].(*uuid.UUID)) = r.document.ID
	*(dest[1].(*uuid.UUID)) = r.document.WorkspaceID
	*(dest[2].(*string)) = r.document.Title
	*(dest[3].(*string)) = r.document.Description
	*(dest[4].(*pgtype.Timestamptz)) = r.document.CreatedAt
	*(dest[5].(*pgtype.Timestamptz)) = r.document.UpdatedAt
	return nil
}

type mockIssueDocumentRows struct {
	documents []searchdb.IssueDocument
	index     int
}

func (r *mockIssueDocumentRows) Close()                                       {}
func (r *mockIssueDocumentRows) Err() error                                   { return nil }
func (r *mockIssueDocumentRows) CommandTag() pgconn.CommandTag                { return pgconn.CommandTag{} }
func (r *mockIssueDocumentRows) FieldDescriptions() []pgconn.FieldDescription { return nil }
func (r *mockIssueDocumentRows) Values() ([]any, error)                       { return nil, nil }
func (r *mockIssueDocumentRows) RawValues() [][]byte                          { return nil }
func (r *mockIssueDocumentRows) Conn() *pgx.Conn                              { return nil }
func (r *mockIssueDocumentRows) Next() bool {
	r.index++
	return r.index <= len(r.documents)
}
func (r *mockIssueDocumentRows) Scan(dest ...any) error {
	return (&mockIssueDocumentRow{document: r.documents[r.index-1]}).Scan(dest...)
}

// — Create unit tests —

func Test_Create_Success_ReturnsDomainIssueDocument(t *testing.T) {
	id := uuid.New()
	workspaceID := uuid.New()
	now := time.Now().UTC()

	returnedRow := searchdb.IssueDocument{
		ID:          id,
		WorkspaceID: workspaceID,
		Title:       "Fix bug",
		Description: "detailed description",
		CreatedAt:   pgtype.Timestamptz{Time: now, Valid: true},
		UpdatedAt:   pgtype.Timestamptz{Time: now, Valid: true},
	}

	mockDatabase := new(mockDBTX)
	mockDatabase.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
		Return(&mockIssueDocumentRow{document: returnedRow})

	repository := &IssueDocumentRepository{db: mockDatabase}

	actual, err := repository.Create(context.Background(), issuedocumentdomain.IssueDocument{
		ID:          id,
		WorkspaceID: workspaceID,
		Title:       "Fix bug",
		Description: "detailed description",
	})

	require.NoError(t, err)
	assert.Equal(t, id, actual.ID)
	assert.Equal(t, "Fix bug", actual.Title)
	assert.Equal(t, "detailed description", actual.Description)
	mockDatabase.AssertExpectations(t)
}

func Test_Create_QueryError_ReturnsWrappedError(t *testing.T) {
	dbErr := errors.New("connection refused")

	mockDatabase := new(mockDBTX)
	mockDatabase.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
		Return(&mockIssueDocumentRow{err: dbErr})

	repository := &IssueDocumentRepository{db: mockDatabase}

	_, err := repository.Create(context.Background(), issuedocumentdomain.IssueDocument{ID: uuid.New()})

	require.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
	assert.Contains(t, err.Error(), "create issue document")
	mockDatabase.AssertExpectations(t)
}

// — Update unit tests —

func Test_Update_Success_ReturnsDomainIssueDocument(t *testing.T) {
	id := uuid.New()
	workspaceID := uuid.New()
	now := time.Now().UTC()

	returnedRow := searchdb.IssueDocument{
		ID:          id,
		WorkspaceID: workspaceID,
		Title:       "Updated title",
		Description: "Updated description",
		CreatedAt:   pgtype.Timestamptz{Time: now, Valid: true},
		UpdatedAt:   pgtype.Timestamptz{Time: now, Valid: true},
	}

	mockDatabase := new(mockDBTX)
	mockDatabase.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
		Return(&mockIssueDocumentRow{document: returnedRow})

	repository := &IssueDocumentRepository{db: mockDatabase}

	actual, err := repository.Update(context.Background(), issuedocumentdomain.IssueDocument{
		ID:          id,
		Title:       "Updated title",
		Description: "Updated description",
	})

	require.NoError(t, err)
	assert.Equal(t, id, actual.ID)
	assert.Equal(t, "Updated title", actual.Title)
	assert.Equal(t, "Updated description", actual.Description)
	mockDatabase.AssertExpectations(t)
}

func Test_Update_QueryError_ReturnsWrappedError(t *testing.T) {
	dbErr := errors.New("connection refused")

	mockDatabase := new(mockDBTX)
	mockDatabase.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
		Return(&mockIssueDocumentRow{err: dbErr})

	repository := &IssueDocumentRepository{db: mockDatabase}

	_, err := repository.Update(context.Background(), issuedocumentdomain.IssueDocument{ID: uuid.New()})

	require.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
	assert.Contains(t, err.Error(), "update issue document")
	mockDatabase.AssertExpectations(t)
}

// — Find unit tests —

func Test_Find_Success_ReturnsDomainIssueDocuments(t *testing.T) {
	now := time.Now().UTC()
	firstID, secondID := uuid.New(), uuid.New()
	workspaceID := uuid.New()

	returnedRows := []searchdb.IssueDocument{
		{ID: firstID, WorkspaceID: workspaceID, Title: "First", Description: "backend issue", CreatedAt: pgtype.Timestamptz{Time: now, Valid: true}, UpdatedAt: pgtype.Timestamptz{Time: now, Valid: true}},
		{ID: secondID, WorkspaceID: workspaceID, Title: "Second", Description: "backend fix", CreatedAt: pgtype.Timestamptz{Time: now, Valid: true}, UpdatedAt: pgtype.Timestamptz{Time: now, Valid: true}},
	}

	mockDatabase := new(mockDBTX)
	mockDatabase.On("Query", mock.Anything, mock.Anything, mock.Anything).
		Return(&mockIssueDocumentRows{documents: returnedRows}, nil)

	repository := &IssueDocumentRepository{db: mockDatabase}

	actual, err := repository.Find(context.Background(), "backend")

	require.NoError(t, err)
	require.Len(t, actual, 2)
	assert.Equal(t, firstID, actual[0].ID)
	assert.Equal(t, secondID, actual[1].ID)
	mockDatabase.AssertExpectations(t)
}

func Test_Find_EmptyResult_ReturnsEmptySlice(t *testing.T) {
	mockDatabase := new(mockDBTX)
	mockDatabase.On("Query", mock.Anything, mock.Anything, mock.Anything).
		Return(&mockIssueDocumentRows{}, nil)

	repository := &IssueDocumentRepository{db: mockDatabase}

	actual, err := repository.Find(context.Background(), "nonexistent")

	require.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Empty(t, actual)
	mockDatabase.AssertExpectations(t)
}

func Test_Find_QueryError_ReturnsWrappedError(t *testing.T) {
	dbErr := errors.New("connection refused")

	mockDatabase := new(mockDBTX)
	mockDatabase.On("Query", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, dbErr)

	repository := &IssueDocumentRepository{db: mockDatabase}

	_, err := repository.Find(context.Background(), "backend")

	require.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
	assert.Contains(t, err.Error(), "find issue documents by description")
	mockDatabase.AssertExpectations(t)
}
