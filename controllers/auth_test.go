package controllers_test

import (
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

// --- Register ---

func TestRegisterPage(t *testing.T) {
	w := getPath("/register", nil)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Create Account")
}

func TestRegisterPost_Success(t *testing.T) {
	t.Cleanup(resetMocks)

	w := postForm("/register", url.Values{
		"username":         {"test_reg_ok"},
		"password":         {"password123"},
		"confirm_password": {"password123"},
	}, nil)

	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/login", w.Header().Get("Location"))
}

func TestRegisterPost_EmptyUsername(t *testing.T) {
	w := postForm("/register", url.Values{
		"username":         {""},
		"password":         {"password123"},
		"confirm_password": {"password123"},
	}, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "required")
}

func TestRegisterPost_EmptyPassword(t *testing.T) {
	w := postForm("/register", url.Values{
		"username":         {"test_nopw_post"},
		"password":         {""},
		"confirm_password": {""},
	}, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "required")
}

func TestRegisterPost_ShortPassword(t *testing.T) {
	w := postForm("/register", url.Values{
		"username":         {"test_shortpw"},
		"password":         {"short"},
		"confirm_password": {"short"},
	}, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "8 characters")
}

func TestRegisterPost_PasswordMismatch(t *testing.T) {
	w := postForm("/register", url.Values{
		"username":         {"test_mismatch"},
		"password":         {"password123"},
		"confirm_password": {"different123"},
	}, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "do not match")
}

func TestRegisterPost_DuplicateUsername(t *testing.T) {
	t.Cleanup(resetMocks)
	setCreateFnError(errors.New("unique constraint violated"))

	w := postForm("/register", url.Values{
		"username":         {"test_dup_reg"},
		"password":         {"password123"},
		"confirm_password": {"password123"},
	}, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "already taken")
}

// --- Login ---

func TestLoginPage(t *testing.T) {
	w := getPath("/login", nil)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Log In")
}

func TestLoginPost_Success(t *testing.T) {
	t.Cleanup(resetMocks)
	setGetByUsernameReturnsUser("lb")

	w := postForm("/login", url.Values{
		"username": {"test_login_ok"},
		"password": {"password123"},
	}, nil)

	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/dashboard", w.Header().Get("Location"))
}

func TestLoginPost_WrongPassword(t *testing.T) {
	t.Cleanup(resetMocks)
	setGetByUsernameReturnsUser("lb")

	w := postForm("/login", url.Values{
		"username": {"test_login_bad"},
		"password": {"wrongpassword"},
	}, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid username or password")
}

func TestLoginPost_UnknownUsername(t *testing.T) {
	// Default mock returns "not found" error — no setup needed.
	w := postForm("/login", url.Values{
		"username": {"nobody_here"},
		"password": {"password123"},
	}, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid username or password")
}

// --- Logout ---

func TestLogout(t *testing.T) {
	t.Cleanup(resetMocks)
	cookies := loginAs(t, "test_logout", "lb")

	w := getPath("/logout", cookies)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/login", w.Header().Get("Location"))
}
