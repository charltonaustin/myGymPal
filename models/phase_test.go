package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testUserForPhase(t *testing.T, username string) *User {
	t.Helper()
	_ = DeleteUserByUsername(username)
	u, err := CreateUser(username, "password123", "lb")
	require.NoError(t, err)
	t.Cleanup(func() { DeleteUserByUsername(username) })
	return u
}

func TestCreateProgram_CreatesPhases(t *testing.T) {
	u := testUserForPhase(t, "phase_created")

	p, err := CreateProgram(u.ID, "Block A", testStartDate, 3, 8, 4, 10, 12)
	require.NoError(t, err)

	phases, err := GetPhasesByProgramID(p.ID)
	require.NoError(t, err)
	require.Len(t, phases, 3)
	assert.Equal(t, 1, phases[0].PhaseNumber)
	assert.Equal(t, 2, phases[1].PhaseNumber)
	assert.Equal(t, 3, phases[2].PhaseNumber)
	// Rep ranges are initialised from the defaults passed to CreateProgram.
	assert.Equal(t, 10, phases[0].RepMin)
	assert.Equal(t, 12, phases[0].RepMax)
}

func TestGetPhasesByProgramID(t *testing.T) {
	u := testUserForPhase(t, "phase_list")

	p, err := CreateProgram(u.ID, "Block B", testStartDate, 4, 8, 4, 10, 12)
	require.NoError(t, err)

	phases, err := GetPhasesByProgramID(p.ID)
	require.NoError(t, err)
	assert.Len(t, phases, 4)
}

func TestUpdatePhaseRepRanges_Success(t *testing.T) {
	u := testUserForPhase(t, "phase_update_ok")

	p, err := CreateProgram(u.ID, "Block C", testStartDate, 2, 8, 4, 10, 12)
	require.NoError(t, err)

	updates := []PhaseUpdate{
		{PhaseNumber: 1, RepMin: 10, RepMax: 12},
		{PhaseNumber: 2, RepMin: 8, RepMax: 10},
	}
	require.NoError(t, UpdatePhaseRepRanges(p.ID, updates))

	phases, err := GetPhasesByProgramID(p.ID)
	require.NoError(t, err)
	assert.Equal(t, 10, phases[0].RepMin)
	assert.Equal(t, 12, phases[0].RepMax)
	assert.Equal(t, 8, phases[1].RepMin)
	assert.Equal(t, 10, phases[1].RepMax)
}

func TestUpdatePhaseRepRanges_ZeroMin(t *testing.T) {
	u := testUserForPhase(t, "phase_update_zero_min")

	p, err := CreateProgram(u.ID, "Block D", testStartDate, 1, 8, 4, 10, 12)
	require.NoError(t, err)

	err = UpdatePhaseRepRanges(p.ID, []PhaseUpdate{{PhaseNumber: 1, RepMin: 0, RepMax: 10}})
	assert.Error(t, err)
}

func TestUpdatePhaseRepRanges_MaxLessThanMin(t *testing.T) {
	u := testUserForPhase(t, "phase_update_bad_range")

	p, err := CreateProgram(u.ID, "Block E", testStartDate, 1, 8, 4, 10, 12)
	require.NoError(t, err)

	err = UpdatePhaseRepRanges(p.ID, []PhaseUpdate{{PhaseNumber: 1, RepMin: 12, RepMax: 10}})
	assert.Error(t, err)
}

func TestDeleteProgram_CascadesPhases(t *testing.T) {
	u := testUserForPhase(t, "phase_cascade")

	p, err := CreateProgram(u.ID, "Block F", testStartDate, 3, 8, 4, 10, 12)
	require.NoError(t, err)

	require.NoError(t, DeleteProgram(p.ID, u.ID))

	phases, err := GetPhasesByProgramID(p.ID)
	require.NoError(t, err)
	assert.Empty(t, phases)
}
