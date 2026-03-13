package controllers_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"myGymPal/models"
)

// --- Create session ---

func TestSessionCreate_Unauthenticated(t *testing.T) {
	w := postForm("/programs/1/sessions", url.Values{}, nil)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/login", w.Header().Get("Location"))
}

func TestSessionCreate_ProgramNotFound(t *testing.T) {
	t.Cleanup(resetMocks)
	setProgramGetByIDError(errors.New("not found"))
	cookies := loginAs(t, "session_create_noprog", "lb")

	w := postForm("/programs/1/sessions", url.Values{}, cookies)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/programs", w.Header().Get("Location"))
}

func TestSessionCreate_Success(t *testing.T) {
	t.Cleanup(resetMocks)
	setProgramGetByIDWithDates("My Program", 4, 8, testProgramDate)
	captureSessionCreate()
	cookies := loginAs(t, "session_create_ok", "lb")

	w := postForm("/programs/1/sessions", url.Values{
		"phase_number":   {"1"},
		"week_number":    {"1"},
		"workout_number": {"1"},
	}, cookies)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, fmt.Sprintf("/sessions/%d", testSessionID), w.Header().Get("Location"))
	assert.Equal(t, 1, lastSessionCreate.workoutNumber)
}

func TestSessionCreate_IncrementsWorkoutNumber(t *testing.T) {
	t.Cleanup(resetMocks)
	setProgramGetByIDWithDates("My Program", 4, 8, testProgramDate)
	captureSessionCreate()
	cookies := loginAs(t, "session_create_workout5", "lb")

	w := postForm("/programs/1/sessions", url.Values{
		"phase_number":   {"1"},
		"week_number":    {"1"},
		"workout_number": {"5"},
	}, cookies)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, 5, lastSessionCreate.workoutNumber)
}

func TestSessionCreate_DeloadWeekFlagged(t *testing.T) {
	t.Cleanup(resetMocks)
	setProgramGetByIDWithDates("My Program", 4, 8, testProgramDate)
	captureSessionCreate()
	cookies := loginAs(t, "session_create_deload", "lb")

	// Submitting week 8 on an 8-week program → deload.
	w := postForm("/programs/1/sessions", url.Values{
		"phase_number":   {"1"},
		"week_number":    {"8"},
		"workout_number": {"1"},
	}, cookies)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.True(t, lastSessionCreate.isDeload)
	assert.Equal(t, 8, lastSessionCreate.weekNumber)
}

func TestSessionCreate_NonDeloadWeek(t *testing.T) {
	t.Cleanup(resetMocks)
	setProgramGetByIDWithDates("My Program", 4, 8, testProgramDate)
	captureSessionCreate()
	cookies := loginAs(t, "session_create_nodeload", "lb")

	w := postForm("/programs/1/sessions", url.Values{
		"phase_number":   {"1"},
		"week_number":    {"4"},
		"workout_number": {"1"},
	}, cookies)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.False(t, lastSessionCreate.isDeload)
}

func TestSessionCreate_InvalidInputRedirectsToNew(t *testing.T) {
	t.Cleanup(resetMocks)
	setProgramGetByIDWithDates("My Program", 4, 8, testProgramDate)
	cookies := loginAs(t, "session_create_invalid", "lb")

	w := postForm("/programs/1/sessions", url.Values{
		"phase_number":   {"0"},
		"week_number":    {"1"},
		"workout_number": {"1"},
	}, cookies)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, fmt.Sprintf("/programs/%d/sessions/new", testProgramID), w.Header().Get("Location"))
}

func TestSessionCreate_NoTemplateSkipsCopy(t *testing.T) {
	t.Cleanup(resetMocks)
	setProgramGetByIDWithDates("My Program", 4, 8, testProgramDate)
	captureSessionCreate()
	captureSessionExerciseCreates()
	cookies := loginAs(t, "session_create_notmpl", "lb")

	w := postForm("/programs/1/sessions", url.Values{
		"phase_number":   {"1"},
		"week_number":    {"1"},
		"workout_number": {"1"},
	}, cookies)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Len(t, sessionExerciseCreateNames, 0)
}

