package core

import (
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
)

func TestBoolValue_JSON(t *testing.T) {
	val := &BoolValue{Val: true}
	str, err := jsoniter.MarshalToString(val)
	assert.NoError(t, err)
	assert.Equal(t, "true", str)

	val2 := &BoolValue{}
	err = jsoniter.UnmarshalFromString("true", val2)
	assert.NoError(t, err)
	assert.EqualValues(t, true, val2.Val)
}

func TestStringMap_JSON(t *testing.T) {
	s := `{"foo":"bar"}`

	val := &StringMap{Vals: map[string]string{"foo": "bar"}}
	str, err := jsoniter.MarshalToString(val)
	assert.NoError(t, err)
	assert.Equal(t, s, str)

	val2 := &StringMap{}
	err = jsoniter.UnmarshalFromString(s, val2)
	assert.NoError(t, err)
	assert.Equal(t, "bar", val2.Vals["foo"])
}

func TestStringsMap_JSON(t *testing.T) {
	s := `{"foo":["bar","ping"]}`
	sv := &StringValues{Vals: []string{"bar", "ping"}}

	val := &StringsMap{Vals: map[string]*StringValues{"foo": sv}}
	str, err := jsoniter.MarshalToString(val)
	assert.NoError(t, err)
	assert.Equal(t, s, str)

	val2 := &StringsMap{}
	err = jsoniter.UnmarshalFromString(s, val2)
	assert.NoError(t, err)
	assert.EqualValues(t, sv, val2.Vals["foo"])
}
