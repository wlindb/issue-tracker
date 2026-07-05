package tracker

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	issuedomain "github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
	"github.com/wlindb/issue-tracker/internal/domain/tracker/label"
	trackerdb "github.com/wlindb/issue-tracker/internal/infrastructure/tracker/generated"
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

type mockRow struct {
	issue trackerdb.Issue
	err   error
}

func (r *mockRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}
	*(dest[0].(*uuid.UUID)) = r.issue.ID
	*(dest[1].(*string)) = r.issue.Identifier
	*(dest[2].(*string)) = r.issue.Title
	*(dest[3].(*pgtype.Text)) = r.issue.Description
	*(dest[4].(*string)) = r.issue.Status
	*(dest[5].(*string)) = r.issue.Priority
	*(dest[6].(*pgtype.UUID)) = r.issue.AssigneeID
	*(dest[7].(*uuid.UUID)) = r.issue.ProjectID
	*(dest[8].(*uuid.UUID)) = r.issue.ReporterID
	*(dest[9].(*pgtype.Timestamptz)) = r.issue.CreatedAt
	*(dest[10].(*pgtype.Timestamptz)) = r.issue.UpdatedAt
	*(dest[11].(*uuid.UUID)) = r.issue.WorkspaceID
	return nil
}

type mockRows struct {
	issues []trackerdb.Issue
	index  int
}

func (r *mockRows) Close()                                       {}
func (r *mockRows) Err() error                                   { return nil }
func (r *mockRows) CommandTag() pgconn.CommandTag                { return pgconn.CommandTag{} }
func (r *mockRows) FieldDescriptions() []pgconn.FieldDescription { return nil }
func (r *mockRows) Values() ([]any, error)                       { return nil, nil }
func (r *mockRows) RawValues() [][]byte                          { return nil }
func (r *mockRows) Conn() *pgx.Conn                              { return nil }
func (r *mockRows) Next() bool                                   { r.index++; return r.index <= len(r.issues) }
func (r *mockRows) Scan(dest ...any) error {
	return (&mockRow{issue: r.issues[r.index-1]}).Scan(dest...)
}

// — CreateIssue unit tests —

func Test_CreateIssue_Success_ReturnsDomainIssue(t *testing.T) {
	projectID := uuid.New()
	reporterID := uuid.New()
	labelID := uuid.New()
	now := time.Now().UTC()

	domainIssue := issuedomain.Issue{
		ID:         uuid.New(),
		Identifier: "test-issue-abc",
		Title:      "Test issue",
		Status:     issuedomain.StatusTodo,
		Priority:   issuedomain.PriorityMedium,
		Labels:     []label.Label{{ID: labelID, Name: "backend"}},
		ProjectID:  projectID,
		ReporterID: reporterID,
	}

	returnedRow := trackerdb.Issue{
		ID:         domainIssue.ID,
		Identifier: domainIssue.Identifier,
		Title:      domainIssue.Title,
		Status:     "todo",
		Priority:   "medium",
		ProjectID:  projectID,
		ReporterID: reporterID,
		CreatedAt:  pgtype.Timestamptz{Time: now, Valid: true},
		UpdatedAt:  pgtype.Timestamptz{Time: now, Valid: true},
	}

	mockDatabase := new(mockDBTX)
	mockDatabase.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
		Return(&mockRow{issue: returnedRow})
	mockDatabase.On("Exec", mock.Anything, mock.Anything, mock.Anything).
		Return(pgconn.CommandTag{}, nil)

	repository := &IssueRepository{db: mockDatabase}

	actual, err := repository.CreateIssue(context.Background(), domainIssue)

	require.NoError(t, err)
	assert.Equal(t, domainIssue.ID, actual.ID)
	assert.Equal(t, domainIssue.Title, actual.Title)
	assert.Equal(t, issuedomain.StatusTodo, actual.Status)
	mockDatabase.AssertExpectations(t)
}

