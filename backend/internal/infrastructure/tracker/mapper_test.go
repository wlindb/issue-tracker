//go:build !integration

package tracker

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	issuedomain "github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
	trackerdb "github.com/wlindb/issue-tracker/internal/infrastructure/tracker/generated"
)

func Test_ProjectToDomain_NilDescription_SetsDescriptionNil(t *testing.T) {
	id, ownerID := uuid.New(), uuid.New()
	now := time.Now().UTC()
	row := trackerdb.Project{
		ID:          id,
		OwnerID:     ownerID,
		Name:        "Test",
		Description: pgtype.Text{Valid: false},
		CreatedAt:   pgtype.Timestamptz{Time: now, Valid: true},
		UpdatedAt:   pgtype.Timestamptz{Time: now, Valid: true},
	}

	actual := projectToDomain(row)

	require.NotNil(t, actual)
	assert.Equal(t, id, actual.ID)
	assert.Equal(t, ownerID, actual.OwnerID)
	assert.Equal(t, "Test", actual.Name)
	assert.Nil(t, actual.Description)
	assert.Equal(t, now, actual.CreatedAt)
}

func Test_ProjectToDomain_WithDescription_SetsDescription(t *testing.T) {
	description := "hello"
	now := time.Now().UTC()
	row := trackerdb.Project{
		ID:          uuid.New(),
		OwnerID:     uuid.New(),
		Name:        "Test",
		Description: pgtype.Text{String: description, Valid: true},
		CreatedAt:   pgtype.Timestamptz{Time: now, Valid: true},
		UpdatedAt:   pgtype.Timestamptz{Time: now, Valid: true},
	}

	actual := projectToDomain(row)

	require.NotNil(t, actual.Description)
	assert.Equal(t, description, *actual.Description)
}

func Test_ProjectsToDomain_Empty_ReturnsEmptySlice(t *testing.T) {
	actual := projectsToDomain([]trackerdb.Project{})

	assert.NotNil(t, actual)
	assert.Empty(t, actual)
}

func Test_ProjectsToDomain_MultipleRows_ReturnsMappedProjects(t *testing.T) {
	now := time.Now().UTC()
	id1, id2 := uuid.New(), uuid.New()
	rows := []trackerdb.Project{
		{
			ID:        id1,
			OwnerID:   uuid.New(),
			Name:      "Alpha",
			CreatedAt: pgtype.Timestamptz{Time: now, Valid: true},
			UpdatedAt: pgtype.Timestamptz{Time: now, Valid: true},
		},
		{
			ID:        id2,
			OwnerID:   uuid.New(),
			Name:      "Beta",
			CreatedAt: pgtype.Timestamptz{Time: now, Valid: true},
			UpdatedAt: pgtype.Timestamptz{Time: now, Valid: true},
		},
	}

	actual := projectsToDomain(rows)

	require.Len(t, actual, 2)
	assert.Equal(t, id1, actual[0].ID)
	assert.Equal(t, "Alpha", actual[0].Name)
	assert.Equal(t, id2, actual[1].ID)
	assert.Equal(t, "Beta", actual[1].Name)
}

// — Issue mapper tests —

func Test_IssueToDomain_NilOptionalFields_SetsNils(t *testing.T) {
	issueID := uuid.New()
	projectID := uuid.New()
	reporterID := uuid.New()
	now := time.Now().UTC()
	row := trackerdb.Issue{
		ID:          issueID,
		Identifier:  "fix-bug-abc123",
		Title:       "Fix bug",
		Description: pgtype.Text{Valid: false},
		Status:      "backlog",
		Priority:    "none",
		Labels:      []string{},
		AssigneeID:  pgtype.UUID{Valid: false},
		ProjectID:   projectID,
		ReporterID:  reporterID,
		CreatedAt:   pgtype.Timestamptz{Time: now, Valid: true},
		UpdatedAt:   pgtype.Timestamptz{Time: now, Valid: true},
	}

	actual := issueToDomain(row)

	require.NotNil(t, actual)
	assert.Equal(t, issueID, actual.ID)
	assert.Equal(t, "fix-bug-abc123", actual.Identifier)
	assert.Equal(t, "Fix bug", actual.Title)
	assert.Nil(t, actual.Description)
	assert.Equal(t, issuedomain.StatusBacklog, actual.Status)
	assert.Equal(t, issuedomain.PriorityNone, actual.Priority)
	assert.Empty(t, actual.Labels)
	assert.Nil(t, actual.AssigneeID)
	assert.Equal(t, projectID, actual.ProjectID)
	assert.Equal(t, reporterID, actual.ReporterID)
	assert.Equal(t, now, actual.CreatedAt)
	assert.Equal(t, now, actual.UpdatedAt)
}

func Test_IssueToDomain_WithOptionalFields_SetsValues(t *testing.T) {
	assigneeID := uuid.New()
	description := "detailed description"
	now := time.Now().UTC()
	row := trackerdb.Issue{
		ID:          uuid.New(),
		Identifier:  "new-feat-abc123",
		Title:       "New feature",
		Description: pgtype.Text{String: description, Valid: true},
		Status:      "todo",
		Priority:    "high",
		Labels:      []string{"frontend", "urgent"},
		AssigneeID:  pgtype.UUID{Bytes: assigneeID, Valid: true},
		ProjectID:   uuid.New(),
		ReporterID:  uuid.New(),
		CreatedAt:   pgtype.Timestamptz{Time: now, Valid: true},
		UpdatedAt:   pgtype.Timestamptz{Time: now, Valid: true},
	}

	actual := issueToDomain(row)

	require.NotNil(t, actual.Description)
	assert.Equal(t, description, *actual.Description)
	require.NotNil(t, actual.AssigneeID)
	assert.Equal(t, assigneeID, *actual.AssigneeID)
	assert.Equal(t, issuedomain.StatusTodo, actual.Status)
	assert.Equal(t, issuedomain.PriorityHigh, actual.Priority)
	assert.Equal(t, []string{"frontend", "urgent"}, actual.Labels)
}

