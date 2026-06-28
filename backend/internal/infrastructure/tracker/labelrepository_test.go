//go:build !integration

package tracker

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	issuedomain "github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
	trackerdb "github.com/wlindb/issue-tracker/internal/infrastructure/tracker/generated"
)

type mockLabelRow struct {
	label trackerdb.Label
	err   error
}

func (r *mockLabelRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}
	*(dest[0].(*uuid.UUID)) = r.label.ID
	*(dest[1].(*uuid.UUID)) = r.label.WorkspaceID
	*(dest[2].(*string)) = r.label.Name
	*(dest[3].(*pgtype.Timestamptz)) = r.label.CreatedAt
	return nil
}

type mockLabelRows struct {
	labels []trackerdb.Label
	index  int
}

func (r *mockLabelRows) Close()                                       {}
func (r *mockLabelRows) Err() error                                   { return nil }
func (r *mockLabelRows) CommandTag() pgconn.CommandTag                { return pgconn.CommandTag{} }
func (r *mockLabelRows) FieldDescriptions() []pgconn.FieldDescription { return nil }
func (r *mockLabelRows) Values() ([]any, error)                       { return nil, nil }
func (r *mockLabelRows) RawValues() [][]byte                          { return nil }
func (r *mockLabelRows) Conn() *pgx.Conn                              { return nil }
func (r *mockLabelRows) Next() bool                                   { r.index++; return r.index <= len(r.labels) }
func (r *mockLabelRows) Scan(dest ...any) error {
	return (&mockLabelRow{label: r.labels[r.index-1]}).Scan(dest...)
}

// — Upsert unit tests —

func Test_Upsert_Success_ReturnsDomainLabel(t *testing.T) {
	labelID := uuid.New()
	workspaceID := uuid.New()
	now := time.Now().UTC()

	returnedRow := trackerdb.Label{
		ID:          labelID,
		WorkspaceID: workspaceID,
		Name:        "backend",
		CreatedAt:   pgtype.Timestamptz{Time: now, Valid: true},
	}

	mockDatabase := new(mockDBTX)
	mockDatabase.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
		Return(&mockLabelRow{label: returnedRow})

	repository := &LabelRepository{db: mockDatabase}

	actual, err := repository.Upsert(context.Background(), issuedomain.Label{ID: labelID, Name: "backend"})

	require.NoError(t, err)
	assert.Equal(t, labelID, actual.ID)
	assert.Equal(t, "backend", actual.Name)
	mockDatabase.AssertExpectations(t)
}

func Test_Upsert_QueryError_ReturnsWrappedError(t *testing.T) {
	dbErr := errors.New("connection refused")

	mockDatabase := new(mockDBTX)
	mockDatabase.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
		Return(&mockLabelRow{err: dbErr})

	repository := &LabelRepository{db: mockDatabase}

	_, err := repository.Upsert(context.Background(), issuedomain.Label{ID: uuid.New(), Name: "backend"})

	require.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
	assert.Contains(t, err.Error(), "upsert label")
	mockDatabase.AssertExpectations(t)
}

// — ListByIDs unit tests —

func Test_ListByIDs_Success_ReturnsDomainLabels(t *testing.T) {
	now := time.Now().UTC()
	labelID1, labelID2 := uuid.New(), uuid.New()
	workspaceID := uuid.New()

	returnedRows := []trackerdb.Label{
		{ID: labelID1, WorkspaceID: workspaceID, Name: "backend", CreatedAt: pgtype.Timestamptz{Time: now, Valid: true}},
		{ID: labelID2, WorkspaceID: workspaceID, Name: "frontend", CreatedAt: pgtype.Timestamptz{Time: now, Valid: true}},
	}

	mockDatabase := new(mockDBTX)
	mockDatabase.On("Query", mock.Anything, mock.Anything, mock.Anything).
		Return(&mockLabelRows{labels: returnedRows}, nil)

	repository := &LabelRepository{db: mockDatabase}

	actual, err := repository.ListByIDs(context.Background(), []uuid.UUID{labelID1, labelID2})

	require.NoError(t, err)
	require.Len(t, actual, 2)
	assert.Equal(t, labelID1, actual[0].ID)
	assert.Equal(t, "backend", actual[0].Name)
	assert.Equal(t, labelID2, actual[1].ID)
	assert.Equal(t, "frontend", actual[1].Name)
	mockDatabase.AssertExpectations(t)
}

func Test_ListByIDs_EmptyResult_ReturnsEmptySlice(t *testing.T) {
	mockDatabase := new(mockDBTX)
	mockDatabase.On("Query", mock.Anything, mock.Anything, mock.Anything).
		Return(&mockLabelRows{}, nil)

	repository := &LabelRepository{db: mockDatabase}

	actual, err := repository.ListByIDs(context.Background(), []uuid.UUID{uuid.New()})

	require.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Empty(t, actual)
	mockDatabase.AssertExpectations(t)
}

func Test_ListByIDs_QueryError_ReturnsWrappedError(t *testing.T) {
	dbErr := errors.New("connection refused")

	mockDatabase := new(mockDBTX)
	mockDatabase.On("Query", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, dbErr)

	repository := &LabelRepository{db: mockDatabase}

	_, err := repository.ListByIDs(context.Background(), []uuid.UUID{uuid.New()})

	require.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
	assert.Contains(t, err.Error(), "list labels by ids")
	mockDatabase.AssertExpectations(t)
}

// — SearchByName unit tests —

func Test_SearchByName_Success_ReturnsDomainLabels(t *testing.T) {
	now := time.Now().UTC()
	labelID := uuid.New()
	workspaceID := uuid.New()

	returnedRows := []trackerdb.Label{
		{ID: labelID, WorkspaceID: workspaceID, Name: "backend", CreatedAt: pgtype.Timestamptz{Time: now, Valid: true}},
	}

	mockDatabase := new(mockDBTX)
	mockDatabase.On("Query", mock.Anything, mock.Anything, mock.Anything).
		Return(&mockLabelRows{labels: returnedRows}, nil)

	repository := &LabelRepository{db: mockDatabase}

	actual, err := repository.SearchByName(context.Background(), "back")

	require.NoError(t, err)
	require.Len(t, actual, 1)
	assert.Equal(t, labelID, actual[0].ID)
	assert.Equal(t, "backend", actual[0].Name)
	mockDatabase.AssertExpectations(t)
}

func Test_SearchByName_EmptyResult_ReturnsEmptySlice(t *testing.T) {
	mockDatabase := new(mockDBTX)
	mockDatabase.On("Query", mock.Anything, mock.Anything, mock.Anything).
		Return(&mockLabelRows{}, nil)

	repository := &LabelRepository{db: mockDatabase}

	actual, err := repository.SearchByName(context.Background(), "nonexistent")

	require.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Empty(t, actual)
	mockDatabase.AssertExpectations(t)
}

func Test_SearchByName_QueryError_ReturnsWrappedError(t *testing.T) {
	dbErr := errors.New("connection refused")

	mockDatabase := new(mockDBTX)
	mockDatabase.On("Query", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, dbErr)

	repository := &LabelRepository{db: mockDatabase}

	_, err := repository.SearchByName(context.Background(), "back")

	require.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
	assert.Contains(t, err.Error(), "search labels by name")
	mockDatabase.AssertExpectations(t)
}
