package controllers_test

import (
	"net/http"
	"net/url"
	"testing"

	"myGymPal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSettings_Unauthenticated(t *testing.T) {
	w := getPath("/settings", nil)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/login", w.Header().Get("Location"))
}

func TestSettings_ShowsCurrentUnit(t *testing.T) {
	cookies := loginAs(t, "test_settings_get", "password123")
	t.Cleanup(func() { models.DeleteUserByUsername("test_settings_get") })

	w := getPath("/settings", cookies)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `value="lb"`)
	assert.Contains(t, w.Body.String(), `value="kg"`)
	// default is lb, so lb radio should be checked and kg should not
	assert.Regexp(t, `id="unit_lb"[^>]*checked`, w.Body.String())
	assert.NotRegexp(t, `id="unit_kg"[^>]*checked`, w.Body.String())
}

func TestSettings_ShowsKgCheckedForKgUser(t *testing.T) {
	_ = models.DeleteUserByUsername("test_settings_kg_get")
	_, err := models.CreateUser("test_settings_kg_get", "password123", "kg")
	require.NoError(t, err)
	t.Cleanup(func() { models.DeleteUserByUsername("test_settings_kg_get") })

	// Log in directly — loginAs always creates with "lb" and would overwrite the kg user.
	w := postForm("/login", url.Values{
		"username": {"test_settings_kg_get"},
		"password": {"password123"},
	}, nil)
	require.Equal(t, http.StatusFound, w.Code)
	cookies := w.Result().Cookies()

	w2 := getPath("/settings", cookies)
	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Regexp(t, `id="unit_kg"[^>]*checked`, w2.Body.String())
	assert.NotRegexp(t, `id="unit_lb"[^>]*checked`, w2.Body.String())
}

func TestSettingsPost_Unauthenticated(t *testing.T) {
	w := postForm("/settings", url.Values{"weight_unit": {"kg"}}, nil)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/login", w.Header().Get("Location"))
}

func TestSettingsPost_UpdateToKg(t *testing.T) {
	cookies := loginAs(t, "test_settings_kg", "password123")
	t.Cleanup(func() { models.DeleteUserByUsername("test_settings_kg") })

	w := postForm("/settings", url.Values{"weight_unit": {"kg"}}, cookies)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/settings", w.Header().Get("Location"))

	// Follow the redirect — success message should appear
	allCookies := append(cookies, w.Result().Cookies()...)
	w2 := getPath("/settings", allCookies)
	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Contains(t, w2.Body.String(), "Settings saved.")

	user, err := models.GetUserByUsername("test_settings_kg")
	require.NoError(t, err)
	assert.Equal(t, "kg", user.WeightUnit)
}

func TestSettingsPost_UpdateToLb(t *testing.T) {
	_ = models.DeleteUserByUsername("test_settings_lb")
	user, err := models.CreateUser("test_settings_lb", "password123", "kg")
	require.NoError(t, err)
	t.Cleanup(func() { models.DeleteUserByUsername("test_settings_lb") })
	_ = user

	cookies := loginAs(t, "test_settings_lb", "password123")

	w := postForm("/settings", url.Values{"weight_unit": {"lb"}}, cookies)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/settings", w.Header().Get("Location"))

	updated, err := models.GetUserByUsername("test_settings_lb")
	require.NoError(t, err)
	assert.Equal(t, "lb", updated.WeightUnit)
}

func TestSettingsPost_InvalidUnit(t *testing.T) {
	cookies := loginAs(t, "test_settings_bad", "password123")
	t.Cleanup(func() { models.DeleteUserByUsername("test_settings_bad") })

	w := postForm("/settings", url.Values{"weight_unit": {"stone"}}, cookies)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid weight unit")
}
