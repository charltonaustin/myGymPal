package controllers_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

// --- Programs list ---

func TestProgramsIndex_Unauthenticated(t *testing.T) {
	w := getPath("/programs", nil)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/login", w.Header().Get("Location"))
}

func TestProgramsIndex_Empty(t *testing.T) {
	t.Cleanup(resetMocks)
	setProgramsGetAllEmpty()
	cookies := loginAs(t, "prog_idx_empty", "lb")

	w := getPath("/programs", cookies)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "No programs yet")
}

func TestProgramsIndex_ShowsPrograms(t *testing.T) {
	t.Cleanup(resetMocks)
	setProgramsGetAllWithOne(1, "My Test Program", 4, 8)
	cookies := loginAs(t, "prog_idx_list", "lb")

	w := getPath("/programs", cookies)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "My Test Program")
	assert.Contains(t, w.Body.String(), "4 phases")
}

// --- New program form ---

func TestProgramsNew_Unauthenticated(t *testing.T) {
	w := getPath("/programs/new", nil)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/login", w.Header().Get("Location"))
}

func TestProgramsNew_ShowsForm(t *testing.T) {
	t.Cleanup(resetMocks)
	cookies := loginAs(t, "prog_new_form", "lb")

	w := getPath("/programs/new", cookies)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "New Training Program")
	assert.Contains(t, w.Body.String(), `name="name"`)
	assert.Contains(t, w.Body.String(), `name="start_date"`)
	assert.Contains(t, w.Body.String(), `name="num_phases"`)
	assert.Contains(t, w.Body.String(), `name="weeks_per_phase"`)
}

// --- Create program ---

func TestProgramsCreate_Unauthenticated(t *testing.T) {
	w := postForm("/programs", url.Values{
		"name":            {"Test"},
		"start_date":      {"2025-01-06"},
		"num_phases":      {"4"},
		"weeks_per_phase": {"8"},
		"default_rep_min": {"10"},
		"default_rep_max": {"12"},
	}, nil)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/login", w.Header().Get("Location"))
}

func TestProgramsCreate_Success(t *testing.T) {
	t.Cleanup(resetMocks)
	captureProgramCreate()
	cookies := loginAs(t, "prog_create_ok", "lb")

	w := postForm("/programs", url.Values{
		"name":            {"Hypertrophy Block"},
		"start_date":      {"2025-01-06"},
		"num_phases":      {"4"},
		"weeks_per_phase": {"6"},
		"default_rep_min": {"10"},
		"default_rep_max": {"12"},
	}, cookies)

	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/programs", w.Header().Get("Location"))
	assert.Equal(t, "Hypertrophy Block", lastProgramCreate.name)
	assert.Equal(t, 4, lastProgramCreate.numPhases)
	assert.Equal(t, 6, lastProgramCreate.weeksPerPhase)

	// Follow the redirect — success flash should appear.
	allCookies := append(cookies, w.Result().Cookies()...)
	w2 := getPath("/programs", allCookies)
	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Contains(t, w2.Body.String(), "Program created.")
}

func TestProgramsCreate_EmptyName(t *testing.T) {
	t.Cleanup(resetMocks)
	cookies := loginAs(t, "prog_create_noname", "lb")

	w := postForm("/programs", url.Values{
		"name":            {""},
		"start_date":      {"2025-01-06"},
		"num_phases":      {"4"},
		"weeks_per_phase": {"8"},
		"default_rep_min": {"10"},
		"default_rep_max": {"12"},
	}, cookies)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "required")
}

func TestProgramsCreate_InvalidDate(t *testing.T) {
	t.Cleanup(resetMocks)
	cookies := loginAs(t, "prog_create_baddate", "lb")

	w := postForm("/programs", url.Values{
		"name":            {"My Program"},
		"start_date":      {"not-a-date"},
		"num_phases":      {"4"},
		"weeks_per_phase": {"8"},
		"default_rep_min": {"10"},
		"default_rep_max": {"12"},
	}, cookies)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "valid date")
}

func TestProgramsCreate_InvalidPhases(t *testing.T) {
	t.Cleanup(resetMocks)
	cookies := loginAs(t, "prog_create_badphases", "lb")

	w := postForm("/programs", url.Values{
		"name":            {"My Program"},
		"start_date":      {"2025-01-06"},
		"num_phases":      {"0"},
		"weeks_per_phase": {"8"},
		"default_rep_min": {"10"},
		"default_rep_max": {"12"},
	}, cookies)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "positive")
}

func TestProgramsCreate_InvalidDefaultRepMin(t *testing.T) {
	t.Cleanup(resetMocks)
	cookies := loginAs(t, "prog_create_bad_repmin", "lb")

	w := postForm("/programs", url.Values{
		"name":            {"My Program"},
		"start_date":      {"2025-01-06"},
		"num_phases":      {"4"},
		"weeks_per_phase": {"8"},
		"default_rep_min": {"0"},
		"default_rep_max": {"12"},
	}, cookies)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "positive")
}

