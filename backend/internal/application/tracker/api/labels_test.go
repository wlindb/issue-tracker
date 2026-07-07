//go:build !integration

package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/wlindb/issue-tracker/internal/application/tracker/api"
	"github.com/wlindb/issue-tracker/internal/application/tracker/api/model"
	labeldomain "github.com/wlindb/issue-tracker/internal/domain/tracker/label"
)

type mockLabelService struct {
	mock.Mock
}

func (m *mockLabelService) Create(ctx context.Context, name string) (labeldomain.Label, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(labeldomain.Label), args.Error(1)
}

func (m *mockLabelService) Search(ctx context.Context, name string) ([]labeldomain.Label, error) {
	args := m.Called(ctx, name)
	result, _ := args.Get(0).([]labeldomain.Label)
	return result, args.Error(1)
}

func newLabelTestServer(t *testing.T, service api.LabelService) *echo.Echo {
	t.Helper()
	e := echo.New()
	h := &api.Handler{LabelHandler: api.NewLabelHandler(service)}
	strict := model.NewStrictHandler(h, nil)
	model.RegisterHandlersWithBaseURL(e, strict, "/api/v1")
	return e
}

// --- POST /workspaces/{workspaceId}/labels ---

func Test_CreateLabel_ValidBody_Returns201(t *testing.T) {
	service := &mockLabelService{}
	service.On("Create", mock.Anything, "Bug").Return(labeldomain.Label{ID: uuid.New(), Name: "Bug"}, nil)
	e := newLabelTestServer(t, service)
	e.Use(injectUser(uuid.New()))

	req := httptest.NewRequest(http.MethodPost, wsPath("/labels"), strings.NewReader(`{"name":"Bug"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	var got model.Label
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	service.AssertExpectations(t)
}

func Test_CreateLabel_EmptyName_Returns400(t *testing.T) {
	service := &mockLabelService{}
	e := newLabelTestServer(t, service)
	e.Use(injectUser(uuid.New()))

	req := httptest.NewRequest(http.MethodPost, wsPath("/labels"), strings.NewReader(`{"name":""}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	var got model.Error
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	assert.Equal(t, "invalid_input", got.Code)
	service.AssertNotCalled(t, "Create")
}

func Test_CreateLabel_NilBody_Returns400(t *testing.T) {
	service := &mockLabelService{}
	e := newLabelTestServer(t, service)
	e.Use(injectUser(uuid.New()))

	req := httptest.NewRequest(http.MethodPost, wsPath("/labels"), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	service.AssertNotCalled(t, "Create")
}

// --- GET /workspaces/{workspaceId}/labels ---

func Test_ListLabels_Returns200WithLabels(t *testing.T) {
	id1, id2 := uuid.New(), uuid.New()
	expected := []labeldomain.Label{
		{ID: id1, Name: "Bug"},
		{ID: id2, Name: "Feature"},
	}
	service := &mockLabelService{}
	service.On("Search", mock.Anything, "").Return(expected, nil)
	e := newLabelTestServer(t, service)
	e.Use(injectUser(uuid.New()))

	req := httptest.NewRequest(http.MethodGet, wsPath("/labels"), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var actual model.LabelPage
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &actual))
	require.Len(t, actual.Items, 2)
	assert.Equal(t, id1, actual.Items[0].Id)
	assert.Equal(t, "Bug", actual.Items[0].Name)
	assert.Equal(t, id2, actual.Items[1].Id)
	assert.Equal(t, "Feature", actual.Items[1].Name)
	service.AssertExpectations(t)
}
