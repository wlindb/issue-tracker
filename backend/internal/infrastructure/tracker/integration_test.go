//go:build integration

package tracker_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	issuedomain "github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
	projectdomain "github.com/wlindb/issue-tracker/internal/domain/tracker/project"
	infradb "github.com/wlindb/issue-tracker/internal/infrastructure/db"
	tracker "github.com/wlindb/issue-tracker/internal/infrastructure/tracker"
)

var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()
	pool, terminate, err := startPostgres(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "start postgres: %v\n", err)
		os.Exit(1)
	}

	if err := tracker.Migrate(ctx, pool); err != nil {
		fmt.Fprintf(os.Stderr, "migrate: %v\n", err)
		os.Exit(1)
	}

	testPool = pool
	code := m.Run()
	terminate()
	os.Exit(code)
}

func startPostgres(ctx context.Context) (*pgxpool.Pool, func(), error) {
	req := testcontainers.ContainerRequest{
		Image: "postgres:17-alpine",
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "test",
		},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
	}
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req, Started: true,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("start container: %w", err)
	}
	port, err := c.MappedPort(ctx, "5432")
	if err != nil {
		return nil, nil, errors.Join(fmt.Errorf("mapped port: %w", err), c.Terminate(ctx))
	}
	host, err := c.Host(ctx)
	if err != nil {
		return nil, nil, errors.Join(fmt.Errorf("host: %w", err), c.Terminate(ctx))
	}
	dsn := fmt.Sprintf("postgres://test:test@%s:%s/test", host, port.Port())
	pool, err := infradb.New(ctx, dsn)
	if err != nil {
		return nil, nil, errors.Join(fmt.Errorf("connect: %w", err), c.Terminate(ctx))
	}
	return pool, func() {
		pool.Close()
		if err := c.Terminate(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "terminate container: %v\n", err)
		}
	}, nil
}

func Test_Create_NoDescription_SuccessfulProjectCreation(t *testing.T) {
	repository := tracker.NewProjectRepository(testPool)
	id, ownerID := uuid.New(), uuid.New()

	actual, err := repository.Create(context.Background(), id, ownerID, "Acme", nil)

	require.NoError(t, err)
	require.NotNil(t, actual)
	assert.Equal(t, id, actual.ID)
	assert.Equal(t, ownerID, actual.OwnerID)
	assert.Equal(t, "Acme", actual.Name)
	assert.Nil(t, actual.Description)
	assert.False(t, actual.CreatedAt.IsZero())
	assert.False(t, actual.UpdatedAt.IsZero())
}

func Test_Create_WithDescription_SuccessfulProjectCreation(t *testing.T) {
	repository := tracker.NewProjectRepository(testPool)
	description := "My description"

	actual, err := repository.Create(context.Background(), uuid.New(), uuid.New(), "Described", &description)

	require.NoError(t, err)
	require.NotNil(t, actual.Description)
	assert.Equal(t, description, *actual.Description)
}

func Test_Create_DuplicateID_ReturnsError(t *testing.T) {
	repository := tracker.NewProjectRepository(testPool)
	ctx := context.Background()
	id := uuid.New()

	_, err := repository.Create(ctx, id, uuid.New(), "First", nil)
	require.NoError(t, err)

	_, err = repository.Create(ctx, id, uuid.New(), "Second", nil)
	require.Error(t, err) // PK violation
}

func Test_List_WithLimit_ReturnsLimitedProjects(t *testing.T) {
	repository := tracker.NewProjectRepository(testPool)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		_, err := repository.Create(ctx, uuid.New(), uuid.New(), "LimitTest", nil)
		require.NoError(t, err)
	}

	limit := 2
	query := projectdomain.NewListProjectQuery(nil, &limit)
	actual, err := repository.List(ctx, query)

	require.NoError(t, err)
	assert.Len(t, actual.Items, 2)
}

func Test_List_LimitExceedsTotal_ReturnsAllProjects(t *testing.T) {
	repository := tracker.NewProjectRepository(testPool)
	ctx := context.Background()

	id1, id2 := uuid.New(), uuid.New()
	_, err := repository.Create(ctx, id1, uuid.New(), "ExceedA", nil)
	require.NoError(t, err)
	_, err = repository.Create(ctx, id2, uuid.New(), "ExceedB", nil)
	require.NoError(t, err)

	limit := 100
	query := projectdomain.NewListProjectQuery(nil, &limit)
	actual, err := repository.List(ctx, query)

	require.NoError(t, err)
	ids := make([]uuid.UUID, len(actual.Items))
	for i, p := range actual.Items {
		ids[i] = p.ID
	}
	assert.Contains(t, ids, id1)
	assert.Contains(t, ids, id2)
}

// — Issue integration helpers —