func TestProgramsCreate_DefaultRepMaxLessThanMin(t *testing.T) {
	t.Cleanup(resetMocks)
	cookies := loginAs(t, "prog_create_bad_repmax", "lb")

	w := postForm("/programs", url.Values{
		"name":            {"My Program"},
		"start_date":      {"2025-01-06"},
		"num_phases":      {"4"},
		"weeks_per_phase": {"8"},
		"default_rep_min": {"12"},
		"default_rep_max": {"8"},
	}, cookies)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "least")
}

func TestProgramsCreate_InvalidWeeksPerPhase(t *testing.T) {
	t.Cleanup(resetMocks)
	cookies := loginAs(t, "prog_create_badweeks", "lb")

	w := postForm("/programs", url.Values{
		"name":            {"My Program"},
		"start_date":      {"2025-01-06"},
		"num_phases":      {"4"},
		"weeks_per_phase": {"0"},
		"default_rep_min": {"10"},
		"default_rep_max": {"12"},
	}, cookies)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "positive")
}

// --- Program detail / show ---

func TestProgramShow_Unauthenticated(t *testing.T) {
	w := getPath("/programs/1", nil)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/login", w.Header().Get("Location"))
}

func TestProgramShow_WrongUser(t *testing.T) {
	t.Cleanup(resetMocks)
	setProgramGetByIDError(errors.New("not found"))
	cookies := loginAs(t, "prog_show_other", "lb")

	w := getPath("/programs/1", cookies)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/programs", w.Header().Get("Location"))
}

func TestProgramShow_ShowsPhases(t *testing.T) {
	t.Cleanup(resetMocks)
	setProgramGetByID("Show Program", 3)
	setPhasesGetByProgram(3)
	cookies := loginAs(t, "prog_show_phases", "lb")

	w := getPath("/programs/1", cookies)
	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, "Show Program")
	assert.Contains(t, body, "Phase 1")
	assert.Contains(t, body, "Phase 2")
	assert.Contains(t, body, "Phase 3")
	assert.Contains(t, body, `name="rep_min_1"`)
	assert.Contains(t, body, `name="rep_max_1"`)
}

// --- Update phase rep ranges ---

func TestUpdatePhases_Unauthenticated(t *testing.T) {
	w := postForm("/programs/1", url.Values{"rep_min_1": {"10"}, "rep_max_1": {"12"}}, nil)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/login", w.Header().Get("Location"))
}

func TestUpdatePhases_Success(t *testing.T) {
	t.Cleanup(resetMocks)
	setProgramGetByID("Phase Update Program", 2)
	setPhasesGetByProgram(2)
	capturePhasesUpdateRepRanges()
	cookies := loginAs(t, "prog_update_phases_ok", "lb")

	w := postForm("/programs/1", url.Values{
		"rep_min_1": {"10"}, "rep_max_1": {"12"},
		"rep_min_2": {"8"}, "rep_max_2": {"10"},
	}, cookies)

	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, fmt.Sprintf("/programs/%d", testProgramID), w.Header().Get("Location"))
	repMin0, repMax0 := getLastPhaseUpdate(0)
	repMin1, repMax1 := getLastPhaseUpdate(1)
	assert.Equal(t, 10, repMin0)
	assert.Equal(t, 12, repMax0)
	assert.Equal(t, 8, repMin1)
	assert.Equal(t, 10, repMax1)
}

func TestUpdatePhases_InvalidRange(t *testing.T) {
	t.Cleanup(resetMocks)
	setProgramGetByID("Bad Range Program", 1)
	setPhasesGetByProgram(1)
	setPhasesUpdateRepRangesError(errors.New("rep_max must be at least rep_min"))
	cookies := loginAs(t, "prog_update_phases_bad", "lb")

	// max < min
	w := postForm("/programs/1", url.Values{
		"rep_min_1": {"12"}, "rep_max_1": {"8"},
	}, cookies)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "rep_max")
}

// --- Delete program ---

func TestProgramDelete_Unauthenticated(t *testing.T) {
	w := postForm("/programs/1/delete", url.Values{}, nil)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/login", w.Header().Get("Location"))
}

func TestProgramDelete_NotFound(t *testing.T) {
	t.Cleanup(resetMocks)
	mockPrograms.DeleteFn = func(id, userID int64) error {
		return errors.New("not found")
	}
	cookies := loginAs(t, "prog_delete_notfound", "lb")

	w := postForm("/programs/1/delete", url.Values{}, cookies)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/programs", w.Header().Get("Location"))
}

func TestProgramDelete_Success(t *testing.T) {
	t.Cleanup(resetMocks)
	cookies := loginAs(t, "prog_delete_ok", "lb")

	w := postForm(fmt.Sprintf("/programs/%d/delete", testProgramID), url.Values{}, cookies)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/programs", w.Header().Get("Location"))
}

func TestProgramsCreate_ReentersFormValues(t *testing.T) {
	t.Cleanup(resetMocks)
	cookies := loginAs(t, "prog_create_reenter", "lb")

	w := postForm("/programs", url.Values{
		"name":            {"My Sticky Program"},
		"start_date":      {"2025-03-01"},
		"num_phases":      {"0"},
		"weeks_per_phase": {"10"},
		"default_rep_min": {"8"},
		"default_rep_max": {"12"},
	}, cookies)

	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, "My Sticky Program")
	assert.Contains(t, body, "2025-03-01")
	assert.Contains(t, body, "10")
}
