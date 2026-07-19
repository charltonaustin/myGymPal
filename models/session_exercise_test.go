package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testSessionExercise creates a user, program, session and one session exercise,
// registering cleanup immediately so a failing assertion cannot leak rows into the
// next run. Deleting the user cascades to everything below it.
func testSessionExercise(t *testing.T, username string) *SessionExercise {
	t.Helper()
	_ = DeleteUserByUsername(username)
	u, err := CreateUser(username, "password123", "lb")
	require.NoError(t, err)
	t.Cleanup(func() { DeleteUserByUsername(username) })

	p, err := CreateProgram(u.ID, "test_link_program", testStartDate, 1, 8, 4, 10, 12, 3)
	require.NoError(t, err)

	s, err := CreateSession(p.ID, u.ID, 1, 1, 1, false, testStartDate)
	require.NoError(t, err)

	ex, err := CreateSessionExercise(s.ID, "test_bench_press", false, 100, "lb", 10, "main", false, 0)
	require.NoError(t, err)
	return ex
}

// readBackLink re-reads the exercise through the view the session page uses, so the
// test proves the column survives the round trip the app actually makes.
func readBackLink(t *testing.T, ex *SessionExercise) bool {
	t.Helper()
	views, err := GetSessionExercisesWithSets(ex.SessionID)
	require.NoError(t, err)
	for _, v := range views {
		if v.Exercise.ID == ex.ID {
			return v.Exercise.LinkedToNext
		}
	}
	t.Fatalf("exercise %d not found in session %d", ex.ID, ex.SessionID)
	return false
}

func TestUpdateSessionExerciseLink_DefaultsToFalse(t *testing.T) {
	ex := testSessionExercise(t, "test_link_default")

	// Proves migration 000034 landed with its NOT NULL DEFAULT FALSE.
	assert.False(t, ex.LinkedToNext)
	assert.False(t, readBackLink(t, ex))
}

func TestUpdateSessionExerciseLink_PersistsTrue(t *testing.T) {
	ex := testSessionExercise(t, "test_link_on")

	require.NoError(t, UpdateSessionExerciseLink(ex.ID, true))

	assert.True(t, readBackLink(t, ex))
}

func TestUpdateSessionExerciseLink_PersistsFalse(t *testing.T) {
	ex := testSessionExercise(t, "test_link_off")
	require.NoError(t, UpdateSessionExerciseLink(ex.ID, true))
	require.True(t, readBackLink(t, ex))

	require.NoError(t, UpdateSessionExerciseLink(ex.ID, false))

	assert.False(t, readBackLink(t, ex))
}

func TestUpdateSessionExerciseLink_MissingIDDoesNotPanic(t *testing.T) {
	// An UPDATE that matches no row is a no-op, not an error.
	assert.NotPanics(t, func() {
		assert.NoError(t, UpdateSessionExerciseLink(-1, true))
	})
}
