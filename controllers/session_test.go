package controllers_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

	w := postForm("/programs/1/sessions", url.Values{}, cookies)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, fmt.Sprintf("/sessions/%d", testSessionID), w.Header().Get("Location"))
	assert.Equal(t, 1, lastSessionCreate.workoutNumber)
}

func TestSessionCreate_IncrementsWorkoutNumber(t *testing.T) {
	t.Cleanup(resetMocks)
	setProgramGetByIDWithDates("My Program", 4, 8, testProgramDate)
	setSessionCountByProgram(4)
	captureSessionCreate()
	cookies := loginAs(t, "session_create_workout5", "lb")

	w := postForm("/programs/1/sessions", url.Values{}, cookies)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, 5, lastSessionCreate.workoutNumber)
}

func TestSessionCreate_DeloadWeekFlagged(t *testing.T) {
	t.Cleanup(resetMocks)
	// Start date 49 days (7 weeks) before today → currently in week 8 of phase 1 → deload.
	startDate := time.Now().UTC().AddDate(0, 0, -49)
	setProgramGetByIDWithDates("My Program", 4, 8, startDate)
	captureSessionCreate()
	cookies := loginAs(t, "session_create_deload", "lb")

	w := postForm("/programs/1/sessions", url.Values{}, cookies)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.True(t, lastSessionCreate.isDeload)
	assert.Equal(t, 8, lastSessionCreate.weekNumber)
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
