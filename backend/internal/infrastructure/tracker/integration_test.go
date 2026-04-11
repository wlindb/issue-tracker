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
	workspacedomain "github.com/wlindb/issue-tracker/internal/domain/tracker/workspace"
	infradb "github.com/wlindb/issue-tracker/internal/infrastructure/db"
	tracker "github.com/wlindb/issue-tracker/internal/infrastructure/tracker"
)

type testContextKey string

const (
	testWorkspaceIDKey testContextKey = "workspace_id"
	testUserIDKey      testContextKey = "user_id"
)

func withWorkspaceContext(workspaceID uuid.UUID, userID uuid.UUID) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, testWorkspaceIDKey, workspaceID)
	ctx = context.WithValue(ctx, testUserIDKey, userID)
	return ctx
}

var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()
	dsn, terminate, err := startPostgres(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "start postgres: %v\n", err)
		os.Exit(1)
	}

	// Run migrations as the superuser (plain pool, no role switching).
	migrationPool, err := infradb.New(ctx, dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "start postgres: connect migration pool: %v\n", err)
		os.Exit(1)
	}
	if err := tracker.Migrate(ctx, migrationPool); err != nil {
		fmt.Fprintf(os.Stderr, "migrate: %v\n", err)
		os.Exit(1)
	}
	migrationPool.Close()

	// Open the application pool after migrations: appuser role now exists.
	testPool, err = infradb.New(ctx, dsn,
		infradb.WithAppSessionVars(
			func(ctx context.Context) (uuid.UUID, error) {
				id, ok := ctx.Value(testWorkspaceIDKey).(uuid.UUID)
				if !ok || id == uuid.Nil {
					return uuid.Nil, errors.New("missing workspace ID")
				}
				return id, nil
			},
			func(ctx context.Context) (uuid.UUID, error) {
				id, ok := ctx.Value(testUserIDKey).(uuid.UUID)
				if !ok || id == uuid.Nil {
					return uuid.Nil, errors.New("missing user ID")
				}
				return id, nil
			},
		),
		infradb.WithAppRole("appuser"),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "start postgres: connect app pool: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()
	testPool.Close()
	terminate()
	os.Exit(code)
}

func startPostgres(ctx context.Context) (string, func(), error) {
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
		return "", nil, fmt.Errorf("start container: %w", err)
	}
	port, err := c.MappedPort(ctx, "5432")
	if err != nil {
		return "", nil, errors.Join(fmt.Errorf("mapped port: %w", err), c.Terminate(ctx))
	}
	host, err := c.Host(ctx)
	if err != nil {
		return "", nil, errors.Join(fmt.Errorf("host: %w", err), c.Terminate(ctx))
	}
	dsn := fmt.Sprintf("postgres://test:test@%s:%s/test", host, port.Port())
	return dsn, func() {
		if err := c.Terminate(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "terminate container: %v\n", err)
		}
	}, nil
}

func createTestWorkspace(t *testing.T) (uuid.UUID, context.Context) {
	t.Helper()
	workspaceRepository := tracker.NewWorkspaceRepository(testPool)
	workspaceID := uuid.New()
	ownerID := uuid.New()
	_, err := workspaceRepository.Create(context.Background(), workspacedomain.Workspace{ID: workspaceID, OwnerID: ownerID, Name: "Test Workspace"})
	require.NoError(t, err)
	return workspaceID, withWorkspaceContext(workspaceID, ownerID)
}

func createTestProject(t *testing.T, ctx context.Context) uuid.UUID {
	t.Helper()
	repository := tracker.NewProjectRepository(testPool)
	id := uuid.New()
	_, err := repository.Create(ctx, projectdomain.Project{
		ID:         id,
		Identifier: "test-" + id.String()[:8],
		OwnerID:    uuid.New(),
		Name:       "TestProject-" + id.String()[:8],
	})
	require.NoError(t, err)
	return id
}

