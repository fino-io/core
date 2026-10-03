package core

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewErrorFrom(t *testing.T) {
	err1 := NewErrorFrom(500, "this is an error")
	assert.Error(t, err1)
}

func TestErrorJSONPreservesStatus(t *testing.T) {
	original := NewErrorFrom(NotFound.Code, "missing")
	data, err := jsoniter.Marshal(original)
	require.NoError(t, err)
	var decoded Error
	require.NoError(t, jsoniter.Unmarshal(data, &decoded))
	require.Equal(t, original.StatusCode(), decoded.StatusCode())
	require.Equal(t, original.Code.Name, decoded.Code.Name)
}

func TestErrorMatching(t *testing.T) {
	err := NewErrorFrom(404, "missing")
	wrapped := fmt.Errorf("lookup: %w", err)
	require.True(t, IsError(wrapped))
	require.Same(t, err, AsError(wrapped))
	require.True(t, errors.Is(wrapped, &Error{}))
	require.False(t, errors.Is(err, errors.New("other")))
	require.False(t, IsError(nil))
	require.False(t, IsError(errors.New("other")))

	typed := NewNotFoundError("missing")
	require.True(t, IsNotFoundError(fmt.Errorf("lookup: %w", typed)))
	require.True(t, IsError(typed))
	require.Same(t, typed.ToError(), AsError(typed))
	require.False(t, IsBadRequestError(typed))
	require.False(t, errors.Is(typed, fmt.Errorf("target: %w", typed)))
}

func TestNewErrorFromDoesNotShareCode(t *testing.T) {
	first := NewErrorFrom(NotFound.Code, "first")
	second := NewErrorFrom(NotFound.Code, "second")
	first.Code.Name = "changed"
	require.Equal(t, NotFound.Name, second.Code.Name)
	require.NotSame(t, NotFound, first.Code)
	typed := NewNotFoundError("missing")
	typed.ToError().Code.Name = "changed"
	require.Equal(t, NotFound.Name, NewNotFoundError("missing").ToError().Code.Name)
}

func TestErrorStatusCodeFallsBackToInternalServerError(t *testing.T) {
	assert.Equal(t, http.StatusInternalServerError, NewErrorFrom(600121001, "business error").StatusCode())
	assert.Equal(t, http.StatusInternalServerError, (*Error)(nil).StatusCode())
}
