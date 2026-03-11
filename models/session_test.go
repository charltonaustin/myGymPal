package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCalculateNextSession_FirstSession(t *testing.T) {
	phase, week, workoutNum, isDeload := CalculateNextSession(0, 8, 4)
	assert.Equal(t, 1, phase)
	assert.Equal(t, 1, week)
	assert.Equal(t, 1, workoutNum)
	assert.False(t, isDeload)
}

func TestCalculateNextSession_SecondWorkout(t *testing.T) {
	phase, week, workoutNum, isDeload := CalculateNextSession(1, 8, 4)
	assert.Equal(t, 1, phase)
	assert.Equal(t, 1, week)
	assert.Equal(t, 2, workoutNum)
	assert.False(t, isDeload)
}

func TestCalculateNextSession_WeekAdvances(t *testing.T) {
	// After 4 sessions (workoutsPerWeek=4), next should be week 2 workout 1.
	phase, week, workoutNum, isDeload := CalculateNextSession(4, 8, 4)
	assert.Equal(t, 1, phase)
	assert.Equal(t, 2, week)
	assert.Equal(t, 1, workoutNum)
	assert.False(t, isDeload)
}

func TestCalculateNextSession_LastWeekOfPhaseIsDeload(t *testing.T) {
	// 7 complete weeks (28 sessions) → week 8 (last) = deload.
	phase, week, workoutNum, isDeload := CalculateNextSession(28, 8, 4)
	assert.Equal(t, 1, phase)
	assert.Equal(t, 8, week)
	assert.Equal(t, 1, workoutNum)
	assert.True(t, isDeload)
}

func TestCalculateNextSession_PhaseAdvances(t *testing.T) {
	// 32 sessions = 8 complete weeks = end of phase 1 → phase 2, week 1.
	phase, week, workoutNum, isDeload := CalculateNextSession(32, 8, 4)
	assert.Equal(t, 2, phase)
	assert.Equal(t, 1, week)
	assert.Equal(t, 1, workoutNum)
	assert.False(t, isDeload)
}

func TestCalculateNextSession_CustomWorkoutsPerWeek(t *testing.T) {
	// 5 workouts/week: after 5 sessions → week 2 workout 1.
	phase, week, workoutNum, isDeload := CalculateNextSession(5, 8, 5)
	assert.Equal(t, 1, phase)
	assert.Equal(t, 2, week)
	assert.Equal(t, 1, workoutNum)
	assert.False(t, isDeload)
}

func TestCalculateNextSession_ZeroDefaultsApplied(t *testing.T) {
	// Zero values should default to weeksPerPhase=8, workoutsPerWeek=4.
	phase, week, workoutNum, isDeload := CalculateNextSession(0, 0, 0)
	assert.Equal(t, 1, phase)
	assert.Equal(t, 1, week)
	assert.Equal(t, 1, workoutNum)
	assert.False(t, isDeload)
}

func TestCalculatePhaseAndWeek_WeekOne(t *testing.T) {
	start := time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC)
	now := start // day 0
	phase, week, isDeload := CalculatePhaseAndWeek(now, start, 8)
	assert.Equal(t, 1, phase)
	assert.Equal(t, 1, week)
	assert.False(t, isDeload)
}

func TestCalculatePhaseAndWeek_MidPhase(t *testing.T) {
	start := time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC)
	now := start.AddDate(0, 0, 21) // 3 weeks in → week 4
	phase, week, isDeload := CalculatePhaseAndWeek(now, start, 8)
	assert.Equal(t, 1, phase)
	assert.Equal(t, 4, week)
	assert.False(t, isDeload)
}

func TestCalculatePhaseAndWeek_LastWeekIsDeload(t *testing.T) {
	start := time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC)
	now := start.AddDate(0, 0, 49) // 7 weeks in → week 8 of phase 1 (last week)
	phase, week, isDeload := CalculatePhaseAndWeek(now, start, 8)
	assert.Equal(t, 1, phase)
	assert.Equal(t, 8, week)
	assert.True(t, isDeload)
}

func TestCalculatePhaseAndWeek_SecondPhase(t *testing.T) {
	start := time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC)
	now := start.AddDate(0, 0, 56) // 8 weeks in → week 1 of phase 2
	phase, week, isDeload := CalculatePhaseAndWeek(now, start, 8)
	assert.Equal(t, 2, phase)
	assert.Equal(t, 1, week)
	assert.False(t, isDeload)
}

func TestCalculatePhaseAndWeek_CustomWeeksPerPhase(t *testing.T) {
	start := time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC)
	// 5 weeks per phase; 4 weeks in → week 5 (last week) → deload
	now := start.AddDate(0, 0, 28)
	phase, week, isDeload := CalculatePhaseAndWeek(now, start, 5)
	assert.Equal(t, 1, phase)
	assert.Equal(t, 5, week)
	assert.True(t, isDeload)
}

func TestCalculatePhaseAndWeek_BeforeStartDate(t *testing.T) {
	start := time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC)
	now := start.AddDate(0, 0, -10) // before start → treated as day 0
	phase, week, isDeload := CalculatePhaseAndWeek(now, start, 8)
	assert.Equal(t, 1, phase)
	assert.Equal(t, 1, week)
	assert.False(t, isDeload)
}

func TestCalculatePhaseAndWeek_ZeroWeeksPerPhaseDefaultsToEight(t *testing.T) {
	start := time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC)
	now := start.AddDate(0, 0, 49) // week 8 with default 8
	phase, week, isDeload := CalculatePhaseAndWeek(now, start, 0)
	assert.Equal(t, 1, phase)
	assert.Equal(t, 8, week)
	assert.True(t, isDeload)
}