func Test_Create_NoDescription_SuccessfulProjectCreation(t *testing.T) {
	repository := tracker.NewProjectRepository(testPool)
	_, ctx := createTestWorkspace(t)
	id, ownerID := uuid.New(), uuid.New()

	actual, err := repository.Create(ctx, projectdomain.Project{
		ID:         id,
		Identifier: "acme-" + id.String()[:8],
		OwnerID:    ownerID,
		Name:       "Acme",
	})

	require.NoError(t, err)
	assert.Equal(t, id, actual.ID)
	assert.Equal(t, ownerID, actual.OwnerID)
	assert.Equal(t, "Acme", actual.Name)
	assert.Nil(t, actual.Description)
	assert.False(t, actual.CreatedAt.IsZero())
	assert.False(t, actual.UpdatedAt.IsZero())
}

func Test_Create_WithDescription_SuccessfulProjectCreation(t *testing.T) {
	repository := tracker.NewProjectRepository(testPool)
	_, ctx := createTestWorkspace(t)
	id := uuid.New()
	description := "My description"

	actual, err := repository.Create(ctx, projectdomain.Project{
		ID:          id,
		Identifier:  "described-" + id.String()[:8],
		OwnerID:     uuid.New(),
		Name:        "Described",
		Description: &description,
	})

	require.NoError(t, err)
	require.NotNil(t, actual.Description)
	assert.Equal(t, description, *actual.Description)
}

func Test_Create_DuplicateID_ReturnsError(t *testing.T) {
	repository := tracker.NewProjectRepository(testPool)
	_, ctx := createTestWorkspace(t)
	id := uuid.New()

	_, err := repository.Create(ctx, projectdomain.Project{
		ID:         id,
		Identifier: "first-" + id.String()[:8],
		OwnerID:    uuid.New(),
		Name:       "First",
	})
	require.NoError(t, err)

	_, err = repository.Create(ctx, projectdomain.Project{
		ID:         id,
		Identifier: "second-" + id.String()[:8],
		OwnerID:    uuid.New(),
		Name:       "Second",
	})
	require.Error(t, err) // PK violation
}

