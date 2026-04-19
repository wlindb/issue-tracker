//go:build integration

package api_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/wlindb/issue-tracker/internal/api"
	"github.com/wlindb/issue-tracker/internal/api/model"
	workspacedomain "github.com/wlindb/issue-tracker/internal/domain/tracker/workspace"
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
	container, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithDatabase("test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
		),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("start container: %w", err)
	}
	dsn, err := container.ConnectionString(ctx)
	if err != nil {
		return nil, nil, errors.Join(fmt.Errorf("connection string: %w", err), container.Terminate(ctx))
	}
	pool, err := infradb.New(ctx, dsn, infradb.WithAppSessionVars(
		api.WorkspaceIDFromContext,
		api.UserIDFromContext,
	))
	if err != nil {
		return nil, nil, errors.Join(fmt.Errorf("connect: %w", err), container.Terminate(ctx))
	}
	return pool, func() {
		pool.Close()
		if err := container.Terminate(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "terminate container: %v\n", err)
		}
	}, nil
}

func newWorkspaceIntegrationServer(t *testing.T) *echo.Echo {
	t.Helper()
	repository := tracker.NewWorkspaceRepository(testPool)
	service := workspacedomain.NewWorkspaceService(repository)
	handler := api.NewWorkspaceHandler(service)
	h := &api.Handler{WorkspaceHandler: handler}
	e := echo.New()
	strict := model.NewStrictHandler(h, nil)
	model.RegisterHandlersWithBaseURL(e, strict, "/api/v1")
	return e
}

func Test_CreateWorkspace_ValidRequest_Returns201(t *testing.T) {
	ownerID := uuid.New()
	e := newWorkspaceIntegrationServer(t)
	e.Use(injectUser(ownerID))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/workspaces", strings.NewReader(`{"name":"Acme"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	var actual model.Workspace
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &actual))
	assert.Equal(t, "Acme", actual.Name)
	assert.NotEqual(t, uuid.Nil, actual.Id)
	assert.Equal(t, ownerID, actual.OwnerId)
}

func Test_ListWorkspaces_WithExistingWorkspace_Returns200(t *testing.T) {
	ownerID := uuid.New()
	e := newWorkspaceIntegrationServer(t)
	e.Use(injectUser(ownerID))

	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/workspaces", strings.NewReader(`{"name":"ListTest"}`))
	createReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	createRec := httptest.NewRecorder()
	e.ServeHTTP(createRec, createReq)
	require.Equal(t, http.StatusCreated, createRec.Code)

	var created model.Workspace
	require.NoError(t, json.Unmarshal(createRec.Body.Bytes(), &created))

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/workspaces", nil)
	listRec := httptest.NewRecorder()
	e.ServeHTTP(listRec, listReq)

	require.Equal(t, http.StatusOK, listRec.Code)
	var actual model.WorkspacePage
	require.NoError(t, json.Unmarshal(listRec.Body.Bytes(), &actual))
	ids := make([]uuid.UUID, len(actual.Items))
	for index, item := range actual.Items {
		ids[index] = item.Id
	}
	assert.Contains(t, ids, created.Id)
}

func Test_GetWorkspace_WithExistingWorkspace_Returns200(t *testing.T) {
	ownerID := uuid.New()
	e := newWorkspaceIntegrationServer(t)
	e.Use(injectUser(ownerID))

	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/workspaces", strings.NewReader(`{"name":"GetTest"}`))
	createReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	createRec := httptest.NewRecorder()
	e.ServeHTTP(createRec, createReq)
	require.Equal(t, http.StatusCreated, createRec.Code)

	var created model.Workspace
	require.NoError(t, json.Unmarshal(createRec.Body.Bytes(), &created))

	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/workspaces/"+created.Id.String(), nil)
	getRec := httptest.NewRecorder()
	e.ServeHTTP(getRec, getReq)

	require.Equal(t, http.StatusOK, getRec.Code)
	var actual model.Workspace
	require.NoError(t, json.Unmarshal(getRec.Body.Bytes(), &actual))
	assert.Equal(t, created.Id, actual.Id)
	assert.Equal(t, "GetTest", actual.Name)
}

func Test_ListWorkspaceMembers_WithExistingWorkspace_Returns200(t *testing.T) {
	ownerID := uuid.New()
	e := newWorkspaceIntegrationServer(t)
	e.Use(injectUser(ownerID))

	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/workspaces", strings.NewReader(`{"name":"MembersTest"}`))
	createReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	createRec := httptest.NewRecorder()
	e.ServeHTTP(createRec, createReq)
	require.Equal(t, http.StatusCreated, createRec.Code)

	var created model.Workspace
	require.NoError(t, json.Unmarshal(createRec.Body.Bytes(), &created))

	membersReq := httptest.NewRequest(http.MethodGet, "/api/v1/workspaces/"+created.Id.String()+"/members", nil)
	membersRec := httptest.NewRecorder()
	e.ServeHTTP(membersRec, membersReq)

	require.Equal(t, http.StatusOK, membersRec.Code)
	var actual []model.WorkspaceMember
	require.NoError(t, json.Unmarshal(membersRec.Body.Bytes(), &actual))
	assert.NotNil(t, actual)
}
