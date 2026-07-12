//go:build !integration

package search_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	rootapi "github.com/wlindb/issue-tracker/internal/application/api"
	"github.com/wlindb/issue-tracker/internal/application/api/model"
	"github.com/wlindb/issue-tracker/internal/application/search"
)

var testWorkspaceID = uuid.MustParse("00000000-0000-0000-0000-000000000099")

func newSearchTestServer(t *testing.T) *echo.Echo {
	t.Helper()
	e := echo.New()
	h := &rootapi.Handler{SearchHandler: search.NewSearchHandler()}
	strict := model.NewStrictHandler(h, nil)
	model.RegisterHandlersWithBaseURL(e, strict, "/api/v1")
	return e
}

func searchIssuesPath() string {
	return "/api/v1/workspaces/" + testWorkspaceID.String() + "/search/issues"
}

func Test_SearchIssues_MissingBody_ReturnsBadRequest(t *testing.T) {
	e := newSearchTestServer(t)

	req := httptest.NewRequest(http.MethodPost, searchIssuesPath(), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func Test_SearchIssues_MissingQuery_ReturnsBadRequest(t *testing.T) {
	e := newSearchTestServer(t)

	req := httptest.NewRequest(http.MethodPost, searchIssuesPath(), strings.NewReader(`{}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func Test_SearchIssues_EmptyQuery_ReturnsBadRequest(t *testing.T) {
	e := newSearchTestServer(t)

	req := httptest.NewRequest(http.MethodPost, searchIssuesPath(), strings.NewReader(`{"query":""}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}
