package controllers_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorPage(t *testing.T) {
	w := getPath("/error", nil)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Something went wrong")
}
