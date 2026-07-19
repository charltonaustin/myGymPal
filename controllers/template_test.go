package controllers_test

import (
	"errors"
	"fmt"
	"myGymPal/models"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Templates list ---

func TestTemplatesIndex_Unauthenticated(t *testing.T) {
	w := getPath("/templates", nil)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/login", w.Header().Get("Location"))
}

func TestTemplatesIndex_Empty(t *testing.T) {
	t.Cleanup(resetMocks)
	setTemplatesGetAllEmpty()
	cookies := loginAs(t, "tmpl_idx_empty", "lb")

	w := getPath("/templates", cookies)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "No templates yet")
}

func TestTemplatesIndex_ShowsTemplates(t *testing.T) {
	t.Cleanup(resetMocks)
	setTemplatesGetAllWithOne(1, "Upper Body A", "Chest & Shoulders")
	cookies := loginAs(t, "tmpl_idx_list", "lb")

	w := getPath("/templates", cookies)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Upper Body A")
	assert.Contains(t, w.Body.String(), "Chest &amp; Shoulders")
}

// --- New template form ---

func TestTemplatesNew_Unauthenticated(t *testing.T) {
	w := getPath("/templates/new", nil)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/login", w.Header().Get("Location"))
}

func TestTemplatesNew_ShowsForm(t *testing.T) {
	t.Cleanup(resetMocks)
	cookies := loginAs(t, "tmpl_new_form", "lb")

	w := getPath("/templates/new", cookies)
	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, "New Workout Template")
	assert.Contains(t, body, "Create Template")
	assert.Contains(t, body, `action="/templates/new"`)
	assert.Contains(t, body, `name="name"`)
	assert.Contains(t, body, `name="focus"`)
	assert.Contains(t, body, `name="exercise_name_0"`)
	assert.Contains(t, body, `name="is_bodyweight_0"`)
}

// --- Create template ---

func TestTemplatesCreate_Unauthenticated(t *testing.T) {
	w := postForm("/templates/new", url.Values{
		"name":            {"Test"},
		"exercise_count":  {"1"},
		"exercise_name_0": {"Bench Press"},
	}, nil)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/login", w.Header().Get("Location"))
}

func TestTemplatesCreate_Success(t *testing.T) {
	t.Cleanup(resetMocks)
	captureTemplateCreate()
	cookies := loginAs(t, "tmpl_create_ok", "lb")

	w := postForm("/templates/new", url.Values{
		"name":            {"Upper Body A"},
		"focus":           {"Chest & Shoulders"},
		"exercise_count":  {"1"},
		"exercise_name_0": {"Bench Press"},
	}, cookies)

	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, fmt.Sprintf("/templates/%d", testTemplateID), w.Header().Get("Location"))
	assert.Equal(t, "Upper Body A", lastTemplateCreate.name)
	assert.Equal(t, "Chest & Shoulders", lastTemplateCreate.focus)
	assert.Equal(t, 1, lastTemplateCreate.numExercises)
}

func TestTemplatesCreate_BodyweightExercise(t *testing.T) {
	t.Cleanup(resetMocks)
	captureTemplateCreate()
	cookies := loginAs(t, "tmpl_create_bw", "lb")

	w := postForm("/templates/new", url.Values{
		"name":            {"Pull Day"},
		"exercise_count":  {"1"},
		"exercise_name_0": {"Pull-ups"},
		"is_bodyweight_0": {"on"},
	}, cookies)

	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, 1, lastTemplateCreate.numExercises)
}

func TestTemplatesCreate_EmptyName(t *testing.T) {
	t.Cleanup(resetMocks)
	cookies := loginAs(t, "tmpl_create_noname", "lb")

	w := postForm("/templates/new", url.Values{
		"name":            {""},
		"exercise_count":  {"1"},
		"exercise_name_0": {"Bench Press"},
	}, cookies)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "required")
}

func TestTemplatesCreate_NoExercises(t *testing.T) {
	t.Cleanup(resetMocks)
	setTemplateCreateError(errors.New("at least one exercise is required"))
	cookies := loginAs(t, "tmpl_create_noex", "lb")

	w := postForm("/templates/new", url.Values{
		"name":           {"My Template"},
		"exercise_count": {"0"},
	}, cookies)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "exercise")
}

