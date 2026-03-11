package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testStartDate = time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC)

func testUserForProgram(t *testing.T, username string) *User {
	t.Helper()
	_ = DeleteUserByUsername(username)
	u, err := CreateUser(username, "password123", "lb")
	require.NoError(t, err)
	t.Cleanup(func() { DeleteUserByUsername(username) })
	return u
}

func TestCreateProgram_Success(t *testing.T) {
	u := testUserForProgram(t, "prog_create_ok")

	p, err := CreateProgram(u.ID, "Hypertrophy Block", testStartDate, 4, 8, 4, 10, 12)
	require.NoError(t, err)
	t.Cleanup(func() { DeleteProgram(p.ID, u.ID) })

	assert.Equal(t, "Hypertrophy Block", p.Name)
	assert.Equal(t, 4, p.NumPhases)
	assert.Equal(t, 8, p.WeeksPerPhase)
	assert.Equal(t, u.ID, p.UserID)
	assert.Positive(t, p.ID)
}

func TestCreateProgram_EmptyName(t *testing.T) {
	u := testUserForProgram(t, "prog_empty_name")

	_, err := CreateProgram(u.ID, "", testStartDate, 4, 8, 4, 10, 12)
	assert.Error(t, err)
}

func TestCreateProgram_ZeroPhases(t *testing.T) {
	u := testUserForProgram(t, "prog_zero_phases")

	_, err := CreateProgram(u.ID, "Bad Program", testStartDate, 0, 8, 4, 10, 12)
	assert.Error(t, err)
}

func TestCreateProgram_NegativePhases(t *testing.T) {
	u := testUserForProgram(t, "prog_neg_phases")

	_, err := CreateProgram(u.ID, "Bad Program", testStartDate, -1, 8, 4, 10, 12)
	assert.Error(t, err)
}

func TestCreateProgram_ZeroWeeksPerPhase(t *testing.T) {
	u := testUserForProgram(t, "prog_zero_weeks")

	_, err := CreateProgram(u.ID, "Bad Program", testStartDate, 4, 0, 4, 10, 12)
	assert.Error(t, err)
}

func TestCreateProgram_NegativeWeeksPerPhase(t *testing.T) {
	u := testUserForProgram(t, "prog_neg_weeks")

	_, err := CreateProgram(u.ID, "Bad Program", testStartDate, 4, -1, 4, 10, 12)
	assert.Error(t, err)
}

func TestGetProgramsByUserID(t *testing.T) {
	u := testUserForProgram(t, "prog_list_user")

	p1, err := CreateProgram(u.ID, "Program A", testStartDate, 3, 8, 4, 10, 12)
	require.NoError(t, err)
	p2, err := CreateProgram(u.ID, "Program B", testStartDate.AddDate(0, 0, 7), 2, 8, 4, 10, 12)
	require.NoError(t, err)
	t.Cleanup(func() {
		DeleteProgram(p1.ID, u.ID)
		DeleteProgram(p2.ID, u.ID)
	})

	programs, err := GetProgramsByUserID(u.ID)
	require.NoError(t, err)
	assert.Len(t, programs, 2)
}

func TestGetProgramsByUserID_OtherUsersNotIncluded(t *testing.T) {
	u1 := testUserForProgram(t, "prog_isolation_u1")
	u2 := testUserForProgram(t, "prog_isolation_u2")

	p, err := CreateProgram(u1.ID, "U1 Program", testStartDate, 2, 8, 4, 10, 12)
	require.NoError(t, err)
	t.Cleanup(func() { DeleteProgram(p.ID, u1.ID) })

	programs, err := GetProgramsByUserID(u2.ID)
	require.NoError(t, err)
	assert.Empty(t, programs)
}

func TestGetProgramByID_Found(t *testing.T) {
	u := testUserForProgram(t, "prog_get_by_id")

	p, err := CreateProgram(u.ID, "My Program", testStartDate, 5, 8, 4, 10, 12)
	require.NoError(t, err)
	t.Cleanup(func() { DeleteProgram(p.ID, u.ID) })

	found, err := GetProgramByID(p.ID, u.ID)
	require.NoError(t, err)
	assert.Equal(t, p.ID, found.ID)
	assert.Equal(t, "My Program", found.Name)
}

func TestGetProgramByID_WrongUser(t *testing.T) {
	u1 := testUserForProgram(t, "prog_wrong_user_u1")
	u2 := testUserForProgram(t, "prog_wrong_user_u2")

	p, err := CreateProgram(u1.ID, "U1 Program", testStartDate, 2, 8, 4, 10, 12)
	require.NoError(t, err)
	t.Cleanup(func() { DeleteProgram(p.ID, u1.ID) })

	_, err = GetProgramByID(p.ID, u2.ID)
	assert.Error(t, err)
}

func TestDeleteProgram(t *testing.T) {
	u := testUserForProgram(t, "prog_delete")

	p, err := CreateProgram(u.ID, "To Delete", testStartDate, 1, 8, 4, 10, 12)
	require.NoError(t, err)

	require.NoError(t, DeleteProgram(p.ID, u.ID))

	programs, err := GetProgramsByUserID(u.ID)
	require.NoError(t, err)
	assert.Empty(t, programs)
}

func TestDeleteProgram_WrongUser(t *testing.T) {
	u1 := testUserForProgram(t, "prog_del_wrong_u1")
	u2 := testUserForProgram(t, "prog_del_wrong_u2")

	p, err := CreateProgram(u1.ID, "Protected", testStartDate, 2, 8, 4, 10, 12)
	require.NoError(t, err)
	t.Cleanup(func() { DeleteProgram(p.ID, u1.ID) })

	err = DeleteProgram(p.ID, u2.ID)
	assert.Error(t, err)
}