func createTestProject(t *testing.T) uuid.UUID {
	t.Helper()
	repo := tracker.NewProjectRepository(testPool)
	id := uuid.New()
	_, err := repo.Create(context.Background(), id, uuid.New(), "TestProject-"+id.String()[:8], nil)
	require.NoError(t, err)
	return id
}

// — CreateIssue full-flow integration tests —

func Test_CreateIssue_FromCommand_HasEmptyLabels(t *testing.T) {
	repository := tracker.NewIssueRepository(testPool)
	projectID := createTestProject(t)
	ctx := context.Background()

	command := issuedomain.CreateIssueCommand{
		ProjectID:  projectID,
		ReporterID: uuid.New(),
		Title:      "New feature",
		Status:     issuedomain.StatusTodo,
		Priority:   issuedomain.PriorityMedium,
	}
	issue := command.ToIssue(uuid.New(), command.Slugify)
	assert.NotNil(t, issue.Labels)
	assert.Empty(t, issue.Labels)

	actual, err := repository.CreateIssue(ctx, issue)

	require.NoError(t, err)
	require.NotNil(t, actual)
	assert.NotNil(t, actual.Labels)
	assert.Empty(t, actual.Labels)
}

// — CreateIssue integration tests —

func Test_CreateIssue_NoOptionalFields_SuccessfulCreation(t *testing.T) {
	repository := tracker.NewIssueRepository(testPool)
	projectID := createTestProject(t)
	ctx := context.Background()

	issue := issuedomain.Issue{
		ID:         uuid.New(),
		Identifier: "create-no-opt-" + uuid.New().String()[:8],
		Title:      "Simple issue",
		Status:     issuedomain.StatusBacklog,
		Priority:   issuedomain.PriorityNone,
		Labels:     []string{},
		ProjectID:  projectID,
		ReporterID: uuid.New(),
	}

	actual, err := repository.CreateIssue(ctx, issue)

	require.NoError(t, err)
	require.NotNil(t, actual)
	assert.Equal(t, issue.ID, actual.ID)
	assert.Equal(t, issue.Identifier, actual.Identifier)
	assert.Equal(t, issue.Title, actual.Title)
	assert.Nil(t, actual.Description)
	assert.Equal(t, issuedomain.StatusBacklog, actual.Status)
	assert.Equal(t, issuedomain.PriorityNone, actual.Priority)
	assert.Empty(t, actual.Labels)
	assert.Nil(t, actual.AssigneeID)
	assert.Equal(t, projectID, actual.ProjectID)
	assert.Equal(t, issue.ReporterID, actual.ReporterID)
	assert.False(t, actual.CreatedAt.IsZero())
	assert.False(t, actual.UpdatedAt.IsZero())
}

func Test_CreateIssue_WithOptionalFields_SuccessfulCreation(t *testing.T) {
	repository := tracker.NewIssueRepository(testPool)
	projectID := createTestProject(t)
	ctx := context.Background()

	description := "detailed description"
	assigneeID := uuid.New()
	issue := issuedomain.Issue{
		ID:          uuid.New(),
		Identifier:  "create-opt-" + uuid.New().String()[:8],
		Title:       "Full issue",
		Description: &description,
		Status:      issuedomain.StatusInProgress,
		Priority:    issuedomain.PriorityHigh,
		Labels:      []string{"backend", "urgent"},
		AssigneeID:  &assigneeID,
		ProjectID:   projectID,
		ReporterID:  uuid.New(),
	}

	actual, err := repository.CreateIssue(ctx, issue)

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.NotNil(t, actual.Description)
	assert.Equal(t, description, *actual.Description)
	require.NotNil(t, actual.AssigneeID)
	assert.Equal(t, assigneeID, *actual.AssigneeID)
	assert.Equal(t, issuedomain.StatusInProgress, actual.Status)
	assert.Equal(t, issuedomain.PriorityHigh, actual.Priority)
	assert.Equal(t, []string{"backend", "urgent"}, actual.Labels)
}

func Test_CreateIssue_DuplicateID_ReturnsError(t *testing.T) {
	repository := tracker.NewIssueRepository(testPool)
	projectID := createTestProject(t)
	ctx := context.Background()

	issueID := uuid.New()
	issue := issuedomain.Issue{
		ID:         issueID,
		Identifier: "dup-id-" + uuid.New().String()[:8],
		Title:      "First",
		Status:     issuedomain.StatusBacklog,
		Priority:   issuedomain.PriorityNone,
		Labels:     []string{},
		ProjectID:  projectID,
		ReporterID: uuid.New(),
	}
	_, err := repository.CreateIssue(ctx, issue)
	require.NoError(t, err)

	issue.Identifier = "dup-id-second-" + uuid.New().String()[:8]
	_, err = repository.CreateIssue(ctx, issue)
	require.Error(t, err) // PK violation
}

