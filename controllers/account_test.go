package controllers_test

import (
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSettings_Unauthenticated(t *testing.T) {
	w := getPath("/settings", nil)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/login", w.Header().Get("Location"))
}

func TestSettings_ShowsCurrentUnit(t *testing.T) {
	t.Cleanup(resetMocks)
	cookies := loginAs(t, "test_settings_get", "lb")

	w := getPath("/settings", cookies)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `value="lb"`)
	assert.Contains(t, w.Body.String(), `value="kg"`)
	assert.Regexp(t, `id="unit_lb"[^>]*checked`, w.Body.String())
	assert.NotRegexp(t, `id="unit_kg"[^>]*checked`, w.Body.String())
}

func TestSettings_ShowsKgCheckedForKgUser(t *testing.T) {
	t.Cleanup(resetMocks)
	cookies := loginAs(t, "test_settings_kg_get", "kg")

	w := getPath("/settings", cookies)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Regexp(t, `id="unit_kg"[^>]*checked`, w.Body.String())
	assert.NotRegexp(t, `id="unit_lb"[^>]*checked`, w.Body.String())
}

func TestSettingsPost_Unauthenticated(t *testing.T) {
	w := postForm("/settings", url.Values{"weight_unit": {"kg"}}, nil)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/login", w.Header().Get("Location"))
}

func TestSettingsPost_UpdateToKg(t *testing.T) {
	t.Cleanup(resetMocks)
	var calledUnit string
	mockUsers.UpdateWeightUnitFn = func(userID int64, unit string) error {
		calledUnit = unit
		return nil
	}
	cookies := loginAs(t, "test_settings_kg", "lb")

	w := postForm("/settings", url.Values{"weight_unit": {"kg"}}, cookies)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/settings", w.Header().Get("Location"))
	assert.Equal(t, "kg", calledUnit)

	// Follow the redirect — flash message should appear.
	allCookies := append(cookies, w.Result().Cookies()...)
	w2 := getPath("/settings", allCookies)
	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Contains(t, w2.Body.String(), "Settings saved.")
}

func TestSettingsPost_UpdateToLb(t *testing.T) {
	t.Cleanup(resetMocks)
	var calledUnit string
	mockUsers.UpdateWeightUnitFn = func(userID int64, unit string) error {
		calledUnit = unit
		return nil
	}
	cookies := loginAs(t, "test_settings_lb", "kg")

	w := postForm("/settings", url.Values{"weight_unit": {"lb"}}, cookies)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/settings", w.Header().Get("Location"))
	assert.Equal(t, "lb", calledUnit)
}

func TestSettingsPost_InvalidUnit(t *testing.T) {
	t.Cleanup(resetMocks)
	mockUsers.UpdateWeightUnitFn = func(userID int64, unit string) error {
		return errors.New("weight_unit must be lb or kg")
	}
	cookies := loginAs(t, "test_settings_bad", "lb")

	w := postForm("/settings", url.Values{"weight_unit": {"stone"}}, cookies)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid weight unit")
}
