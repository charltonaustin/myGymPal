package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