func TestTemplatesCreate_ReentersFormValues(t *testing.T) {
	t.Cleanup(resetMocks)
	setTemplateCreateError(errors.New("exercise name is required"))
	cookies := loginAs(t, "tmpl_create_reenter", "lb")

	w := postForm("/templates/new", url.Values{
		"name":            {"Sticky Template"},
		"focus":           {"Legs"},
		"exercise_count":  {"1"},
		"exercise_name_0": {"Squat"},
	}, cookies)

	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, "Sticky Template")
	assert.Contains(t, body, "Squat")
	// The error path renders through the shared form partial too, so it needs
	// the same chrome keys as the initial GET.
	assert.Contains(t, body, "Create Template")
	assert.Contains(t, body, `action="/templates/new"`)
}

// --- Edit template form ---

func TestTemplatesEdit_Unauthenticated(t *testing.T) {
	w := getPath("/templates/10/edit", nil)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/login", w.Header().Get("Location"))
}

func TestTemplatesEdit_NotFound(t *testing.T) {
	t.Cleanup(resetMocks)
	setTemplateGetByIDError(errors.New("not found"))
	cookies := loginAs(t, "tmpl_edit_404", "lb")

	w := getPath("/templates/10/edit", cookies)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/templates", w.Header().Get("Location"))
}

func TestTemplatesEdit_ShowsPrefilledForm(t *testing.T) {
	t.Cleanup(resetMocks)
	setTemplateGetByID(testTemplateID, "Upper Body A", "Chest", 2)
	cookies := loginAs(t, "tmpl_edit_form", "lb")

	w := getPath("/templates/10/edit", cookies)
	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, "Edit Template")
	assert.Contains(t, body, "Save Changes")
	assert.Contains(t, body, fmt.Sprintf(`action="/templates/%d"`, testTemplateID))
	assert.Contains(t, body, `value="Upper Body A"`)
	assert.Contains(t, body, `value="Chest"`)
	assert.Contains(t, body, `name="exercise_name_1"`)
}

// --- Update template ---

func TestTemplatesUpdate_Unauthenticated(t *testing.T) {
	w := postForm("/templates/10", url.Values{"name": {"Test"}}, nil)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/login", w.Header().Get("Location"))
}

func TestTemplatesUpdate_Success(t *testing.T) {
	t.Cleanup(resetMocks)
	setTemplateGetByID(testTemplateID, "Upper Body A", "Chest", 1)
	captureTemplateUpdate()
	cookies := loginAs(t, "tmpl_update_ok", "lb")

	w := postForm("/templates/10", url.Values{
		"name":            {"Upper Body B"},
		"focus":           {"Back"},
		"exercise_count":  {"2"},
		"exercise_name_0": {"Plank"},
		"is_time_based_0": {"on"},
		"exercise_name_1": {"Bench Press"},
		"block_1":         {"main"},
	}, cookies)

	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, fmt.Sprintf("/templates/%d", testTemplateID), w.Header().Get("Location"))
	assert.Equal(t, "Upper Body B", lastTemplateUpdate.name)
	assert.Equal(t, "Back", lastTemplateUpdate.focus)

	// The time-based flag must survive the trip through the controller as well
	// as the model — the form posts it, so the repository must receive it.
	assert.Len(t, lastTemplateUpdate.exercises, 2)
	assert.True(t, lastTemplateUpdate.exercises[0].IsTimeBased)
	assert.False(t, lastTemplateUpdate.exercises[1].IsTimeBased)
}

func TestTemplatesUpdate_EmptyName(t *testing.T) {
	t.Cleanup(resetMocks)
	setTemplateGetByID(testTemplateID, "Upper Body A", "Chest", 1)
	cookies := loginAs(t, "tmpl_update_noname", "lb")

	w := postForm("/templates/10", url.Values{
		"name":            {""},
		"exercise_count":  {"1"},
		"exercise_name_0": {"Bench Press"},
	}, cookies)

	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, "required")
	// Same chrome check as the create error path: a c.Data key missed here
	// renders an empty submit button instead of failing.
	assert.Contains(t, body, "Save Changes")
	assert.Contains(t, body, fmt.Sprintf(`action="/templates/%d"`, testTemplateID))
	assert.Contains(t, body, "Bench Press")
}

