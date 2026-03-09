package controllers_test

import (
	"net/http"
	"net/url"
	"testing"

	"myGymPal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Programs list ---

func TestProgramsIndex_Unauthenticated(t *testing.T) {
	w := getPath("/programs", nil)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/login", w.Header().Get("Location"))
}

func TestProgramsIndex_Empty(t *testing.T) {
	cookies := loginAs(t, "prog_idx_empty", "password123")
	t.Cleanup(func() { models.DeleteUserByUsername("prog_idx_empty") })

	w := getPath("/programs", cookies)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "No programs yet")
}

func TestProgramsIndex_ShowsPrograms(t *testing.T) {
	cookies := loginAs(t, "prog_idx_list", "password123")
	t.Cleanup(func() { models.DeleteUserByUsername("prog_idx_list") })

	// Create a program directly via model so we can verify it appears.
	user, err := models.GetUserByUsername("prog_idx_list")
	require.NoError(t, err)
	p, err := models.CreateProgram(user.ID, "My Test Program", testProgramDate, 4, 8)
	require.NoError(t, err)
	t.Cleanup(func() { models.DeleteProgram(p.ID, user.ID) })

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
	cookies := loginAs(t, "prog_new_form", "password123")
	t.Cleanup(func() { models.DeleteUserByUsername("prog_new_form") })

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
	}, nil)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/login", w.Header().Get("Location"))
}

func TestProgramsCreate_Success(t *testing.T) {
	cookies := loginAs(t, "prog_create_ok", "password123")
	t.Cleanup(func() { models.DeleteUserByUsername("prog_create_ok") })

	w := postForm("/programs", url.Values{
		"name":            {"Hypertrophy Block"},
		"start_date":      {"2025-01-06"},
		"num_phases":      {"4"},
		"weeks_per_phase": {"6"},
	}, cookies)

	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/programs", w.Header().Get("Location"))

	// Follow the redirect — success message should appear on the programs page.
	allCookies := append(cookies, w.Result().Cookies()...)
	w2 := getPath("/programs", allCookies)
	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Contains(t, w2.Body.String(), "Program created.")

	user, err := models.GetUserByUsername("prog_create_ok")
	require.NoError(t, err)
	programs, err := models.GetProgramsByUserID(user.ID)
	require.NoError(t, err)
	require.Len(t, programs, 1)
	assert.Equal(t, "Hypertrophy Block", programs[0].Name)
	assert.Equal(t, 4, programs[0].NumPhases)
	assert.Equal(t, 6, programs[0].WeeksPerPhase)
}

func TestProgramsCreate_EmptyName(t *testing.T) {
	cookies := loginAs(t, "prog_create_noname", "password123")
	t.Cleanup(func() { models.DeleteUserByUsername("prog_create_noname") })

	w := postForm("/programs", url.Values{
		"name":            {""},
		"start_date":      {"2025-01-06"},
		"num_phases":      {"4"},
		"weeks_per_phase": {"8"},
	}, cookies)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "required")
}

func TestProgramsCreate_InvalidDate(t *testing.T) {
	cookies := loginAs(t, "prog_create_baddate", "password123")
	t.Cleanup(func() { models.DeleteUserByUsername("prog_create_baddate") })

	w := postForm("/programs", url.Values{
		"name":            {"My Program"},
		"start_date":      {"not-a-date"},
		"num_phases":      {"4"},
		"weeks_per_phase": {"8"},
	}, cookies)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "valid date")
}

func TestProgramsCreate_InvalidPhases(t *testing.T) {
	cookies := loginAs(t, "prog_create_badphases", "password123")
	t.Cleanup(func() { models.DeleteUserByUsername("prog_create_badphases") })

	w := postForm("/programs", url.Values{
		"name":            {"My Program"},
		"start_date":      {"2025-01-06"},
		"num_phases":      {"0"},
		"weeks_per_phase": {"8"},
	}, cookies)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "positive")
}

func TestProgramsCreate_InvalidWeeksPerPhase(t *testing.T) {
	cookies := loginAs(t, "prog_create_badweeks", "password123")
	t.Cleanup(func() { models.DeleteUserByUsername("prog_create_badweeks") })

	w := postForm("/programs", url.Values{
		"name":            {"My Program"},
		"start_date":      {"2025-01-06"},
		"num_phases":      {"4"},
		"weeks_per_phase": {"0"},
	}, cookies)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "positive")
}

func TestProgramsCreate_ReentersFormValues(t *testing.T) {
	cookies := loginAs(t, "prog_create_reenter", "password123")
	t.Cleanup(func() { models.DeleteUserByUsername("prog_create_reenter") })

	w := postForm("/programs", url.Values{
		"name":            {"My Sticky Program"},
		"start_date":      {"2025-03-01"},
		"num_phases":      {"0"},
		"weeks_per_phase": {"10"},
	}, cookies)

	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, "My Sticky Program")
	assert.Contains(t, body, "2025-03-01")
	assert.Contains(t, body, "10")
}
