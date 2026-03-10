package controllers_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDashboard_Unauthenticated(t *testing.T) {
	w := getPath("/dashboard", nil)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/login", w.Header().Get("Location"))
}

func TestDashboard_Authenticated(t *testing.T) {
	t.Cleanup(resetMocks)
	cookies := loginAs(t, "test_dash", "lb")

	w := getPath("/dashboard", cookies)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "test_dash")
}
