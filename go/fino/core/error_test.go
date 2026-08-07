package core

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewErrorFrom(t *testing.T) {
	err1 := NewErrorFrom(500, "this is an error")
	assert.Error(t, err1)
}

func TestErrorStatusCodeFallsBackToInternalServerError(t *testing.T) {
	assert.Equal(t, http.StatusInternalServerError, NewErrorFrom(600121001, "business error").StatusCode())
	assert.Equal(t, http.StatusInternalServerError, (*Error)(nil).StatusCode())
}