func Test_List_WithLimit_ReturnsLimitedProjects(t *testing.T) {
	repository := tracker.NewProjectRepository(testPool)
	_, ctx := createTestWorkspace(t)

	for i := 0; i < 3; i++ {
		id := uuid.New()
		_, err := repository.Create(ctx, projectdomain.Project{
			ID:         id,
			Identifier: fmt.Sprintf("limit-%d-%s", i, id.String()[:8]),
			OwnerID:    uuid.New(),
			Name:       "LimitTest",
		})
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
	_, ctx := createTestWorkspace(t)

	id1, id2 := uuid.New(), uuid.New()
	_, err := repository.Create(ctx, projectdomain.Project{
		ID:         id1,
		Identifier: "exceed-a-" + id1.String()[:8],
		OwnerID:    uuid.New(),
		Name:       "ExceedA",
	})
	require.NoError(t, err)
	_, err = repository.Create(ctx, projectdomain.Project{
		ID:         id2,
		Identifier: "exceed-b-" + id2.String()[:8],
		OwnerID:    uuid.New(),
		Name:       "ExceedB",
	})
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

// — CreateIssue full-flow integration tests —

func Test_CreateIssue_FromCommand_HasEmptyLabels(t *testing.T) {
	repository := tracker.NewIssueRepository(testPool)
	_, ctx := createTestWorkspace(t)
	projectID := createTestProject(t, ctx)

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
	_, ctx := createTestWorkspace(t)
	projectID := createTestProject(t, ctx)

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
	_, ctx := createTestWorkspace(t)
	projectID := createTestProject(t, ctx)

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
	_, ctx := createTestWorkspace(t)
	projectID := createTestProject(t, ctx)

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
	_, ctx := createTestWorkspace(t)
	projectID := createTestProject(t, ctx)

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
	_, ctx := createTestWorkspace(t)
	projectA := createTestProject(t, ctx)
	projectB := createTestProject(t, ctx)

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
	require.Error(t, err)
}

// — RLS integration tests —

func Test_Projects_RLS_NonMember_HidesRows(t *testing.T) {
	workspaceID, ctx := createTestWorkspace(t)
	createTestProject(t, ctx)

	nonMemberID := uuid.New()
	ctxNonMember := withWorkspaceContext(workspaceID, nonMemberID)
	repository := tracker.NewProjectRepository(testPool)
	limit := 100
	query := projectdomain.NewListProjectQuery(nil, &limit)

	actual, err := repository.List(ctxNonMember, query)

	require.NoError(t, err)
	assert.Empty(t, actual.Items)
}

func Test_GetWorkspace_AsMember_ReturnsWorkspace(t *testing.T) {
	workspaceID, ctx := createTestWorkspace(t)
	repository := tracker.NewWorkspaceRepository(testPool)

	actual, err := repository.Get(ctx, workspaceID)

	require.NoError(t, err)
	require.NotNil(t, actual)
	assert.Equal(t, workspaceID, actual.ID)
}

func Test_GetWorkspace_AsNonMember_ReturnsErrWorkspaceNotFound(t *testing.T) {
	workspaceID, _ := createTestWorkspace(t)
	nonMemberCtx := withWorkspaceContext(workspaceID, uuid.New())
	repository := tracker.NewWorkspaceRepository(testPool)

	_, err := repository.Get(nonMemberCtx, workspaceID)

	require.Error(t, err)
	assert.ErrorIs(t, err, workspacedomain.ErrWorkspaceNotFound)
}

func Test_Issues_CrossWorkspaceProjectID_ReturnsError(t *testing.T) {
	_, ctxA := createTestWorkspace(t)
	_, ctxB := createTestWorkspace(t)
	projectA := createTestProject(t, ctxA)

	// Issue workspace_id comes from current_setting (workspace B), but project_id
	// belongs to workspace A. The composite FK (workspace_id, project_id) →
	// projects(workspace_id, id) rejects this mismatch.
	_, err := testPool.Exec(ctxB,
		`INSERT INTO issues (id, identifier, title, status, priority, labels, project_id, reporter_id, workspace_id, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, current_setting('app.workspace_id')::uuid, NOW(), NOW())`,
		uuid.New(), "cross-ws-"+uuid.New().String()[:8], "Cross Workspace Issue", "backlog", "none", []string{}, projectA, uuid.New(),
	)

	require.Error(t, err)
}

// — ListIssues integration tests —

func Test_ListIssues_EmptyProject_ReturnsEmptyPage(t *testing.T) {
	repository := tracker.NewIssueRepository(testPool)
	_, ctx := createTestWorkspace(t)
	projectID := createTestProject(t, ctx)

	query := issuedomain.ListIssueQuery{}
	actual, err := repository.ListIssues(ctx, projectID, query)

	require.NoError(t, err)
	assert.Empty(t, actual.Items)
}

func Test_ListIssues_WithIssues_ReturnsAllIssues(t *testing.T) {
	repository := tracker.NewIssueRepository(testPool)
	_, ctx := createTestWorkspace(t)
	projectID := createTestProject(t, ctx)

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
	_, ctx := createTestWorkspace(t)
	projectID := createTestProject(t, ctx)

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
	_, ctx := createTestWorkspace(t)
	projectID := createTestProject(t, ctx)

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
	_, ctx := createTestWorkspace(t)
	projectID := createTestProject(t, ctx)

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

func Test_ListIssues_IsolatedByWorkspace_ReturnsEmpty(t *testing.T) {
	repository := tracker.NewIssueRepository(testPool)
	_, ctxA := createTestWorkspace(t)
	_, ctxB := createTestWorkspace(t)
	projectID := createTestProject(t, ctxA)

	_, err := repository.CreateIssue(ctxA, issuedomain.Issue{
		ID:         uuid.New(),
		Identifier: "isolated-" + uuid.New().String()[:8],
		Title:      "Issue in Workspace A",
		Status:     issuedomain.StatusBacklog,
		Priority:   issuedomain.PriorityNone,
		Labels:     []string{},
		ProjectID:  projectID,
		ReporterID: uuid.New(),
	})
	require.NoError(t, err)

	query := issuedomain.ListIssueQuery{}
	actual, err := repository.ListIssues(ctxB, projectID, query)

	require.NoError(t, err)
	assert.Empty(t, actual.Items)
}

func Test_ListIssues_IsolatesByProject_ReturnsOnlyProjectIssues(t *testing.T) {
	repository := tracker.NewIssueRepository(testPool)
	_, ctx := createTestWorkspace(t)
	projectA := createTestProject(t, ctx)
	projectB := createTestProject(t, ctx)

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
