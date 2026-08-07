package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNotFoundError(t *testing.T) {
	key := "key"
	// args := []any{"key", "value"}
	notFoundError := NewNotFoundError(key)
	assert.Equal(t, notFoundError, NewNotFoundError(key))
}
