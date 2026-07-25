package models

import (
	"testing"

	"github.com/beego/beego/v2/client/orm"
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

	ex, err := CreateSessionExercise(SessionExerciseInput{
		SessionID:  s.ID,
		Name:       "test_bench_press",
		GoalWeight: 100,
		WeightUnit: "lb",
		GoalReps:   10,
		Block:      "main",
	})
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

// --- Circuits ---

// testSessionFor creates a user, program and session under a test-prefixed
// username, registering cleanup immediately. Deleting the user cascades to
// everything below it, so a failing assertion cannot leak rows into the next run.
func testSessionFor(t *testing.T, username string) (*User, *Session) {
	t.Helper()
	_ = DeleteUserByUsername(username)
	u, err := CreateUser(username, "password123", "lb")
	require.NoError(t, err)
	t.Cleanup(func() { DeleteUserByUsername(username) })

	p, err := CreateProgram(u.ID, "test_circuit_program", testStartDate, 1, 8, 4, 10, 12, 3)
	require.NoError(t, err)

	s, err := CreateSession(p.ID, u.ID, 1, 1, 1, false, testStartDate)
	require.NoError(t, err)
	return u, s
}

// morningStretchBody is the fixture the runner is built around: one circuit, one
// round, a 5s transition, and three stretches held for different lengths. The
// circuit ids are the *template* ids the copy remaps.
func morningStretchBody() ([]SessionCircuitInput, []SessionExerciseInput) {
	const tmplCircuitID int64 = 4242
	circuitID := tmplCircuitID
	circuits := []SessionCircuitInput{{
		TemplateCircuitID: tmplCircuitID,
		Name:              "test_morning_stretch",
		Rounds:            1,
		TransitionSeconds: 5,
		SortOrder:         0,
	}}
	exercises := []SessionExerciseInput{
		{Name: "test_shoulder_stretch", Block: "stretch", IsTimeBased: true, CircuitID: &circuitID, WorkSeconds: 30},
		{Name: "test_hip_flexor", Block: "stretch", IsTimeBased: true, CircuitID: &circuitID, WorkSeconds: 45},
		{Name: "test_hamstring", Block: "stretch", IsTimeBased: true, CircuitID: &circuitID, WorkSeconds: 30},
	}
	return circuits, exercises
}

// exercisesByName re-reads a session through the view the session page uses, so
// the tests prove the columns survive the round trip the app actually makes.
func exercisesByName(t *testing.T, sessionID int64) map[string]*SessionExercise {
	t.Helper()
	views, err := GetSessionExercisesWithSets(sessionID)
	require.NoError(t, err)
	byName := make(map[string]*SessionExercise, len(views))
	for _, v := range views {
		byName[v.Exercise.Name] = v.Exercise
	}
	return byName
}

func TestCreateSessionBody_CopiesCircuitFromTemplate(t *testing.T) {
	_, s := testSessionFor(t, "test_circuit_copy")
	circuits, exercises := morningStretchBody()

	require.NoError(t, CreateSessionBody(s.ID, circuits, exercises))

	got, err := GetSessionCircuits(s.ID)
	require.NoError(t, err)
	require.Len(t, got, 1, "the template's circuit must be copied into the session")
	assert.Equal(t, "test_morning_stretch", got[0].Name)
	assert.Equal(t, 1, got[0].Rounds, "rounds must survive the copy")
	assert.Equal(t, 5, got[0].TransitionSeconds, "transition must survive the copy")
	assert.Equal(t, s.ID, got[0].SessionID)

	byName := exercisesByName(t, s.ID)
	require.Len(t, byName, 3)
	for name, wantSecs := range map[string]int{
		"test_shoulder_stretch": 30,
		"test_hip_flexor":       45,
		"test_hamstring":        30,
	} {
		ex := byName[name]
		require.NotNil(t, ex, "%s must be copied into the session", name)
		require.NotNil(t, ex.CircuitID, "%s must be filed under the circuit", name)
		assert.Equal(t, got[0].ID, *ex.CircuitID, "%s must point at the *session* circuit, not the template one", name)
		assert.Equal(t, wantSecs, ex.WorkSeconds, "%s work seconds", name)
	}
}

func TestCreateSessionBody_LooseExerciseIsNotInACircuit(t *testing.T) {
	_, s := testSessionFor(t, "test_circuit_loose")
	circuits, exercises := morningStretchBody()
	exercises = append(exercises, SessionExerciseInput{
		Name: "test_bench_press", Block: "main", WeightUnit: "lb", GoalSeconds: 0,
	})

	require.NoError(t, CreateSessionBody(s.ID, circuits, exercises))

	bench := exercisesByName(t, s.ID)["test_bench_press"]
	require.NotNil(t, bench)
	assert.Nil(t, bench.CircuitID, "a loose exercise must not be filed under a circuit")
}

// A loose exercise is *allowed* to carry work seconds — the template repository
// stores what it is given rather than forcing 0 — so a non-zero work_seconds must
// never be read as circuit membership. Membership is CircuitID and nothing else.
func TestCreateSessionBody_LooseExerciseWithWorkSecondsStaysLoose(t *testing.T) {
	_, s := testSessionFor(t, "test_circuit_loose_secs")
	circuits, exercises := morningStretchBody()
	exercises = append(exercises, SessionExerciseInput{
		Name: "test_plank", Block: "abs", IsTimeBased: true, GoalSeconds: 60, WorkSeconds: 90,
	})

	require.NoError(t, CreateSessionBody(s.ID, circuits, exercises))

	plank := exercisesByName(t, s.ID)["test_plank"]
	require.NotNil(t, plank)
	assert.Nil(t, plank.CircuitID, "work seconds alone must not put an exercise in a circuit")
	assert.Equal(t, 90, plank.WorkSeconds, "the value is stored as given, not zeroed")
	assert.Equal(t, 60, plank.GoalSeconds, "a loose exercise keeps its own goal seconds")
}

// An exercise naming a circuit that was not submitted alongside it becomes loose.
// Storing a dangling id would fail the foreign key and roll the whole copy back,
// losing the rest of the workout over one bad row.
func TestCreateSessionBody_UnknownCircuitBecomesLoose(t *testing.T) {
	_, s := testSessionFor(t, "test_circuit_unknown")
	stray := int64(999999)
	exercises := []SessionExerciseInput{
		{Name: "test_orphan", Block: "main", CircuitID: &stray, WorkSeconds: 20},
	}

	require.NoError(t, CreateSessionBody(s.ID, nil, exercises))

	orphan := exercisesByName(t, s.ID)["test_orphan"]
	require.NotNil(t, orphan)
	assert.Nil(t, orphan.CircuitID)
}

func TestCreateSessionBody_PreservesTemplateOrder(t *testing.T) {
	_, s := testSessionFor(t, "test_circuit_order")
	circuits, exercises := morningStretchBody()

	require.NoError(t, CreateSessionBody(s.ID, circuits, exercises))

	views, err := GetSessionExercisesWithSets(s.ID)
	require.NoError(t, err)
	var names []string
	for _, v := range views {
		names = append(names, v.Exercise.Name)
	}
	assert.Equal(t, []string{"test_shoulder_stretch", "test_hip_flexor", "test_hamstring"}, names)
}

func TestCreateSessionBody_AppendsAfterExistingExercises(t *testing.T) {
	_, s := testSessionFor(t, "test_circuit_append")
	_, err := CreateSessionExercise(SessionExerciseInput{SessionID: s.ID, Name: "test_warmup", Block: "main", WeightUnit: "lb"})
	require.NoError(t, err)
	circuits, exercises := morningStretchBody()

	require.NoError(t, CreateSessionBody(s.ID, circuits, exercises))

	views, err := GetSessionExercisesWithSets(s.ID)
	require.NoError(t, err)
	require.Len(t, views, 4)
	assert.Equal(t, "test_warmup", views[0].Exercise.Name, "the copy must not renumber over what is already there")
	assert.Equal(t, "test_shoulder_stretch", views[1].Exercise.Name)
}

func TestCreateSessionExercise_DefaultsToNoCircuit(t *testing.T) {
	_, s := testSessionFor(t, "test_circuit_default")

	ex, err := CreateSessionExercise(SessionExerciseInput{SessionID: s.ID, Name: "test_row", Block: "main", WeightUnit: "lb"})
	require.NoError(t, err)

	// Proves migration 000036 landed with a nullable circuit_id and work_seconds
	// defaulting to 0 — the shape every pre-existing row has.
	assert.Nil(t, ex.CircuitID)
	assert.Equal(t, 0, ex.WorkSeconds)
	readBack := exercisesByName(t, s.ID)["test_row"]
	require.NotNil(t, readBack)
	assert.Nil(t, readBack.CircuitID)
	assert.Equal(t, 0, readBack.WorkSeconds)
}

func TestDeleteSession_CascadesCircuits(t *testing.T) {
	u, s := testSessionFor(t, "test_circuit_cascade")
	circuits, exercises := morningStretchBody()
	require.NoError(t, CreateSessionBody(s.ID, circuits, exercises))
	created, err := GetSessionCircuits(s.ID)
	require.NoError(t, err)
	require.Len(t, created, 1)
	circuitID := created[0].ID

	require.NoError(t, DeleteSession(s.ID, u.ID))

	// Queried by circuit id rather than by session, so a row orphaned from its
	// deleted session would still fail this.
	o := orm.NewOrm()
	n, err := o.QueryTable(&SessionCircuit{}).Filter("ID", circuitID).Count()
	require.NoError(t, err)
	assert.Equal(t, int64(0), n, "deleting a session must cascade its circuits away")
}

func TestGetSessionCircuits_ScopedToOneUsersSession(t *testing.T) {
	_, sA := testSessionFor(t, "test_circuit_user_a")
	_, sB := testSessionFor(t, "test_circuit_user_b")
	circuits, exercises := morningStretchBody()
	require.NoError(t, CreateSessionBody(sA.ID, circuits, exercises))

	got, err := GetSessionCircuits(sB.ID)

	require.NoError(t, err)
	assert.Empty(t, got, "user B's session must not see user A's circuits")
}