func TestTemplatesUpdate_RepositoryError(t *testing.T) {
	t.Cleanup(resetMocks)
	setTemplateGetByID(testTemplateID, "Upper Body A", "Chest", 1)
	setTemplateUpdateError(errors.New("exercise name is required"))
	cookies := loginAs(t, "tmpl_update_err", "lb")

	w := postForm("/templates/10", url.Values{
		"name":            {"Upper Body A"},
		"exercise_count":  {"1"},
		"exercise_name_0": {"Bench Press"},
	}, cookies)

	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, "exercise name is required")
	assert.Contains(t, body, "Save Changes")
}

// --- Template show ---

func TestTemplatesShow_Unauthenticated(t *testing.T) {
	w := getPath("/templates/1", nil)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/login", w.Header().Get("Location"))
}

func TestTemplatesShow_NotFound(t *testing.T) {
	t.Cleanup(resetMocks)
	setTemplateGetByIDError(errors.New("not found"))
	cookies := loginAs(t, "tmpl_show_notfound", "lb")

	w := getPath("/templates/1", cookies)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/templates", w.Header().Get("Location"))
}

func TestTemplatesShow_ShowsTemplate(t *testing.T) {
	t.Cleanup(resetMocks)
	setTemplateGetByID(1, "Upper Body A", "Chest & Shoulders", 3)
	cookies := loginAs(t, "tmpl_show_ok", "lb")

	w := getPath("/templates/1", cookies)
	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, "Upper Body A")
	assert.Contains(t, body, "Exercise 1")
	assert.Contains(t, body, "Exercise 2")
	assert.Contains(t, body, "Exercise 3")
}

// --- Delete template ---

func TestTemplateDelete_Unauthenticated(t *testing.T) {
	w := postForm(fmt.Sprintf("/templates/%d/delete", testTemplateID), url.Values{}, nil)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/login", w.Header().Get("Location"))
}

func TestTemplateDelete_Error(t *testing.T) {
	t.Cleanup(resetMocks)
	mockTemplates.DeleteFn = func(id int64) error {
		return errors.New("not found")
	}
	cookies := loginAs(t, "tmpl_delete_err", "lb")

	w := postForm(fmt.Sprintf("/templates/%d/delete", testTemplateID), url.Values{}, cookies)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/templates", w.Header().Get("Location"))
}

func TestTemplateDelete_Success(t *testing.T) {
	t.Cleanup(resetMocks)
	cookies := loginAs(t, "tmpl_delete_ok", "lb")

	w := postForm(fmt.Sprintf("/templates/%d/delete", testTemplateID), url.Values{}, cookies)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/templates", w.Header().Get("Location"))
}

// --- Circuits ---

// circuitID is a stable id for the circuit the mocks hand back.
const circuitID = int64(77)

// stretchCircuitForm is the "Morning Stretch" circuit as the browser submits it:
// one circuit, three exercises inside it, each with its own work duration.
func stretchCircuitForm() url.Values {
	return url.Values{
		"name":                 {"Mobility"},
		"circuit_count":        {"1"},
		"circuit_name_0":       {"Morning Stretch"},
		"circuit_rounds_0":     {"1"},
		"circuit_transition_0": {"5"},
		"exercise_count":       {"3"},
		"exercise_name_0":      {"Shoulder Stretch"},
		"is_time_based_0":      {"on"},
		"work_seconds_0":       {"30"},
		"circuit_index_0":      {"0"},
		"exercise_name_1":      {"Hip Flexor"},
		"is_time_based_1":      {"on"},
		"work_seconds_1":       {"45"},
		"circuit_index_1":      {"0"},
		"exercise_name_2":      {"Hamstring"},
		"is_time_based_2":      {"on"},
		"work_seconds_2":       {"30"},
		"circuit_index_2":      {"0"},
	}
}

