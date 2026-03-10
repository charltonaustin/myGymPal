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
	assert.Contains(t, body, `name="name"`)
	assert.Contains(t, body, `name="focus"`)
	assert.Contains(t, body, `name="exercise_name_0"`)
	assert.Contains(t, body, `name="rep_min_0"`)
	assert.Contains(t, body, `name="rep_max_0"`)
	assert.Contains(t, body, "lb")
}

func TestTemplatesNew_ShowsKgUnit(t *testing.T) {
	t.Cleanup(resetMocks)
	cookies := loginAs(t, "tmpl_new_form_kg", "kg")

	w := getPath("/templates/new", cookies)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "kg")
}

// --- Create template ---

func TestTemplatesCreate_Unauthenticated(t *testing.T) {
	w := postForm("/templates/new", url.Values{
		"name":            {"Test"},
		"exercise_count":  {"1"},
		"exercise_name_0": {"Bench Press"},
		"rep_min_0":       {"8"},
		"rep_max_0":       {"10"},
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
		"goal_weight_0":   {"60"},
		"rep_min_0":       {"8"},
		"rep_max_0":       {"10"},
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
		"rep_min_0":       {"8"},
		"rep_max_0":       {"12"},
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
		"rep_min_0":       {"8"},
		"rep_max_0":       {"10"},
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

func TestTemplatesCreate_InvalidRepRange(t *testing.T) {
	t.Cleanup(resetMocks)
	setTemplateCreateError(errors.New("rep_max must be at least rep_min"))
	cookies := loginAs(t, "tmpl_create_badrange", "lb")

	w := postForm("/templates/new", url.Values{
		"name":            {"My Template"},
		"exercise_count":  {"1"},
		"exercise_name_0": {"Bench Press"},
		"rep_min_0":       {"12"},
		"rep_max_0":       {"8"},
	}, cookies)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "rep_max")
}

func TestTemplatesCreate_ReentersFormValues(t *testing.T) {
	t.Cleanup(resetMocks)
	setTemplateCreateError(errors.New("rep_max must be at least rep_min"))
	cookies := loginAs(t, "tmpl_create_reenter", "lb")

	w := postForm("/templates/new", url.Values{
		"name":            {"Sticky Template"},
		"focus":           {"Legs"},
		"exercise_count":  {"1"},
		"exercise_name_0": {"Squat"},
		"rep_min_0":       {"5"},
		"rep_max_0":       {"3"},
	}, cookies)

	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, "Sticky Template")
	assert.Contains(t, body, "Squat")
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