func Test_IssuesToDomain_Empty_ReturnsEmptySlice(t *testing.T) {
	actual := issuesToDomain([]trackerdb.Issue{})

	assert.NotNil(t, actual)
	assert.Empty(t, actual)
}

func Test_IssuesToDomain_MultipleRows_ReturnsMappedIssues(t *testing.T) {
	now := time.Now().UTC()
	firstID, secondID := uuid.New(), uuid.New()
	rows := []trackerdb.Issue{
		{
			ID:         firstID,
			Identifier: "issue-a",
			Title:      "Issue A",
			Status:     "backlog",
			Priority:   "low",
			Labels:     []string{},
			ProjectID:  uuid.New(),
			ReporterID: uuid.New(),
			CreatedAt:  pgtype.Timestamptz{Time: now, Valid: true},
			UpdatedAt:  pgtype.Timestamptz{Time: now, Valid: true},
		},
		{
			ID:         secondID,
			Identifier: "issue-b",
			Title:      "Issue B",
			Status:     "done",
			Priority:   "urgent",
			Labels:     []string{"backend"},
			ProjectID:  uuid.New(),
			ReporterID: uuid.New(),
			CreatedAt:  pgtype.Timestamptz{Time: now, Valid: true},
			UpdatedAt:  pgtype.Timestamptz{Time: now, Valid: true},
		},
	}

	actual := issuesToDomain(rows)

	require.Len(t, actual, 2)
	assert.Equal(t, firstID, actual[0].ID)
	assert.Equal(t, "Issue A", actual[0].Title)
	assert.Equal(t, secondID, actual[1].ID)
	assert.Equal(t, "Issue B", actual[1].Title)
}

// — CreateIssueParamsFromDomain tests —

func Test_CreateIssueParamsFromDomain_NoOptionalFields_MapsCorrectly(t *testing.T) {
	issueID := uuid.New()
	projectID := uuid.New()
	reporterID := uuid.New()
	domainIssue := issuedomain.Issue{
		ID:         issueID,
		Identifier: "test-issue-abc",
		Title:      "Test issue",
		Status:     issuedomain.StatusBacklog,
		Priority:   issuedomain.PriorityNone,
		Labels:     []string{},
		ProjectID:  projectID,
		ReporterID: reporterID,
	}

	actual := createIssueParamsFromDomain(domainIssue)

	assert.Equal(t, issueID, actual.ID)
	assert.Equal(t, "test-issue-abc", actual.Identifier)
	assert.Equal(t, "Test issue", actual.Title)
	assert.False(t, actual.Description.Valid)
	assert.Equal(t, "backlog", actual.Status)
	assert.Equal(t, "none", actual.Priority)
	assert.Equal(t, []string{}, actual.Labels)
	assert.False(t, actual.AssigneeID.Valid)
	assert.Equal(t, projectID, actual.ProjectID)
	assert.Equal(t, reporterID, actual.ReporterID)
}

func Test_CreateIssueParamsFromDomain_WithOptionalFields_MapsCorrectly(t *testing.T) {
	issueID := uuid.New()
	assigneeID := uuid.New()
	description := "some description"
	domainIssue := issuedomain.Issue{
		ID:          issueID,
		Identifier:  "full-issue-abc",
		Title:       "Full issue",
		Description: &description,
		Status:      issuedomain.StatusInProgress,
		Priority:    issuedomain.PriorityHigh,
		Labels:      []string{"backend", "urgent"},
		AssigneeID:  &assigneeID,
		ProjectID:   uuid.New(),
		ReporterID:  uuid.New(),
	}

	actual := createIssueParamsFromDomain(domainIssue)

	assert.True(t, actual.Description.Valid)
	assert.Equal(t, description, actual.Description.String)
	assert.True(t, actual.AssigneeID.Valid)
	assert.Equal(t, assigneeID, uuid.UUID(actual.AssigneeID.Bytes))
	assert.Equal(t, "in_progress", actual.Status)
	assert.Equal(t, "high", actual.Priority)
	assert.Equal(t, []string{"backend", "urgent"}, actual.Labels)
}

// — ListIssuesParamsFromDomain tests —

func Test_ListIssuesParamsFromDomain_NoFilters_MapsCorrectly(t *testing.T) {
	projectID := uuid.New()
	query := issuedomain.ListIssueQuery{}

	actual := listIssuesParamsFromDomain(projectID, query)

	assert.Equal(t, projectID, actual.ProjectID)
	assert.False(t, actual.Status.Valid)
	assert.False(t, actual.Priority.Valid)
	assert.False(t, actual.AssigneeID.Valid)
}

func Test_ListIssuesParamsFromDomain_AllFilters_MapsCorrectly(t *testing.T) {
	projectID := uuid.New()
	assigneeID := uuid.New()
	status := issuedomain.StatusTodo
	priority := issuedomain.PriorityHigh
	query := issuedomain.ListIssueQuery{
		Status:     &status,
		Priority:   &priority,
		AssigneeID: &assigneeID,
	}

	actual := listIssuesParamsFromDomain(projectID, query)

	assert.Equal(t, projectID, actual.ProjectID)
	assert.True(t, actual.Status.Valid)
	assert.Equal(t, "todo", actual.Status.String)
	assert.True(t, actual.Priority.Valid)
	assert.Equal(t, "high", actual.Priority.String)
	assert.True(t, actual.AssigneeID.Valid)
	assert.Equal(t, assigneeID, uuid.UUID(actual.AssigneeID.Bytes))
}