func TestSessionCreate_CopiesTemplateExercises(t *testing.T) {
	t.Cleanup(resetMocks)
	setProgramGetByIDWithDates("My Program", 4, 8, testProgramDate)
	captureSessionCreate()
	setTemplateGetByID(testTemplateID, "Upper Body A", "Upper", 3)
	captureSessionExerciseCreates()
	cookies := loginAs(t, "session_create_tmpl", "lb")

	w := postForm("/programs/1/sessions", url.Values{
		"phase_number":   {"1"},
		"week_number":    {"1"},
		"workout_number": {"1"},
		"template_id":    {fmt.Sprintf("%d", testTemplateID)},
	}, cookies)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Len(t, sessionExerciseCreateNames, 3)
	assert.Equal(t, "Exercise 1", sessionExerciseCreateNames[0])
	assert.Equal(t, "Exercise 3", sessionExerciseCreateNames[2])
}

func TestSessionCreate_SetsGoalRepsFromPhase(t *testing.T) {
	t.Cleanup(resetMocks)
	setProgramGetByIDWithDates("My Program", 4, 8, testProgramDate)
	captureSessionCreate()
	setTemplateGetByID(testTemplateID, "Upper Body A", "Upper", 2)
	setPhasesGetByProgram(4) // phases 1–4 all have RepMin=10
	var capturedGoalReps []int
	mockSessionExercises.CreateFn = func(sessionID int64, name string, isBodyweight bool, goalWeight float64, weightUnit string, goalReps int, block string, isTimeBased bool, goalSeconds int) (*models.SessionExercise, error) {
		capturedGoalReps = append(capturedGoalReps, goalReps)
		return &models.SessionExercise{ID: testExerciseID, SessionID: sessionID, Name: name, GoalReps: goalReps}, nil
	}
	cookies := loginAs(t, "session_create_goalreps", "lb")

	w := postForm("/programs/1/sessions", url.Values{
		"phase_number":   {"1"},
		"week_number":    {"1"},
		"workout_number": {"1"},
		"template_id":    {fmt.Sprintf("%d", testTemplateID)},
	}, cookies)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Len(t, capturedGoalReps, 2)
	// setPhasesGetByProgram gives RepMin=10 for all phases
	assert.Equal(t, 10, capturedGoalReps[0])
	assert.Equal(t, 10, capturedGoalReps[1])
}

// --- New session form ---

func TestSessionNew_Unauthenticated(t *testing.T) {
	w := getPath("/programs/1/sessions/new", nil)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/login", w.Header().Get("Location"))
}

func TestSessionNew_ProgramNotFound(t *testing.T) {
	t.Cleanup(resetMocks)
	setProgramGetByIDError(errors.New("not found"))
	cookies := loginAs(t, "session_new_noprog", "lb")

	w := getPath("/programs/1/sessions/new", cookies)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/programs", w.Header().Get("Location"))
}

func TestSessionNew_ShowsForm(t *testing.T) {
	t.Cleanup(resetMocks)
	setProgramGetByIDWithDates("My Program", 4, 8, testProgramDate)
	setSessionLatestByProgram(1, 1, 3)
	cookies := loginAs(t, "session_new_ok", "lb")

	w := getPath("/programs/1/sessions/new", cookies)
	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, "My Program")
	assert.Contains(t, body, `name="phase_number"`)
	assert.Contains(t, body, `name="week_number"`)
	assert.Contains(t, body, `name="workout_number"`)
	// Last workout was #3, so next is #4
	assert.Contains(t, body, `value="4"`)
}

func TestSessionNew_ShowsTemplates(t *testing.T) {
	t.Cleanup(resetMocks)
	setProgramGetByIDWithDates("My Program", 4, 8, testProgramDate)
	setTemplatesGetAllWithOne(testTemplateID, "Upper Body A", "Upper")
	cookies := loginAs(t, "session_new_templates", "lb")

	w := getPath("/programs/1/sessions/new", cookies)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Upper Body A")
}

// --- Show session ---

func TestSessionShow_Unauthenticated(t *testing.T) {
	w := getPath("/sessions/1", nil)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/login", w.Header().Get("Location"))
}

