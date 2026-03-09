package controllers_test

import (
	"net/http"
	"testing"

	"myGymPal/models"

	"github.com/stretchr/testify/assert"
)

func TestHomepage_Unauthenticated(t *testing.T) {
	w := getPath("/", nil)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `href="/login"`)
	assert.Contains(t, w.Body.String(), `href="/register"`)
	assert.NotContains(t, w.Body.String(), `href="/logout"`)
}

func TestHomepage_Authenticated(t *testing.T) {
	cookies := loginAs(t, "test_home", "password123")
	t.Cleanup(func() { models.DeleteUserByUsername("test_home") })

	w := getPath("/", cookies)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `href="/logout"`)
	assert.NotContains(t, w.Body.String(), `href="/login"`)
	assert.NotContains(t, w.Body.String(), `href="/register"`)
}
