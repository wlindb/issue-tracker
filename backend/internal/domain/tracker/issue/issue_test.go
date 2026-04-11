//go:build !integration

package issue_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
)

func Test_ToIssue_NewCommand_GeneratesNonZeroID(t *testing.T) {
	command := issue.CreateIssueCommand{
		ProjectID:  uuid.New(),
		ReporterID: uuid.New(),
		Title:      "New feature",
		Status:     issue.StatusTodo,
		Priority:   issue.PriorityMedium,
	}

	result := command.ToIssue(uuid.New(), command.Slugify)

	assert.NotEqual(t, uuid.Nil, result.ID)
}

func Test_ToIssue_Title_GeneratesSlugIdentifier(t *testing.T) {
	fixedID := uuid.MustParse("a1b2c3d4-0000-0000-0000-000000000000")
	idSuffix := fixedID.String()[:8] // "a1b2c3d4"

	tests := []struct {
		title      string
		identifier string
	}{
		{"New feature", "new-feature-" + idSuffix},
		{"Fix login bug", "fix-login-bug-" + idSuffix},
		{"  leading and trailing  ", "leading-and-trailing-" + idSuffix},
		{"Special chars: foo@bar!", "special-chars-foobar-" + idSuffix},
		{"Multiple   spaces", "multiple-spaces-" + idSuffix},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			command := issue.CreateIssueCommand{
				ProjectID:  uuid.New(),
				ReporterID: uuid.New(),
				Title:      tt.title,
				Status:     issue.StatusTodo,
				Priority:   issue.PriorityMedium,
			}

			result := command.ToIssue(fixedID, command.Slugify)

			assert.Equal(t, tt.identifier, result.Identifier)
		})
	}
}

func Test_ToIssue_FullCommand_PopulatesAllFields(t *testing.T) {
	projectID := uuid.New()
	reporterID := uuid.New()
	assigneeID := uuid.New()
	description := "some description"
	command := issue.CreateIssueCommand{
		ProjectID:   projectID,
		ReporterID:  reporterID,
		Title:       "New feature",
		Description: &description,
		Status:      issue.StatusTodo,
		Priority:    issue.PriorityMedium,
		AssigneeID:  &assigneeID,
	}

	result := command.ToIssue(uuid.New(), command.Slugify)

	assert.Equal(t, projectID, result.ProjectID)
	assert.Equal(t, reporterID, result.ReporterID)
	assert.Equal(t, "New feature", result.Title)
	assert.Equal(t, &description, result.Description)
	assert.Equal(t, issue.StatusTodo, result.Status)
	assert.Equal(t, issue.PriorityMedium, result.Priority)
	assert.Equal(t, &assigneeID, result.AssigneeID)
}

func Test_ToIssue_NewIssue_HasEmptyLabels(t *testing.T) {
	command := issue.CreateIssueCommand{
		ProjectID:  uuid.New(),
		ReporterID: uuid.New(),
		Title:      "New feature",
		Status:     issue.StatusTodo,
		Priority:   issue.PriorityMedium,
	}

	result := command.ToIssue(uuid.New(), command.Slugify)

	assert.NotNil(t, result.Labels)
	assert.Empty(t, result.Labels)
}

// — Status.Valid —

func Test_Status_Valid_KnownValues_ReturnTrue(t *testing.T) {
	statuses := []issue.Status{
		issue.StatusBacklog,
		issue.StatusTodo,
		issue.StatusInProgress,
		issue.StatusDone,
		issue.StatusCancelled,
	}
	for _, s := range statuses {
		assert.True(t, s.Valid(), "expected %q to be valid", s)
	}
}

func Test_Status_Valid_UnknownValue_ReturnFalse(t *testing.T) {
	assert.False(t, issue.Status("unknown").Valid())
}

// — Priority.Valid —

func Test_Priority_Valid_KnownValues_ReturnTrue(t *testing.T) {
	priorities := []issue.Priority{
		issue.PriorityNone,
		issue.PriorityLow,
		issue.PriorityMedium,
		issue.PriorityHigh,
		issue.PriorityUrgent,
	}
	for _, p := range priorities {
		assert.True(t, p.Valid(), "expected %q to be valid", p)
	}
}