func TestSessionShow_NotFound(t *testing.T) {
	t.Cleanup(resetMocks)
	setSessionGetByIDError(errors.New("not found"))
	cookies := loginAs(t, "session_show_notfound", "lb")

	w := getPath("/sessions/1", cookies)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/programs", w.Header().Get("Location"))
}

func TestSessionShow_ShowsSession(t *testing.T) {
	t.Cleanup(resetMocks)
	setSessionGetByID(2, 3, 5, false)
	setProgramGetByID("My Program", 4)
	cookies := loginAs(t, "session_show_ok", "lb")

	w := getPath("/sessions/99", cookies)
	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, "My Program")
	assert.Contains(t, body, "Phase 2")
	assert.Contains(t, body, "Week 3")
	assert.Contains(t, body, "Session #5")
}

func TestSessionShow_DeloadBadge(t *testing.T) {
	t.Cleanup(resetMocks)
	setSessionGetByID(1, 8, 8, true)
	setProgramGetByID("My Program", 4)
	cookies := loginAs(t, "session_show_deload", "lb")

	w := getPath("/sessions/99", cookies)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Deload")
}

func TestSessionShow_ShowsExercises(t *testing.T) {
	t.Cleanup(resetMocks)
	setSessionGetByID(1, 1, 1, false)
	setProgramGetByID("My Program", 4)
	setSessionExerciseGetBySessionWithOne("Bench Press", "lb")
	cookies := loginAs(t, "session_show_exercises", "lb")

	w := getPath("/sessions/99", cookies)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Bench Press")
}

func TestSessionShow_ShowsPhaseRepRange(t *testing.T) {
	t.Cleanup(resetMocks)
	setSessionGetByID(2, 1, 1, false) // phase 2
	setProgramGetByID("My Program", 4)
	setPhasesGetByProgram(4) // all phases RepMin=10, RepMax=12
	setSessionExerciseGetBySessionWithOne("Bench Press", "lb")
	cookies := loginAs(t, "session_show_reprange", "lb")

	w := getPath("/sessions/99", cookies)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "10–12 reps")
}

func TestSessionShow_DefaultsWeightAndReps(t *testing.T) {
	t.Cleanup(resetMocks)
	setSessionGetByID(1, 1, 1, false)
	setProgramGetByID("My Program", 4)
	mockSessionExercises.GetBySessionFn = func(sessionID int64) ([]*models.SessionExerciseView, error) {
		ex := &models.SessionExercise{
			ID: testExerciseID, SessionID: sessionID,
			Name: "Squat", GoalWeight: 135, WeightUnit: "lb", GoalReps: 10,
		}
		return []*models.SessionExerciseView{{Exercise: ex}}, nil
	}
	cookies := loginAs(t, "session_show_defaults", "lb")

	w := getPath("/sessions/99", cookies)
	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, `value="135"`)
	assert.Contains(t, body, `value="10"`)
}

// --- Add exercise ---

func TestSessionAddExercise_Unauthenticated(t *testing.T) {
	w := postForm("/sessions/99/exercises", url.Values{"name": {"Squat"}}, nil)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/login", w.Header().Get("Location"))
}

func TestSessionAddExercise_SessionNotFound(t *testing.T) {
	t.Cleanup(resetMocks)
	setSessionGetByIDError(errors.New("not found"))
	cookies := loginAs(t, "add_exercise_nosession", "lb")

	w := postForm("/sessions/99/exercises", url.Values{"name": {"Squat"}}, cookies)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/programs", w.Header().Get("Location"))
}

func TestSessionAddExercise_Success(t *testing.T) {
	t.Cleanup(resetMocks)
	setSessionGetByID(1, 1, 1, false)
	cookies := loginAs(t, "add_exercise_ok", "lb")

	w := postForm("/sessions/99/exercises", url.Values{
		"name":        {"Bench Press"},
		"goal_weight": {"135"},
		"weight_unit": {"lb"},
	}, cookies)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, fmt.Sprintf("/sessions/%d", testSessionID), w.Header().Get("Location"))
}

func TestSessionAddExercise_EmptyNameRedirects(t *testing.T) {
	t.Cleanup(resetMocks)
	setSessionGetByID(1, 1, 1, false)
	cookies := loginAs(t, "add_exercise_empty", "lb")

	w := postForm("/sessions/99/exercises", url.Values{"name": {""}}, cookies)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, fmt.Sprintf("/sessions/%d", testSessionID), w.Header().Get("Location"))
}

