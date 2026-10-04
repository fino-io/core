package core

import (
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestStringsMapNull(t *testing.T) {
	values := &StringsMap{Vals: map[string]*StringValues{"old": {Vals: []string{"stale"}}}}
	require.NoError(t, jsoniter.UnmarshalFromString(`null`, values))
	require.Nil(t, values.Vals)
	data, err := jsoniter.MarshalToString(values)
	require.NoError(t, err)
	require.Equal(t, `null`, data)
	require.NoError(t, jsoniter.UnmarshalFromString(`{}`, values))
	require.NotNil(t, values.Vals)
	data, err = jsoniter.MarshalToString(values)
	require.NoError(t, err)
	require.Equal(t, `{}`, data)
}
