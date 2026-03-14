package controllers_test

import (
	"net/http"
	"testing"

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
	t.Cleanup(resetMocks)
	cookies := loginAs(t, "test_home", "lb")

	w := getPath("/", cookies)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/dashboard", w.Header().Get("Location"))
}
