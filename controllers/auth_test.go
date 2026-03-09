package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"myGymPal/models"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// postForm is a helper that submits a form POST and returns the response.
func postForm(path string, data url.Values, cookies []*http.Cookie) *httptest.ResponseRecorder {
	r, _ := http.NewRequest("POST", path, strings.NewReader(data.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for _, c := range cookies {
		r.AddCookie(c)
	}
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	return w
}

func getPath(path string, cookies []*http.Cookie) *httptest.ResponseRecorder {
	r, _ := http.NewRequest("GET", path, nil)
	for _, c := range cookies {
		r.AddCookie(c)
	}
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	return w
}

// loginAs creates a user, logs in, and returns the session cookies.
// The caller is responsible for cleanup via t.Cleanup.
func loginAs(t *testing.T, username, password string) []*http.Cookie {
	t.Helper()
	// Remove any leftover user from a previous failed run before creating.
	_ = models.DeleteUserByUsername(username)
	_, err := models.CreateUser(username, password, "lb")
	require.NoError(t, err)

	w := postForm("/login", url.Values{
		"username": {username},
		"password": {password},
	}, nil)
	require.Equal(t, http.StatusFound, w.Code)
	return w.Result().Cookies()
}

// --- Register ---

func TestRegisterPage(t *testing.T) {
	w := getPath("/register", nil)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Create Account")
}

func TestRegisterPost_Success(t *testing.T) {
	t.Cleanup(func() { models.DeleteUserByUsername("test_reg_ok") })

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
	_, err := models.CreateUser("test_dup_reg", "password123", "lb")
	require.NoError(t, err)
	t.Cleanup(func() { models.DeleteUserByUsername("test_dup_reg") })

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
	_, err := models.CreateUser("test_login_ok", "password123", "lb")
	require.NoError(t, err)
	t.Cleanup(func() { models.DeleteUserByUsername("test_login_ok") })

	w := postForm("/login", url.Values{
		"username": {"test_login_ok"},
		"password": {"password123"},
	}, nil)

	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/dashboard", w.Header().Get("Location"))
}

func TestLoginPost_WrongPassword(t *testing.T) {
	_, err := models.CreateUser("test_login_bad", "password123", "lb")
	require.NoError(t, err)
	t.Cleanup(func() { models.DeleteUserByUsername("test_login_bad") })

	w := postForm("/login", url.Values{
		"username": {"test_login_bad"},
		"password": {"wrongpassword"},
	}, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid username or password")
}

func TestLoginPost_UnknownUsername(t *testing.T) {
	w := postForm("/login", url.Values{
		"username": {"nobody_here"},
		"password": {"password123"},
	}, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid username or password")
}

// --- Logout ---

func TestLogout(t *testing.T) {
	cookies := loginAs(t, "test_logout", "password123")
	t.Cleanup(func() { models.DeleteUserByUsername("test_logout") })

	w := getPath("/logout", cookies)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/login", w.Header().Get("Location"))
}
