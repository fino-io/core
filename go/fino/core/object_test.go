package core

import (
	"fmt"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
)

func TestNewObject(t *testing.T) {
	v := map[string]any{"name": "apple", "age": 20}
	got, err := NewObjectFromMap(v)
	assert.NoError(t, err)
	m := got.AsMap()
	assert.Equal(t, "apple", m["name"])
	assert.EqualValues(t, 20, m["age"])
}

func TestObject_AsMap(t *testing.T) {
	v := map[string]any{
		"map": map[string]any{
			"inmap": struct{ a int }{a: 11},
		},
	}
	got, err := NewObjectFromMap(v)
	assert.NoError(t, err)
	bytes, err := jsoniter.Marshal(got)
	assert.NoError(t, err)
	fmt.Println(string(bytes))
}