func TestTemplatesCreate_WithCircuit(t *testing.T) {
	t.Cleanup(resetMocks)
	captureTemplateCreate()
	cookies := loginAs(t, "tmpl_create_circuit", "lb")

	w := postForm("/templates/new", stretchCircuitForm(), cookies)
	assert.Equal(t, http.StatusFound, w.Code)

	require.Len(t, lastTemplateCreate.circuits, 1)
	circuit := lastTemplateCreate.circuits[0]
	assert.Equal(t, "Morning Stretch", circuit.Name)
	assert.Equal(t, 1, circuit.Rounds)
	assert.Equal(t, 5, circuit.TransitionSeconds)

	require.Len(t, lastTemplateCreate.exercises, 3)
	for i, wantSeconds := range []int{30, 45, 30} {
		ex := lastTemplateCreate.exercises[i]
		assert.Equal(t, 0, ex.CircuitIndex, "%s should be in circuit 0", ex.Name)
		assert.Equal(t, wantSeconds, ex.WorkSeconds, "%s work seconds", ex.Name)
	}
}

// The positive control: an exercise outside any circuit must be passed through
// as NoCircuit. A CircuitIndex of 0 would file it under the first circuit.
func TestTemplatesCreate_LooseExerciseIsNotInACircuit(t *testing.T) {
	t.Cleanup(resetMocks)
	captureTemplateCreate()
	cookies := loginAs(t, "tmpl_create_loose", "lb")

	form := stretchCircuitForm()
	form.Set("exercise_count", "4")
	form.Set("exercise_name_3", "Bench Press")
	form.Set("circuit_index_3", "-1")

	w := postForm("/templates/new", form, cookies)
	assert.Equal(t, http.StatusFound, w.Code)

	require.Len(t, lastTemplateCreate.exercises, 4)
	bench := lastTemplateCreate.exercises[3]
	assert.Equal(t, "Bench Press", bench.Name)
	assert.Equal(t, models.NoCircuit, bench.CircuitIndex, "a loose exercise must not be filed under a circuit")
	assert.Equal(t, 0, bench.WorkSeconds)
}

// A template with no circuits at all must still submit exactly as it did before
// circuits existed — this is the regression guard for every template the
// engineer already has.
func TestTemplatesCreate_NoCircuitsAtAll(t *testing.T) {
	t.Cleanup(resetMocks)
	captureTemplateCreate()
	cookies := loginAs(t, "tmpl_create_nocircuit", "lb")

	w := postForm("/templates/new", url.Values{
		"name":            {"Upper Body A"},
		"exercise_count":  {"1"},
		"exercise_name_0": {"Bench Press"},
	}, cookies)

	assert.Equal(t, http.StatusFound, w.Code)
	assert.Empty(t, lastTemplateCreate.circuits)
	require.Len(t, lastTemplateCreate.exercises, 1)
	assert.Equal(t, models.NoCircuit, lastTemplateCreate.exercises[0].CircuitIndex)
}

// The form is never trusted: zero rounds and negative seconds are coerced to the
// bounds the table enforces rather than passed to the repository.
func TestTemplatesCreate_ClampsRoundsAndSeconds(t *testing.T) {
	t.Cleanup(resetMocks)
	captureTemplateCreate()
	cookies := loginAs(t, "tmpl_create_clamp", "lb")

	form := stretchCircuitForm()
	form.Set("circuit_rounds_0", "0")
	form.Set("circuit_transition_0", "-30")
	form.Set("work_seconds_0", "-10")

	w := postForm("/templates/new", form, cookies)
	assert.Equal(t, http.StatusFound, w.Code)

	require.Len(t, lastTemplateCreate.circuits, 1)
	assert.Equal(t, 1, lastTemplateCreate.circuits[0].Rounds, "zero rounds must clamp to one")
	assert.Equal(t, 0, lastTemplateCreate.circuits[0].TransitionSeconds, "a negative transition must clamp to zero")
	assert.Equal(t, 0, lastTemplateCreate.exercises[0].WorkSeconds, "negative work seconds must clamp to zero")
}