func Test_CreateIssue_DuplicateIdentifierSameProject_ReturnsError(t *testing.T) {
	repository := tracker.NewIssueRepository(testPool)
	projectID := createTestProject(t)
	ctx := context.Background()

	identifier := "dup-ident-" + uuid.New().String()[:8]
	issue := issuedomain.Issue{
		ID:         uuid.New(),
		Identifier: identifier,
		Title:      "First",
		Status:     issuedomain.StatusBacklog,
		Priority:   issuedomain.PriorityNone,
		Labels:     []string{},
		ProjectID:  projectID,
		ReporterID: uuid.New(),
	}
	_, err := repository.CreateIssue(ctx, issue)
	require.NoError(t, err)

	issue.ID = uuid.New()
	_, err = repository.CreateIssue(ctx, issue)
	require.Error(t, err) // UNIQUE violation on (project_id, identifier)
}

func Test_CreateIssue_DuplicateIdentifierDifferentProject_Succeeds(t *testing.T) {
	repository := tracker.NewIssueRepository(testPool)
	projectA := createTestProject(t)
	projectB := createTestProject(t)
	ctx := context.Background()

	identifier := "cross-proj-" + uuid.New().String()[:8]
	issueA := issuedomain.Issue{
		ID:         uuid.New(),
		Identifier: identifier,
		Title:      "Issue in A",
		Status:     issuedomain.StatusBacklog,
		Priority:   issuedomain.PriorityNone,
		Labels:     []string{},
		ProjectID:  projectA,
		ReporterID: uuid.New(),
	}
	_, err := repository.CreateIssue(ctx, issueA)
	require.NoError(t, err)

	issueB := issuedomain.Issue{
		ID:         uuid.New(),
		Identifier: identifier,
		Title:      "Issue in B",
		Status:     issuedomain.StatusBacklog,
		Priority:   issuedomain.PriorityNone,
		Labels:     []string{},
		ProjectID:  projectB,
		ReporterID: uuid.New(),
	}
	_, err = repository.CreateIssue(ctx, issueB)
	require.NoError(t, err) // same identifier, different project is OK
}

func Test_CreateIssue_InvalidProjectID_ReturnsError(t *testing.T) {
	repository := tracker.NewIssueRepository(testPool)
	ctx := context.Background()

	issue := issuedomain.Issue{
		ID:         uuid.New(),
		Identifier: "invalid-proj-" + uuid.New().String()[:8],
		Title:      "Bad project ref",
		Status:     issuedomain.StatusBacklog,
		Priority:   issuedomain.PriorityNone,
		Labels:     []string{},
		ProjectID:  uuid.New(), // nonexistent project
		ReporterID: uuid.New(),
	}

	_, err := repository.CreateIssue(ctx, issue)
	require.Error(t, err) // FK violation
}

// — ListIssues integration tests —

func Test_ListIssues_EmptyProject_ReturnsEmptyPage(t *testing.T) {
	repository := tracker.NewIssueRepository(testPool)
	projectID := createTestProject(t)
	ctx := context.Background()

	query := issuedomain.ListIssueQuery{}
	actual, err := repository.ListIssues(ctx, projectID, query)

	require.NoError(t, err)
	assert.Empty(t, actual.Items)
}

func Test_ListIssues_WithIssues_ReturnsAllIssues(t *testing.T) {
	repository := tracker.NewIssueRepository(testPool)
	projectID := createTestProject(t)
	ctx := context.Background()

	for idx := 0; idx < 3; idx++ {
		_, err := repository.CreateIssue(ctx, issuedomain.Issue{
			ID:         uuid.New(),
			Identifier: fmt.Sprintf("list-all-%d-%s", idx, uuid.New().String()[:8]),
			Title:      fmt.Sprintf("Issue %d", idx),
			Status:     issuedomain.StatusBacklog,
			Priority:   issuedomain.PriorityNone,
			Labels:     []string{},
			ProjectID:  projectID,
			ReporterID: uuid.New(),
		})
		require.NoError(t, err)
	}

	query := issuedomain.ListIssueQuery{}
	actual, err := repository.ListIssues(ctx, projectID, query)

	require.NoError(t, err)
	assert.Len(t, actual.Items, 3)
}

