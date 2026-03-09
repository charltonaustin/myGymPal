package models

import (
	"os"
	"testing"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// Tests run from the package directory; move to project root so
	// conf/app.conf and migrations/ are found at their expected paths.
	if err := os.Chdir(".."); err != nil {
		panic(err)
	}
	beego.LoadAppConfig("ini", "conf/app.test.conf")
	if err := Init(); err != nil {
		panic("failed to init DB: " + err.Error())
	}
	os.Exit(m.Run())
}

func TestCreateUser_Success(t *testing.T) {
	user, err := CreateUser("test_create", "password123", "lb")
	require.NoError(t, err)
	t.Cleanup(func() { DeleteUserByUsername("test_create") })

	assert.Equal(t, "test_create", user.Username)
	assert.Equal(t, "lb", user.WeightUnit)
	assert.NotEmpty(t, user.PasswordHash)
	assert.NotEqual(t, "password123", user.PasswordHash, "password must be hashed")
	assert.Positive(t, user.ID)
}

func TestCreateUser_DuplicateUsername(t *testing.T) {
	_, err := CreateUser("test_dup", "password123", "lb")
	require.NoError(t, err)
	t.Cleanup(func() { DeleteUserByUsername("test_dup") })

	_, err = CreateUser("test_dup", "differentpassword", "kg")
	assert.Error(t, err)
}

func TestCreateUser_EmptyUsername(t *testing.T) {
	_, err := CreateUser("", "password123", "lb")
	assert.Error(t, err)
}

func TestCreateUser_EmptyPassword(t *testing.T) {
	_, err := CreateUser("test_nopw", "", "lb")
	assert.Error(t, err)
}

func TestGetUserByUsername_Found(t *testing.T) {
	_, err := CreateUser("test_get", "password123", "kg")
	require.NoError(t, err)
	t.Cleanup(func() { DeleteUserByUsername("test_get") })

	user, err := GetUserByUsername("test_get")
	require.NoError(t, err)
	assert.Equal(t, "test_get", user.Username)
	assert.Equal(t, "kg", user.WeightUnit)
}

func TestGetUserByUsername_NotFound(t *testing.T) {
	_, err := GetUserByUsername("does_not_exist")
	assert.Error(t, err)
}

func TestCheckPassword_Correct(t *testing.T) {
	user, err := CreateUser("test_pw_ok", "correctpassword", "lb")
	require.NoError(t, err)
	t.Cleanup(func() { DeleteUserByUsername("test_pw_ok") })

	assert.True(t, user.CheckPassword("correctpassword"))
}

func TestCheckPassword_Wrong(t *testing.T) {
	user, err := CreateUser("test_pw_wrong", "correctpassword", "lb")
	require.NoError(t, err)
	t.Cleanup(func() { DeleteUserByUsername("test_pw_wrong") })

	assert.False(t, user.CheckPassword("wrongpassword"))
}

func TestGetUserByID_Found(t *testing.T) {
	user, err := CreateUser("test_get_id", "password123", "kg")
	require.NoError(t, err)
	t.Cleanup(func() { DeleteUserByUsername("test_get_id") })

	found, err := GetUserByID(user.ID)
	require.NoError(t, err)
	assert.Equal(t, user.ID, found.ID)
	assert.Equal(t, "test_get_id", found.Username)
}

func TestGetUserByID_NotFound(t *testing.T) {
	_, err := GetUserByID(-1)
	assert.Error(t, err)
}

func TestUpdateWeightUnit_Valid(t *testing.T) {
	user, err := CreateUser("test_update_unit", "password123", "lb")
	require.NoError(t, err)
	t.Cleanup(func() { DeleteUserByUsername("test_update_unit") })

	require.NoError(t, UpdateWeightUnit(user.ID, "kg"))
	updated, err := GetUserByID(user.ID)
	require.NoError(t, err)
	assert.Equal(t, "kg", updated.WeightUnit)

	require.NoError(t, UpdateWeightUnit(user.ID, "lb"))
	updated, err = GetUserByID(user.ID)
	require.NoError(t, err)
	assert.Equal(t, "lb", updated.WeightUnit)
}

func TestUpdateWeightUnit_Invalid(t *testing.T) {
	user, err := CreateUser("test_update_bad_unit", "password123", "lb")
	require.NoError(t, err)
	t.Cleanup(func() { DeleteUserByUsername("test_update_bad_unit") })

	assert.Error(t, UpdateWeightUnit(user.ID, "stone"))
}
