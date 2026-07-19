package controllers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"myGymPal/models"
)

// blockWithLinks builds one block's sort-ordered views from raw linked_to_next flags.
func blockWithLinks(links ...bool) []*models.SessionExerciseView {
	views := make([]*models.SessionExerciseView, len(links))
	for i, linked := range links {
		views[i] = &models.SessionExerciseView{
			Exercise: &models.SessionExercise{
				ID:           int64(i + 1),
				Block:        "main",
				SortOrder:    i,
				LinkedToNext: linked,
			},
		}
	}
	return views
}

func linkedFlags(views []*models.SessionExerciseView) []bool {
	out := make([]bool, len(views))
	for i, v := range views {
		out[i] = v.SupersetLinked
	}
	return out
}

func labels(views []*models.SessionExerciseView) []string {
	out := make([]string, len(views))
	for i, v := range views {
		out[i] = v.SupersetLabel
	}
	return out
}

func TestComputeSupersetRuns_WorkedExample(t *testing.T) {
	// bench, row, face pull, squat, curl, tricep ext, plank — the plank's link is
	// stale: it is last in the block, so it has nothing to flow into.
	views := blockWithLinks(true, true, false, false, true, false, true)

	computeSupersetRuns(views)

	assert.Equal(t, []bool{true, true, false, false, true, false, false}, linkedFlags(views))
	assert.Equal(t, []string{"A1", "A2", "A3", "", "B1", "B2", ""}, labels(views))
}

func TestComputeSupersetRuns_StaleLinkOnLastExerciseIsIgnored(t *testing.T) {
	views := blockWithLinks(false, true)

	computeSupersetRuns(views)

	// The rest timer must fire after the last exercise even though its column says
	// otherwise — there is no next card to flow into.
	assert.False(t, views[1].SupersetLinked)
	assert.Empty(t, views[1].SupersetLabel)
}

func TestComputeSupersetRuns_CapsRunAtFourMembers(t *testing.T) {
	// Six exercises all claiming a link: only the first four may form a run.
	views := blockWithLinks(true, true, true, true, true, true)

	computeSupersetRuns(views)

	assert.Equal(t, []bool{true, true, true, false, true, false}, linkedFlags(views))
	// The 4th member's link is forced off, so rest fires after it and a new run starts.
	assert.Equal(t, []string{"A1", "A2", "A3", "A4", "B1", "B2"}, labels(views))
}

func TestComputeSupersetRuns_SoloExercisesUnlabelled(t *testing.T) {
	views := blockWithLinks(false, false, false)

	computeSupersetRuns(views)

	assert.Equal(t, []bool{false, false, false}, linkedFlags(views))
	assert.Equal(t, []string{"", "", ""}, labels(views))
}

func TestComputeSupersetRuns_EmptyBlock(t *testing.T) {
	assert.NotPanics(t, func() { computeSupersetRuns(nil) })
}

func TestGroupSessionExercises_ComputesRunsPerBlock(t *testing.T) {
	mainA := &models.SessionExerciseView{Exercise: &models.SessionExercise{ID: 1, Block: "main", LinkedToNext: true}}
	mainB := &models.SessionExerciseView{Exercise: &models.SessionExercise{ID: 2, Block: "main"}}
	absA := &models.SessionExerciseView{Exercise: &models.SessionExercise{ID: 3, Block: "abs", LinkedToNext: true}}
	absB := &models.SessionExerciseView{Exercise: &models.SessionExercise{ID: 4, Block: "abs"}}

	blocks := groupSessionExercises([]*models.SessionExerciseView{mainA, mainB, absA, absB})
	require.NotEmpty(t, blocks)

	// Runs are lettered per block, so each block starts again at A.
	assert.True(t, mainA.SupersetLinked)
	assert.Equal(t, "A1", mainA.SupersetLabel)
	assert.Equal(t, "A2", mainB.SupersetLabel)
	assert.True(t, absA.SupersetLinked)
	assert.Equal(t, "A1", absA.SupersetLabel)
	assert.Equal(t, "A2", absB.SupersetLabel)
}

func TestGroupSessionExercises_LinkDoesNotCrossBlocks(t *testing.T) {
	// The only exercise in the main block claims a link; the abs block below it is a
	// different block, so the link is ignored rather than chaining across.
	main := &models.SessionExerciseView{Exercise: &models.SessionExercise{ID: 1, Block: "main", LinkedToNext: true}}
	abs := &models.SessionExerciseView{Exercise: &models.SessionExercise{ID: 2, Block: "abs"}}

	groupSessionExercises([]*models.SessionExerciseView{main, abs})

	assert.False(t, main.SupersetLinked)
	assert.Empty(t, main.SupersetLabel)
}

func TestGroupSessionExercises_UntypedBlockDefaultsToMain(t *testing.T) {
	// An exercise with no block falls into "main" and can still be linked.
	first := &models.SessionExerciseView{Exercise: &models.SessionExercise{ID: 1, LinkedToNext: true}}
	second := &models.SessionExerciseView{Exercise: &models.SessionExercise{ID: 2}}

	groupSessionExercises([]*models.SessionExerciseView{first, second})

	assert.True(t, first.SupersetLinked)
	assert.Equal(t, "A1", first.SupersetLabel)
	assert.Equal(t, "A2", second.SupersetLabel)
}

func TestSupersetRunSize(t *testing.T) {
	// A full run of four: bench→row→face→curl, with curl not linking onward.
	full := blockWithLinks(true, true, true, false, false)
	// Turning on the 4th member's link would pull in a 5th.
	assert.Equal(t, 5, supersetRunSize(full, 3))

	// A lone pair.
	pair := blockWithLinks(false, false, false)
	assert.Equal(t, 2, supersetRunSize(pair, 0))
}
