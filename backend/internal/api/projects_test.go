package api_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/wlindb/issue-tracker/internal/api"
	"github.com/wlindb/issue-tracker/internal/api/generated"
	trackerdomain "github.com/wlindb/issue-tracker/internal/domain/tracker/project"
)

// mockProjectService implements the projectServicer interface for testing.
type mockProjectService struct {
	mock.Mock
}

func (m *mockProjectService) Create(ctx context.Context, ownerID uuid.UUID, name string, description *string) (*trackerdomain.Project, error) {
	args := m.Called(ctx, ownerID, name, description)
	if p, ok := args.Get(0).(*trackerdomain.Project); ok {
		return p, args.Error(1)
	}
	return nil, args.Error(1)
}

// newTestServer builds a minimal Echo server wired to the given project service.
func newTestServer(t *testing.T, svc api.ProjectServicer) *echo.Echo {
	t.Helper()
	e := echo.New()
	h := &api.Handler{
		ProjectHandler: api.NewProjectHandler(svc),
	}
	strict := generated.NewStrictHandler(h, nil)
	generated.RegisterHandlersWithBaseURL(e, strict, "/api/v1")
	return e
}

// injectUser returns Echo middleware that injects a fixed caller UUID into the
// request context, simulating what a future JWT middleware would do.
func injectUser(id uuid.UUID) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := api.WithUserID(c.Request().Context(), id)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

func Test_CreateProject_ValidBody_Returns201(t *testing.T) {
	svc := &mockProjectService{}
	now := time.Now().UTC()
	ownerID := uuid.New()

	svc.On("Create", mock.Anything, ownerID, "Acme", (*string)(nil)).
		Return(&trackerdomain.Project{
			ID:        uuid.New(),
			Name:      "Acme",
			OwnerID:   ownerID,
			CreatedAt: now,
			UpdatedAt: now,
		}, nil)

	e := newTestServer(t, svc)
	e.Use(injectUser(ownerID))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/projects", strings.NewReader(`{"name":"Acme"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	var got generated.Project
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	assert.Equal(t, "Acme", got.Name)
	svc.AssertExpectations(t)
}

func Test_CreateProject_EmptyName_BadRequest(t *testing.T) {
	svc := &mockProjectService{}

	e := newTestServer(t, svc)
	e.Use(injectUser(uuid.New()))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/projects", strings.NewReader(`{}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	svc.AssertNotCalled(t, "Create")
}

func Test_CreateProject_ServiceError_InternalServerError(t *testing.T) {
	svc := &mockProjectService{}
	ownerID := uuid.New()

	svc.On("Create", mock.Anything, ownerID, "Acme", (*string)(nil)).
		Return(nil, errors.New("db down"))

	e := newTestServer(t, svc)
	e.Use(injectUser(ownerID))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/projects", strings.NewReader(`{"name":"Acme"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	svc.AssertExpectations(t)
}