func Test_ListIssues_FilterByStatus_ReturnsFilteredIssues(t *testing.T) {
	repository := tracker.NewIssueRepository(testPool)
	projectID := createTestProject(t)
	ctx := context.Background()

	statuses := []issuedomain.Status{issuedomain.StatusBacklog, issuedomain.StatusTodo, issuedomain.StatusBacklog}
	for idx, status := range statuses {
		_, err := repository.CreateIssue(ctx, issuedomain.Issue{
			ID:         uuid.New(),
			Identifier: fmt.Sprintf("filter-status-%d-%s", idx, uuid.New().String()[:8]),
			Title:      fmt.Sprintf("Issue %d", idx),
			Status:     status,
			Priority:   issuedomain.PriorityNone,
			Labels:     []string{},
			ProjectID:  projectID,
			ReporterID: uuid.New(),
		})
		require.NoError(t, err)
	}

	filterStatus := issuedomain.StatusBacklog
	query := issuedomain.ListIssueQuery{Status: &filterStatus}
	actual, err := repository.ListIssues(ctx, projectID, query)

	require.NoError(t, err)
	assert.Len(t, actual.Items, 2)
	for _, item := range actual.Items {
		assert.Equal(t, issuedomain.StatusBacklog, item.Status)
	}
}

func Test_ListIssues_FilterByPriority_ReturnsFilteredIssues(t *testing.T) {
	repository := tracker.NewIssueRepository(testPool)
	projectID := createTestProject(t)
	ctx := context.Background()

	priorities := []issuedomain.Priority{issuedomain.PriorityHigh, issuedomain.PriorityLow, issuedomain.PriorityHigh}
	for idx, priority := range priorities {
		_, err := repository.CreateIssue(ctx, issuedomain.Issue{
			ID:         uuid.New(),
			Identifier: fmt.Sprintf("filter-priority-%d-%s", idx, uuid.New().String()[:8]),
			Title:      fmt.Sprintf("Issue %d", idx),
			Status:     issuedomain.StatusBacklog,
			Priority:   priority,
			Labels:     []string{},
			ProjectID:  projectID,
			ReporterID: uuid.New(),
		})
		require.NoError(t, err)
	}

	filterPriority := issuedomain.PriorityHigh
	query := issuedomain.ListIssueQuery{Priority: &filterPriority}
	actual, err := repository.ListIssues(ctx, projectID, query)

	require.NoError(t, err)
	assert.Len(t, actual.Items, 2)
	for _, item := range actual.Items {
		assert.Equal(t, issuedomain.PriorityHigh, item.Priority)
	}
}

func Test_ListIssues_FilterByAssignee_ReturnsFilteredIssues(t *testing.T) {
	repository := tracker.NewIssueRepository(testPool)
	projectID := createTestProject(t)
	ctx := context.Background()

	assigneeID := uuid.New()
	otherID := uuid.New()
	assignees := []*uuid.UUID{&assigneeID, &otherID, &assigneeID}
	for idx, assignee := range assignees {
		_, err := repository.CreateIssue(ctx, issuedomain.Issue{
			ID:         uuid.New(),
			Identifier: fmt.Sprintf("filter-assignee-%d-%s", idx, uuid.New().String()[:8]),
			Title:      fmt.Sprintf("Issue %d", idx),
			Status:     issuedomain.StatusBacklog,
			Priority:   issuedomain.PriorityNone,
			Labels:     []string{},
			AssigneeID: assignee,
			ProjectID:  projectID,
			ReporterID: uuid.New(),
		})
		require.NoError(t, err)
	}

	query := issuedomain.ListIssueQuery{AssigneeID: &assigneeID}
	actual, err := repository.ListIssues(ctx, projectID, query)

	require.NoError(t, err)
	assert.Len(t, actual.Items, 2)
	for _, item := range actual.Items {
		require.NotNil(t, item.AssigneeID)
		assert.Equal(t, assigneeID, *item.AssigneeID)
	}
}

func Test_ListIssues_IsolatesByProject_ReturnsOnlyProjectIssues(t *testing.T) {
	repository := tracker.NewIssueRepository(testPool)
	projectA := createTestProject(t)
	projectB := createTestProject(t)
	ctx := context.Background()

	_, err := repository.CreateIssue(ctx, issuedomain.Issue{
		ID:         uuid.New(),
		Identifier: "proj-a-" + uuid.New().String()[:8],
		Title:      "Issue in A",
		Status:     issuedomain.StatusBacklog,
		Priority:   issuedomain.PriorityNone,
		Labels:     []string{},
		ProjectID:  projectA,
		ReporterID: uuid.New(),
	})
	require.NoError(t, err)

	_, err = repository.CreateIssue(ctx, issuedomain.Issue{
		ID:         uuid.New(),
		Identifier: "proj-b-" + uuid.New().String()[:8],
		Title:      "Issue in B",
		Status:     issuedomain.StatusBacklog,
		Priority:   issuedomain.PriorityNone,
		Labels:     []string{},
		ProjectID:  projectB,
		ReporterID: uuid.New(),
	})
	require.NoError(t, err)

	query := issuedomain.ListIssueQuery{}
	actual, err := repository.ListIssues(ctx, projectA, query)

	require.NoError(t, err)
	assert.Len(t, actual.Items, 1)
	assert.Equal(t, "Issue in A", actual.Items[0].Title)
}