// If the user clears a circuit's name, that circuit is dropped. The exercises
// that pointed at it must become loose rather than shift onto whichever circuit
// slid into its index.
func TestTemplatesCreate_ExerciseOfDroppedCircuitBecomesLoose(t *testing.T) {
	t.Cleanup(resetMocks)
	captureTemplateCreate()
	cookies := loginAs(t, "tmpl_create_dropped", "lb")

	w := postForm("/templates/new", url.Values{
		"name":                 {"Mobility"},
		"circuit_count":        {"2"},
		"circuit_name_0":       {""}, // emptied out by the user
		"circuit_rounds_0":     {"3"},
		"circuit_name_1":       {"HIIT"},
		"circuit_rounds_1":     {"4"},
		"circuit_transition_1": {"15"},
		"exercise_count":       {"2"},
		"exercise_name_0":      {"Orphan"},
		"circuit_index_0":      {"0"}, // pointed at the circuit that is now gone
		"exercise_name_1":      {"Burpees"},
		"circuit_index_1":      {"1"}, // HIIT, which compacts from index 1 to index 0
	}, cookies)

	assert.Equal(t, http.StatusFound, w.Code)

	require.Len(t, lastTemplateCreate.circuits, 1)
	assert.Equal(t, "HIIT", lastTemplateCreate.circuits[0].Name)

	require.Len(t, lastTemplateCreate.exercises, 2)
	assert.Equal(t, models.NoCircuit, lastTemplateCreate.exercises[0].CircuitIndex,
		"an exercise whose circuit was dropped must become loose, not point at another circuit")
	assert.Equal(t, 0, lastTemplateCreate.exercises[1].CircuitIndex,
		"Burpees must follow HIIT to its compacted index")
}

// --- Circuits on the show page ---

// savedStretchCircuit sets up the mocks to return the Morning Stretch circuit as
// it comes back from the database.
func savedStretchCircuit() {
	id := circuitID
	setTemplateGetByIDExercises(testTemplateID, "Mobility", "",
		&models.TemplateExercise{ID: 1, Name: "shoulder stretch", IsTimeBased: true, Block: "stretch", SortOrder: 0, CircuitID: &id, WorkSeconds: 30},
		&models.TemplateExercise{ID: 2, Name: "hip flexor", IsTimeBased: true, Block: "stretch", SortOrder: 1, CircuitID: &id, WorkSeconds: 45},
		&models.TemplateExercise{ID: 3, Name: "hamstring", IsTimeBased: true, Block: "stretch", SortOrder: 2, CircuitID: &id, WorkSeconds: 30},
		&models.TemplateExercise{ID: 4, Name: "bench press", Block: "main", SortOrder: 3},
	)
	setTemplateGetCircuits(&models.TemplateCircuit{
		ID: circuitID, TemplateID: testTemplateID, Name: "Morning Stretch", Rounds: 1, TransitionSeconds: 5, SortOrder: 0,
	})
}

func TestTemplatesShow_RendersCircuitAsOneCard(t *testing.T) {
	t.Cleanup(resetMocks)
	savedStretchCircuit()
	cookies := loginAs(t, "tmpl_show_circuit", "lb")

	w := getPath(fmt.Sprintf("/templates/%d", testTemplateID), cookies)
	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()

	assert.Contains(t, body, "Morning Stretch")
	assert.Contains(t, body, "1 round &middot; 5s transition")
	assert.Contains(t, body, "shoulder stretch")
	assert.Contains(t, body, "30s")
	assert.Contains(t, body, "45s")
}

// The circuit's exercises are rendered inside the circuit card, so they must not
// also appear in the block list underneath as loose exercises.
func TestTemplatesShow_CircuitExercisesAreNotAlsoListedLoose(t *testing.T) {
	t.Cleanup(resetMocks)
	savedStretchCircuit()
	cookies := loginAs(t, "tmpl_show_nodupe", "lb")

	w := getPath(fmt.Sprintf("/templates/%d", testTemplateID), cookies)
	body := w.Body.String()

	assert.Equal(t, 1, strings.Count(body, "hip flexor"), "a circuit exercise must be listed exactly once")
	// The loose exercise still renders in its block, so this is not passing just
	// because the block list is empty.
	assert.Contains(t, body, "bench press")
	assert.Contains(t, body, "Stretch")
}

