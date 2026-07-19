package controllers_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
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