func Test_Priority_Valid_UnknownValue_ReturnFalse(t *testing.T) {
	assert.False(t, issue.Priority("critical").Valid())
}

// — Issue.UpdateDescription —

func Test_UpdateDescription_NonNilValue_SetsDescription(t *testing.T) {
	base := baseIssue()
	desc := "new description"

	actual, err := base.UpdateDescription(&desc)

	require.NoError(t, err)
	require.NotNil(t, actual.Description)
	assert.Equal(t, desc, *actual.Description)
}

func Test_UpdateDescription_NilDescription_ClearsDescription(t *testing.T) {
	existing := "old description"
	base := baseIssue()
	base.Description = &existing

	actual, err := base.UpdateDescription(nil)

	require.NoError(t, err)
	assert.Nil(t, actual.Description)
}

func Test_UpdateDescription_NonNilValue_OriginalUnchanged(t *testing.T) {
	base := baseIssue()
	desc := "changed"

	_, err := base.UpdateDescription(&desc)

	require.NoError(t, err)
	assert.Nil(t, base.Description)
}

// — Issue.UpdatePriority —

func Test_UpdatePriority_ValidPriority_SetsPriority(t *testing.T) {
	base := baseIssue()

	actual, err := base.UpdatePriority(issue.PriorityHigh)

	require.NoError(t, err)
	assert.Equal(t, issue.PriorityHigh, actual.Priority)
}

func Test_UpdatePriority_InvalidPriority_ReturnsError(t *testing.T) {
	base := baseIssue()

	_, err := base.UpdatePriority(issue.Priority("extreme"))

	require.Error(t, err)
	assert.ErrorIs(t, err, issue.ErrInvalidIssue)
}

func Test_UpdatePriority_ValidPriority_OriginalUnchanged(t *testing.T) {
	base := baseIssue()
	base.Priority = issue.PriorityNone

	_, err := base.UpdatePriority(issue.PriorityUrgent)

	require.NoError(t, err)
	assert.Equal(t, issue.PriorityNone, base.Priority)
}

// — Issue.UpdateStatus —

func Test_UpdateStatus_ValidStatus_SetsStatus(t *testing.T) {
	base := baseIssue()

	actual, err := base.UpdateStatus(issue.StatusDone)

	require.NoError(t, err)
	assert.Equal(t, issue.StatusDone, actual.Status)
}

func Test_UpdateStatus_InvalidStatus_ReturnsError(t *testing.T) {
	base := baseIssue()

	_, err := base.UpdateStatus(issue.Status("archived"))

	require.Error(t, err)
	assert.ErrorIs(t, err, issue.ErrInvalidIssue)
}

func Test_UpdateStatus_ValidStatus_OriginalUnchanged(t *testing.T) {
	base := baseIssue()
	base.Status = issue.StatusTodo

	_, err := base.UpdateStatus(issue.StatusDone)

	require.NoError(t, err)
	assert.Equal(t, issue.StatusTodo, base.Status)
}

// — Issue.UpdateAssignee —

func Test_UpdateAssignee_ValidAssigneeID_SetsAssigneeID(t *testing.T) {
	base := baseIssue()
	assigneeID := uuid.New()

	actual, err := base.UpdateAssignee(assigneeID)

	require.NoError(t, err)
	require.NotNil(t, actual.AssigneeID)
	assert.Equal(t, assigneeID, *actual.AssigneeID)
}

func Test_UpdateAssignee_ValidAssigneeID_OriginalUnchanged(t *testing.T) {
	base := baseIssue()
	assigneeID := uuid.New()

	_, err := base.UpdateAssignee(assigneeID)

	require.NoError(t, err)
	assert.Nil(t, base.AssigneeID)
}

// baseIssue returns a minimal valid Issue for use in tests.
func baseIssue() issue.Issue {
	return issue.Issue{
		ID:         uuid.New(),
		Identifier: "test-issue-abc",
		Title:      "Test issue",
		Status:     issue.StatusTodo,
		Priority:   issue.PriorityNone,
		Labels:     []string{},
		ProjectID:  uuid.New(),
		ReporterID: uuid.New(),
	}
}