func TestTemplatesShow_PluralisesRounds(t *testing.T) {
	t.Cleanup(resetMocks)
	setTemplateGetByIDExercises(testTemplateID, "HIIT Day", "")
	setTemplateGetCircuits(&models.TemplateCircuit{
		ID: circuitID, Name: "HIIT", Rounds: 4, TransitionSeconds: 15, SortOrder: 0,
	})
	cookies := loginAs(t, "tmpl_show_plural", "lb")

	w := getPath(fmt.Sprintf("/templates/%d", testTemplateID), cookies)
	assert.Contains(t, w.Body.String(), "4 rounds &middot; 15s transition")
}

// --- Circuits on the edit form ---

func TestTemplatesEdit_PrefillsCircuit(t *testing.T) {
	t.Cleanup(resetMocks)
	savedStretchCircuit()
	cookies := loginAs(t, "tmpl_edit_circuit", "lb")

	w := getPath(fmt.Sprintf("/templates/%d/edit", testTemplateID), cookies)
	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()

	assert.Contains(t, body, `name="circuit_name_0" value="Morning Stretch"`)
	assert.Contains(t, body, `name="circuit_rounds_0" value="1"`)
	assert.Contains(t, body, `name="circuit_transition_0" value="5"`)
	assert.Contains(t, body, `id="circuit_count" value="1"`)
	// Each exercise keeps its own duration, and the loose one is counted too.
	assert.Contains(t, body, `name="work_seconds_0" value="30"`)
	assert.Contains(t, body, `name="work_seconds_1" value="45"`)
	assert.Contains(t, body, `id="exercise_count" value="4"`)
}

// Criterion 6: submitting with an empty name and a half-built circuit must come
// back with the circuit still filled in, not silently emptied. This is the path
// a missing c.Data key would break.
func TestTemplatesCreate_EmptyNameKeepsCircuitRows(t *testing.T) {
	t.Cleanup(resetMocks)
	cookies := loginAs(t, "tmpl_create_err_circuit", "lb")

	form := stretchCircuitForm()
	form.Set("name", "")

	w := postForm("/templates/new", form, cookies)
	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()

	assert.Contains(t, body, "required")
	assert.Contains(t, body, `name="circuit_name_0" value="Morning Stretch"`)
	assert.Contains(t, body, `name="circuit_rounds_0" value="1"`)
	assert.Contains(t, body, `name="circuit_transition_0" value="5"`)
	assert.Contains(t, body, `name="work_seconds_0" value="30"`)
	assert.Contains(t, body, `name="work_seconds_1" value="45"`)
	assert.Contains(t, body, `id="circuit_count" value="1"`)
}

func TestTemplatesUpdate_EmptyNameKeepsCircuitRows(t *testing.T) {
	t.Cleanup(resetMocks)
	setTemplateGetByID(testTemplateID, "Mobility", "", 1)
	cookies := loginAs(t, "tmpl_update_err_circuit", "lb")

	form := stretchCircuitForm()
	form.Set("name", "")

	w := postForm(fmt.Sprintf("/templates/%d", testTemplateID), form, cookies)
	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()

	assert.Contains(t, body, "required")
	assert.Contains(t, body, `name="circuit_name_0" value="Morning Stretch"`)
	assert.Contains(t, body, `name="circuit_transition_0" value="5"`)
	assert.Contains(t, body, `name="work_seconds_1" value="45"`)
}

// The update path is a second, separate parse of the same form.
func TestTemplatesUpdate_WithCircuit(t *testing.T) {
	t.Cleanup(resetMocks)
	setTemplateGetByID(testTemplateID, "Mobility", "", 1)
	captureTemplateUpdate()
	cookies := loginAs(t, "tmpl_update_circuit", "lb")

	w := postForm(fmt.Sprintf("/templates/%d", testTemplateID), stretchCircuitForm(), cookies)
	assert.Equal(t, http.StatusFound, w.Code)

	require.Len(t, lastTemplateUpdate.circuits, 1)
	assert.Equal(t, "Morning Stretch", lastTemplateUpdate.circuits[0].Name)
	assert.Equal(t, 5, lastTemplateUpdate.circuits[0].TransitionSeconds)

	require.Len(t, lastTemplateUpdate.exercises, 3)
	assert.Equal(t, 45, lastTemplateUpdate.exercises[1].WorkSeconds)
	assert.Equal(t, 0, lastTemplateUpdate.exercises[1].CircuitIndex)
}
