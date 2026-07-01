//go:build integration

package tracker_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	issuedomain "github.com/wlindb/issue-tracker/internal/domain/tracker/issue"
	tracker "github.com/wlindb/issue-tracker/internal/infrastructure/tracker"
)

// — LabelRepository integration tests —

func Test_Upsert_NewLabel_ReturnsLabel(t *testing.T) {
	repository := tracker.NewLabelRepository(testPool)
	_, ctx := createTestWorkspace(t)
	labelID := uuid.New()

	actual, err := repository.Upsert(ctx, issuedomain.Label{ID: labelID, Name: "backend"})

	require.NoError(t, err)
	assert.Equal(t, labelID, actual.ID)
	assert.Equal(t, "backend", actual.Name)
}

func Test_Upsert_ExistingID_UpdatesName(t *testing.T) {
	repository := tracker.NewLabelRepository(testPool)
	_, ctx := createTestWorkspace(t)
	labelID := uuid.New()

	_, err := repository.Upsert(ctx, issuedomain.Label{ID: labelID, Name: "original"})
	require.NoError(t, err)

	actual, err := repository.Upsert(ctx, issuedomain.Label{ID: labelID, Name: "updated"})

	require.NoError(t, err)
	assert.Equal(t, labelID, actual.ID)
	assert.Equal(t, "updated", actual.Name)
}

func Test_ListByIDs_ExistingLabels_ReturnsLabels(t *testing.T) {
	repository := tracker.NewLabelRepository(testPool)
	_, ctx := createTestWorkspace(t)

	labelID1 := createTestLabel(t, ctx, "alpha")
	labelID2 := createTestLabel(t, ctx, "beta")

	actual, err := repository.ListByIDs(ctx, []uuid.UUID{labelID1, labelID2})

	require.NoError(t, err)
	require.Len(t, actual, 2)
	ids := []uuid.UUID{actual[0].ID, actual[1].ID}
	assert.Contains(t, ids, labelID1)
	assert.Contains(t, ids, labelID2)
}

func Test_ListByIDs_NonExistentIDs_ReturnsEmptySlice(t *testing.T) {
	repository := tracker.NewLabelRepository(testPool)
	_, ctx := createTestWorkspace(t)

	actual, err := repository.ListByIDs(ctx, []uuid.UUID{uuid.New(), uuid.New()})

	require.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Empty(t, actual)
}

func Test_SearchByName_SubstringMatch_ReturnsLabels(t *testing.T) {
	repository := tracker.NewLabelRepository(testPool)
	_, ctx := createTestWorkspace(t)

	createTestLabel(t, ctx, "backend-service")
	createTestLabel(t, ctx, "backend-api")
	createTestLabel(t, ctx, "frontend")

	actual, err := repository.SearchByName(ctx, "backend")

	require.NoError(t, err)
	require.Len(t, actual, 2)
	for _, label := range actual {
		assert.Contains(t, label.Name, "backend")
	}
}

func Test_SearchByName_NoMatch_ReturnsEmptySlice(t *testing.T) {
	repository := tracker.NewLabelRepository(testPool)
	_, ctx := createTestWorkspace(t)

	createTestLabel(t, ctx, "alpha")

	actual, err := repository.SearchByName(ctx, "zzz-no-match")

	require.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Empty(t, actual)
}

func Test_ListByIDs_CrossWorkspaceIsolation_ReturnsEmpty(t *testing.T) {
	repository := tracker.NewLabelRepository(testPool)
	_, ctxA := createTestWorkspace(t)
	_, ctxB := createTestWorkspace(t)

	labelID := createTestLabel(t, ctxA, "workspace-a-label")

	actual, err := repository.ListByIDs(ctxB, []uuid.UUID{labelID})

	require.NoError(t, err)
	assert.Empty(t, actual)
}
