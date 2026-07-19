package models

import (
	"testing"

	"github.com/beego/beego/v2/client/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// exerciseByName pulls a single exercise row out of a template by name, so the
// assertions below do not depend on sort order.
func exerciseByName(t *testing.T, exercises []*TemplateExercise, name string) *TemplateExercise {
	t.Helper()
	for _, ex := range exercises {
		if ex.Name == name {
			return ex
		}
	}
	t.Fatalf("exercise %q not found in template", name)
	return nil
}

// looseExercise builds an exercise input that is not part of any circuit. The
// zero value of CircuitIndex is 0, which means "the first circuit" — never "no
// circuit" — so every non-circuit input has to say NoCircuit explicitly.
func looseExercise(name string, sortOrder int) TemplateExerciseInput {
	return TemplateExerciseInput{Name: name, Block: "main", SortOrder: sortOrder, CircuitIndex: NoCircuit}
}

// UpdateTemplate deletes and re-inserts every exercise row, so any field it
// forgets to copy is silently cleared on every edit. Both flags are asserted:
// a test that only checked the time-based exercise would also pass if the
// update forced every row to true.
func TestUpdateTemplate_PreservesExerciseType(t *testing.T) {
	exercises := []TemplateExerciseInput{
		{Name: "test_plank", IsTimeBased: true, Block: "abs", SortOrder: 0, CircuitIndex: NoCircuit},
		{Name: "test_bench_press", Block: "main", SortOrder: 1, CircuitIndex: NoCircuit},
	}

	tmpl, err := CreateTemplate("test_upper_a", "chest", nil, exercises)
	require.NoError(t, err)
	t.Cleanup(func() { DeleteTemplate(tmpl.ID) })

	_, err = UpdateTemplate(tmpl.ID, "test_upper_a_renamed", "chest", nil, exercises)
	require.NoError(t, err)

	_, saved, err := GetTemplateByID(tmpl.ID)
	require.NoError(t, err)
	require.Len(t, saved, 2)

	plank := exerciseByName(t, saved, "test_plank")
	assert.True(t, plank.IsTimeBased, "time-based flag must survive an update")

	bench := exerciseByName(t, saved, "test_bench_press")
	assert.False(t, bench.IsTimeBased, "a weighted exercise must not become time-based")
}

// --- Circuits ---

// stretchCircuit is the "Morning Stretch" circuit from the feature spec: one
// round, a 5s transition, three timed exercises.
func stretchCircuit() ([]TemplateCircuitInput, []TemplateExerciseInput) {
	circuits := []TemplateCircuitInput{
		{Name: "test_morning_stretch", Rounds: 1, TransitionSeconds: 5, SortOrder: 0},
	}
	exercises := []TemplateExerciseInput{
		{Name: "test_shoulder_stretch", IsTimeBased: true, Block: "stretch", SortOrder: 0, CircuitIndex: 0, WorkSeconds: 30},
		{Name: "test_hip_flexor", IsTimeBased: true, Block: "stretch", SortOrder: 1, CircuitIndex: 0, WorkSeconds: 45},
		{Name: "test_hamstring", IsTimeBased: true, Block: "stretch", SortOrder: 2, CircuitIndex: 0, WorkSeconds: 30},
	}
	return circuits, exercises
}

func TestCreateTemplate_CircuitRoundTrips(t *testing.T) {
	circuits, exercises := stretchCircuit()

	tmpl, err := CreateTemplate("test_stretch_template", "mobility", circuits, exercises)
	require.NoError(t, err)
	t.Cleanup(func() { DeleteTemplate(tmpl.ID) })

	saved, err := GetTemplateCircuits(tmpl.ID)
	require.NoError(t, err)
	require.Len(t, saved, 1)
	assert.Equal(t, "test_morning_stretch", saved[0].Name)
	assert.Equal(t, 1, saved[0].Rounds)
	assert.Equal(t, 5, saved[0].TransitionSeconds)

	_, savedEx, err := GetTemplateByID(tmpl.ID)
	require.NoError(t, err)
	require.Len(t, savedEx, 3)

	// Every exercise points at the circuit that was created in the same submit,
	// and each keeps its own work duration.
	wantSeconds := map[string]int{
		"test_shoulder_stretch": 30,
		"test_hip_flexor":       45,
		"test_hamstring":        30,
	}
	for name, want := range wantSeconds {
		ex := exerciseByName(t, savedEx, name)
		require.NotNil(t, ex.CircuitID, "%s should be in a circuit", name)
		assert.Equal(t, saved[0].ID, *ex.CircuitID, "%s should point at the stretch circuit", name)
		assert.Equal(t, want, ex.WorkSeconds, "%s work seconds", name)
	}
}

// The template body is deleted and re-inserted on every update, so this is the
// shape the is_time_based bug took. Only the template's name changes; everything
// else must come back identical.
func TestUpdateTemplate_PreservesCircuitsAndWorkSeconds(t *testing.T) {
	circuits, exercises := stretchCircuit()

	tmpl, err := CreateTemplate("test_stretch_template", "mobility", circuits, exercises)
	require.NoError(t, err)
	t.Cleanup(func() { DeleteTemplate(tmpl.ID) })

	_, err = UpdateTemplate(tmpl.ID, "test_stretch_renamed", "mobility", circuits, exercises)
	require.NoError(t, err)

	saved, err := GetTemplateCircuits(tmpl.ID)
	require.NoError(t, err)
	require.Len(t, saved, 1, "the circuit must survive an update, not be dropped or duplicated")
	assert.Equal(t, "test_morning_stretch", saved[0].Name)
	assert.Equal(t, 1, saved[0].Rounds, "round count must survive an update")
	assert.Equal(t, 5, saved[0].TransitionSeconds, "transition must survive an update")

	_, savedEx, err := GetTemplateByID(tmpl.ID)
	require.NoError(t, err)
	require.Len(t, savedEx, 3)

	hip := exerciseByName(t, savedEx, "test_hip_flexor")
	require.NotNil(t, hip.CircuitID, "circuit membership must survive an update")
	assert.Equal(t, saved[0].ID, *hip.CircuitID, "membership must point at the re-created circuit")
	assert.Equal(t, 45, hip.WorkSeconds, "work seconds must survive an update")
	assert.True(t, hip.IsTimeBased, "the time-based flag must still survive an update")
}

// The positive control for the circuit work: a template with no circuits must
// still save and read back exactly as it did before circuits existed.
func TestCreateTemplate_NoCircuitLeavesCircuitIDNull(t *testing.T) {
	exercises := []TemplateExerciseInput{looseExercise("test_squat", 0)}

	tmpl, err := CreateTemplate("test_plain_template", "legs", nil, exercises)
	require.NoError(t, err)
	t.Cleanup(func() { DeleteTemplate(tmpl.ID) })

	_, saved, err := GetTemplateByID(tmpl.ID)
	require.NoError(t, err)
	require.Len(t, saved, 1)
	assert.Nil(t, saved[0].CircuitID, "an exercise outside a circuit must have circuit_id NULL")
	assert.Equal(t, 0, saved[0].WorkSeconds)

	circuits, err := GetTemplateCircuits(tmpl.ID)
	require.NoError(t, err)
	assert.Empty(t, circuits)
}

// A circuit index that points at nothing is a bug in the form. Writing NULL
// instead would quietly drop the exercise out of its circuit.
func TestCreateTemplate_CircuitIndexOutOfRangeIsError(t *testing.T) {
	circuits := []TemplateCircuitInput{
		{Name: "test_only_circuit", Rounds: 2, TransitionSeconds: 10, SortOrder: 0},
	}
	exercises := []TemplateExerciseInput{
		{Name: "test_burpees", SortOrder: 0, CircuitIndex: 1, WorkSeconds: 40}, // only index 0 exists
	}

	tmpl, err := CreateTemplate("test_bad_index", "hiit", circuits, exercises)
	require.Error(t, err, "an out-of-range circuit index must be an error, not a silent NULL")
	assert.Nil(t, tmpl)
	assert.Contains(t, err.Error(), "does not exist")
}

// The same guard has to hold on the update path, which is a separate code path
// into the same insert.
func TestUpdateTemplate_CircuitIndexOutOfRangeIsError(t *testing.T) {
	circuits, exercises := stretchCircuit()

	tmpl, err := CreateTemplate("test_stretch_template", "mobility", circuits, exercises)
	require.NoError(t, err)
	t.Cleanup(func() { DeleteTemplate(tmpl.ID) })

	bad := []TemplateExerciseInput{
		{Name: "test_burpees", SortOrder: 0, CircuitIndex: 3, WorkSeconds: 40},
	}
	_, err = UpdateTemplate(tmpl.ID, "test_stretch_template", "mobility", circuits, bad)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")

	// The failed update must not have destroyed the body it was replacing.
	saved, err := GetTemplateCircuits(tmpl.ID)
	require.NoError(t, err)
	assert.Len(t, saved, 1, "a rejected update must leave the existing circuit intact")
}

func TestCreateTemplate_RejectsCircuitWithoutName(t *testing.T) {
	circuits := []TemplateCircuitInput{{Name: "  ", Rounds: 1, SortOrder: 0}}
	exercises := []TemplateExerciseInput{
		{Name: "test_burpees", SortOrder: 0, CircuitIndex: 0, WorkSeconds: 40},
	}

	_, err := CreateTemplate("test_unnamed_circuit", "hiit", circuits, exercises)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "circuit name is required")
}

func TestDeleteTemplate_CascadesCircuits(t *testing.T) {
	circuits, exercises := stretchCircuit()

	tmpl, err := CreateTemplate("test_stretch_template", "mobility", circuits, exercises)
	require.NoError(t, err)

	saved, err := GetTemplateCircuits(tmpl.ID)
	require.NoError(t, err)
	require.Len(t, saved, 1)

	require.NoError(t, DeleteTemplate(tmpl.ID))

	// Query the table directly: GetTemplateCircuits filters by template_id, which
	// would return empty even if the rows were orphaned rather than deleted.
	var count int
	err = orm.NewOrm().Raw("SELECT COUNT(*) FROM template_circuits WHERE id = ?", saved[0].ID).QueryRow(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count, "deleting a template must cascade its circuits away")
}

func TestValidRoundsAndSeconds(t *testing.T) {
	assert.Equal(t, 1, ValidRounds(0), "rounds below one coerce to one")
	assert.Equal(t, 1, ValidRounds(-4))
	assert.Equal(t, 4, ValidRounds(4))

	assert.Equal(t, 0, ValidSeconds(-1), "negative seconds coerce to zero")
	assert.Equal(t, 30, ValidSeconds(30))
}