func Test_CreateIssue_QueryError_ReturnsWrappedError(t *testing.T) {
	dbErr := errors.New("unique constraint violation")

	mockDatabase := new(mockDBTX)
	mockDatabase.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
		Return(&mockRow{err: dbErr})
	repository := &IssueRepository{db: mockDatabase}

	_, err := repository.CreateIssue(context.Background(), issuedomain.Issue{
		ID:         uuid.New(),
		Identifier: "err-test",
		Title:      "Error test",
		Status:     issuedomain.StatusBacklog,
		Priority:   issuedomain.PriorityNone,
		Labels:     []label.Label{},
		ProjectID:  uuid.New(),
		ReporterID: uuid.New(),
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
	assert.Contains(t, err.Error(), "create issue")
	mockDatabase.AssertExpectations(t)
}

// — ListIssues unit tests —

func Test_ListIssues_Success_ReturnsDomainPage(t *testing.T) {
	projectID := uuid.New()
	now := time.Now().UTC()

	returnedRows := []trackerdb.Issue{
		{
			ID:         uuid.New(),
			Identifier: "issue-1",
			Title:      "First issue",
			Status:     "backlog",
			Priority:   "none",
			ProjectID:  projectID,
			ReporterID: uuid.New(),
			CreatedAt:  pgtype.Timestamptz{Time: now, Valid: true},
			UpdatedAt:  pgtype.Timestamptz{Time: now, Valid: true},
		},
	}

	mockDatabase := new(mockDBTX)
	mockDatabase.On("Query", mock.Anything, mock.Anything, mock.Anything).
		Return(&mockRows{issues: returnedRows}, nil)
	repository := &IssueRepository{db: mockDatabase}

	actual, err := repository.ListIssues(context.Background(), projectID, issuedomain.ListIssueQuery{})

	require.NoError(t, err)
	require.Len(t, actual.Items, 1)
	assert.Equal(t, "First issue", actual.Items[0].Title)
	assert.Equal(t, issuedomain.StatusBacklog, actual.Items[0].Status)
	mockDatabase.AssertExpectations(t)
}

func Test_ListIssues_EmptyResult_ReturnsEmptyPage(t *testing.T) {
	mockDatabase := new(mockDBTX)
	mockDatabase.On("Query", mock.Anything, mock.Anything, mock.Anything).
		Return(&mockRows{}, nil)
	repository := &IssueRepository{db: mockDatabase}

	actual, err := repository.ListIssues(context.Background(), uuid.New(), issuedomain.ListIssueQuery{})

	require.NoError(t, err)
	assert.Empty(t, actual.Items)
	mockDatabase.AssertExpectations(t)
}

func Test_ListIssues_QueryError_ReturnsWrappedError(t *testing.T) {
	dbErr := errors.New("connection refused")

	mockDatabase := new(mockDBTX)
	mockDatabase.On("Query", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, dbErr)
	repository := &IssueRepository{db: mockDatabase}

	_, err := repository.ListIssues(context.Background(), uuid.New(), issuedomain.ListIssueQuery{})

	require.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
	assert.Contains(t, err.Error(), "list issues")
	mockDatabase.AssertExpectations(t)
}

func Test_ListIssues_WithFilters_PassesCorrectParams(t *testing.T) {
	projectID := uuid.New()
	assigneeID := uuid.New()
	status := issuedomain.StatusInProgress
	priority := issuedomain.PriorityHigh

	query := issuedomain.ListIssueQuery{
		Status:     &status,
		Priority:   &priority,
		AssigneeID: &assigneeID,
	}
	expectedParams := listIssuesParamsFromDomain(projectID, query)

	var capturedArgs []any
	mockDatabase := new(mockDBTX)
	mockDatabase.On("Query", mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			capturedArgs = args.Get(2).([]any)
		}).
		Return(&mockRows{}, nil)
	repository := &IssueRepository{db: mockDatabase}

	_, err := repository.ListIssues(context.Background(), projectID, query)

	require.NoError(t, err)
	require.Len(t, capturedArgs, 4)
	assert.Equal(t, expectedParams.ProjectID, capturedArgs[0])
	assert.Equal(t, expectedParams.Status, capturedArgs[1])
	assert.Equal(t, expectedParams.Priority, capturedArgs[2])
	assert.Equal(t, expectedParams.AssigneeID, capturedArgs[3])
	mockDatabase.AssertExpectations(t)
}

// — Update unit tests —

func Test_Update_Success_ReturnsDomainIssue(t *testing.T) {
	projectID := uuid.New()
	reporterID := uuid.New()
	description := "updated desc"
	now := time.Now().UTC()

	domainIssue := issuedomain.Issue{
		ID:          uuid.New(),
		Identifier:  "test-issue-abc",
		Title:       "Test issue",
		Description: &description,
		Status:      issuedomain.StatusInProgress,
		Priority:    issuedomain.PriorityHigh,
		Labels:      []label.Label{},
		ProjectID:   projectID,
		ReporterID:  reporterID,
	}

	returnedRow := trackerdb.Issue{
		ID:          domainIssue.ID,
		Identifier:  domainIssue.Identifier,
		Title:       domainIssue.Title,
		Description: pgtype.Text{String: description, Valid: true},
		Status:      "in_progress",
		Priority:    "high",
		ProjectID:   projectID,
		ReporterID:  reporterID,
		CreatedAt:   pgtype.Timestamptz{Time: now, Valid: true},
		UpdatedAt:   pgtype.Timestamptz{Time: now, Valid: true},
	}

	mockDatabase := new(mockDBTX)
	mockDatabase.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
		Return(&mockRow{issue: returnedRow})
	repository := &IssueRepository{db: mockDatabase}

	actual, err := repository.Update(context.Background(), domainIssue)

	require.NoError(t, err)
	assert.Equal(t, domainIssue.ID, actual.ID)
	assert.Equal(t, issuedomain.StatusInProgress, actual.Status)
	assert.Equal(t, issuedomain.PriorityHigh, actual.Priority)
	require.NotNil(t, actual.Description)
	assert.Equal(t, description, *actual.Description)
	mockDatabase.AssertExpectations(t)
}

func Test_Update_QueryError_ReturnsWrappedError(t *testing.T) {
	dbErr := errors.New("update conflict")

	mockDatabase := new(mockDBTX)
	mockDatabase.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
		Return(&mockRow{err: dbErr})
	repository := &IssueRepository{db: mockDatabase}

	_, err := repository.Update(context.Background(), issuedomain.Issue{
		ID:       uuid.New(),
		Status:   issuedomain.StatusDone,
		Priority: issuedomain.PriorityNone,
		Labels:   []label.Label{},
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
	assert.Contains(t, err.Error(), "update issue")
	mockDatabase.AssertExpectations(t)
}

// — AddLabel unit tests —

func Test_AddLabel_Success_ReturnsNil(t *testing.T) {
	mockDatabase := new(mockDBTX)
	mockDatabase.On("Exec", mock.Anything, mock.Anything, mock.Anything).
		Return(pgconn.CommandTag{}, nil)
	repository := &IssueRepository{db: mockDatabase}

	err := repository.AddLabel(context.Background(), uuid.New(), label.Label{ID: uuid.New()})

	require.NoError(t, err)
	mockDatabase.AssertExpectations(t)
}

func Test_AddLabel_LabelForeignKeyViolation_ReturnsErrLabelNotFound(t *testing.T) {
	pgErr := &pgconn.PgError{Code: pgerrcode.ForeignKeyViolation, ConstraintName: "issue_labels_label_id_fkey"}

	mockDatabase := new(mockDBTX)
	mockDatabase.On("Exec", mock.Anything, mock.Anything, mock.Anything).
		Return(pgconn.CommandTag{}, pgErr)
	repository := &IssueRepository{db: mockDatabase}

	err := repository.AddLabel(context.Background(), uuid.New(), label.Label{ID: uuid.New()})

	require.Error(t, err)
	assert.ErrorIs(t, err, label.ErrLabelNotFound)
	mockDatabase.AssertExpectations(t)
}

func Test_AddLabel_IssueForeignKeyViolation_ReturnsErrIssueNotFound(t *testing.T) {
	pgErr := &pgconn.PgError{Code: pgerrcode.ForeignKeyViolation, ConstraintName: "issue_labels_issue_id_fkey"}

	mockDatabase := new(mockDBTX)
	mockDatabase.On("Exec", mock.Anything, mock.Anything, mock.Anything).
		Return(pgconn.CommandTag{}, pgErr)
	repository := &IssueRepository{db: mockDatabase}

	err := repository.AddLabel(context.Background(), uuid.New(), label.Label{ID: uuid.New()})

	require.Error(t, err)
	assert.ErrorIs(t, err, issuedomain.ErrIssueNotFound)
	mockDatabase.AssertExpectations(t)
}

func Test_AddLabel_LabelWorkspaceMismatch_ReturnsErrLabelNotFound(t *testing.T) {
	pgErr := &pgconn.PgError{Code: pgerrcode.ForeignKeyViolation, ConstraintName: "issue_labels_workspace_matches_label"}

	mockDatabase := new(mockDBTX)
	mockDatabase.On("Exec", mock.Anything, mock.Anything, mock.Anything).
		Return(pgconn.CommandTag{}, pgErr)
	repository := &IssueRepository{db: mockDatabase}

	err := repository.AddLabel(context.Background(), uuid.New(), label.Label{ID: uuid.New()})

	require.Error(t, err)
	assert.ErrorIs(t, err, label.ErrLabelNotFound)
	mockDatabase.AssertExpectations(t)
}

func Test_AddLabel_IssueWorkspaceMismatch_ReturnsErrIssueNotFound(t *testing.T) {
	pgErr := &pgconn.PgError{Code: pgerrcode.ForeignKeyViolation, ConstraintName: "issue_labels_workspace_matches_issue"}

	mockDatabase := new(mockDBTX)
	mockDatabase.On("Exec", mock.Anything, mock.Anything, mock.Anything).
		Return(pgconn.CommandTag{}, pgErr)
	repository := &IssueRepository{db: mockDatabase}

	err := repository.AddLabel(context.Background(), uuid.New(), label.Label{ID: uuid.New()})

	require.Error(t, err)
	assert.ErrorIs(t, err, issuedomain.ErrIssueNotFound)
	mockDatabase.AssertExpectations(t)
}

func Test_AddLabel_OtherError_ReturnsWrappedError(t *testing.T) {
	dbErr := errors.New("db down")

	mockDatabase := new(mockDBTX)
	mockDatabase.On("Exec", mock.Anything, mock.Anything, mock.Anything).
		Return(pgconn.CommandTag{}, dbErr)
	repository := &IssueRepository{db: mockDatabase}

	err := repository.AddLabel(context.Background(), uuid.New(), label.Label{ID: uuid.New()})

	require.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
	assert.Contains(t, err.Error(), "add label")
	mockDatabase.AssertExpectations(t)
}