// --- Log set ---

func TestSessionLogSet_Unauthenticated(t *testing.T) {
	w := postForm("/sessions/99/exercises/77/sets", url.Values{"actual_reps": {"8"}}, nil)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/login", w.Header().Get("Location"))
}

func TestSessionLogSet_SessionNotFound(t *testing.T) {
	t.Cleanup(resetMocks)
	setSessionGetByIDError(errors.New("not found"))
	cookies := loginAs(t, "log_set_nosession", "lb")

	w := postForm("/sessions/99/exercises/77/sets", url.Values{"actual_reps": {"8"}}, cookies)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/programs", w.Header().Get("Location"))
}

func TestSessionLogSet_ExerciseNotFound(t *testing.T) {
	t.Cleanup(resetMocks)
	setSessionGetByID(1, 1, 1, false)
	setSessionExerciseGetByIDError(errors.New("not found"))
	cookies := loginAs(t, "log_set_noexercise", "lb")

	w := postForm("/sessions/99/exercises/77/sets", url.Values{"actual_reps": {"8"}}, cookies)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, fmt.Sprintf("/sessions/%d", testSessionID), w.Header().Get("Location"))
}

func TestSessionLogSet_Success(t *testing.T) {
	t.Cleanup(resetMocks)
	setSessionGetByID(1, 1, 1, false)
	setSessionExerciseGetByID(testSessionID)
	captureLogSet()
	cookies := loginAs(t, "log_set_ok", "lb")

	w := postForm("/sessions/99/exercises/77/sets", url.Values{
		"actual_weight": {"135"},
		"weight_unit":   {"lb"},
		"actual_reps":   {"10"},
	}, cookies)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, fmt.Sprintf("/sessions/%d", testSessionID), w.Header().Get("Location"))
	assert.Equal(t, float64(135), lastLogSet.actualWeight)
	assert.Equal(t, 10, lastLogSet.actualReps)
	assert.Equal(t, "lb", lastLogSet.weightUnit)
	assert.Equal(t, 1, lastLogSet.setNumber)
}

func TestSessionLogSet_IncrementsSetNumber(t *testing.T) {
	t.Cleanup(resetMocks)
	setSessionGetByID(1, 1, 1, false)
	setSessionExerciseGetByID(testSessionID)
	setLogSetCountByExercise(2)
	captureLogSet()
	cookies := loginAs(t, "log_set_increment", "lb")

	w := postForm("/sessions/99/exercises/77/sets", url.Values{
		"actual_weight": {"135"},
		"weight_unit":   {"lb"},
		"actual_reps":   {"8"},
	}, cookies)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, 3, lastLogSet.setNumber)
}

// --- Delete session ---

func TestSessionDelete_Unauthenticated(t *testing.T) {
	w := postForm("/sessions/99/delete", url.Values{}, nil)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/login", w.Header().Get("Location"))
}

func TestSessionDelete_NotFound(t *testing.T) {
	t.Cleanup(resetMocks)
	setSessionGetByIDError(errors.New("not found"))
	cookies := loginAs(t, "session_delete_notfound", "lb")

	w := postForm("/sessions/99/delete", url.Values{}, cookies)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/programs", w.Header().Get("Location"))
}

func TestSessionDelete_Success(t *testing.T) {
	t.Cleanup(resetMocks)
	setSessionGetByID(1, 1, 1, false)
	cookies := loginAs(t, "session_delete_ok", "lb")

	w := postForm("/sessions/99/delete", url.Values{}, cookies)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, fmt.Sprintf("/programs/%d", testProgramID), w.Header().Get("Location"))
}

func TestSessionDelete_DeleteError(t *testing.T) {
	t.Cleanup(resetMocks)
	setSessionGetByID(1, 1, 1, false)
	mockSessions.DeleteFn = func(id, userID int64) error {
		return errors.New("db error")
	}
	cookies := loginAs(t, "session_delete_err", "lb")

	w := postForm("/sessions/99/delete", url.Values{}, cookies)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, fmt.Sprintf("/programs/%d", testProgramID), w.Header().Get("Location"))
}
