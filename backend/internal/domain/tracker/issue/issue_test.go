package issue_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

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
