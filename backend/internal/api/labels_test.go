//go:build !integration

package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wlindb/issue-tracker/internal/api"
	"github.com/wlindb/issue-tracker/internal/api/model"
)

func newLabelTestServer(t *testing.T) *echo.Echo {
	t.Helper()
	e := echo.New()
	h := &api.Handler{LabelHandler: api.NewLabelHandler()}
	strict := model.NewStrictHandler(h, nil)
	model.RegisterHandlersWithBaseURL(e, strict, "/api/v1")
	return e
}

// --- POST /workspaces/{workspaceId}/labels ---

func Test_CreateLabel_ValidBody_Returns201(t *testing.T) {
	e := newLabelTestServer(t)
	e.Use(injectUser(uuid.New()))

	req := httptest.NewRequest(http.MethodPost, wsPath("/labels"), strings.NewReader(`{"name":"Bug"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	var got model.Label
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
}

func Test_CreateLabel_EmptyName_Returns400(t *testing.T) {
	e := newLabelTestServer(t)
	e.Use(injectUser(uuid.New()))

	req := httptest.NewRequest(http.MethodPost, wsPath("/labels"), strings.NewReader(`{"name":""}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	var got model.Error
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	assert.Equal(t, "invalid_input", got.Code)
}

func Test_CreateLabel_NilBody_Returns400(t *testing.T) {
	e := newLabelTestServer(t)
	e.Use(injectUser(uuid.New()))

	req := httptest.NewRequest(http.MethodPost, wsPath("/labels"), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

// --- GET /workspaces/{workspaceId}/labels ---

func Test_ListLabels_Returns200WithEmptySlice(t *testing.T) {
	e := newLabelTestServer(t)
	e.Use(injectUser(uuid.New()))

	req := httptest.NewRequest(http.MethodGet, wsPath("/labels"), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got model.LabelPage
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	assert.NotNil(t, got.Items)
	assert.Empty(t, got.Items)
}
