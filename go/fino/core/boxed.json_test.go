package core

import (
	"encoding/json"
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

func TestBoxedJSONFailurePreservesValue(t *testing.T) {
	for _, tt := range []struct {
		name  string
		value json.Unmarshaler
		input string
	}{
		{"scalar", &Int32Value{Val: 9}, `1.5`},
		{"array", &Int32Values{Vals: []int32{9}}, `[1,"bad"]`},
		{"strings", &StringValues{Vals: []string{"old"}}, `["new",1]`},
		{"map", &StringMap{Vals: map[string]string{"old": "kept"}}, `{"new":"value","bad":1}`},
		{"nested map", &StringsMap{Vals: map[string]*StringValues{"old": {Vals: []string{"kept"}}}}, `{"new":["value"],"bad":1}`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			before, err := json.Marshal(tt.value)
			require.NoError(t, err)
			for _, decode := range []func([]byte, any) error{json.Unmarshal, jsoniter.Unmarshal} {
				require.Error(t, decode([]byte(tt.input), tt.value))
				after, err := json.Marshal(tt.value)
				require.NoError(t, err)
				require.JSONEq(t, string(before), string(after))
			}
		})
	}
}

func TestBoxedJSONNullClearsValue(t *testing.T) {
	for _, tt := range []struct {
		value json.Unmarshaler
		want  string
	}{
		{&BoolValue{Val: true}, `false`},
		{&Int32Value{Val: 9}, `0`},
		{&Int64Value{Val: 9}, `0`},
		{&Uint32Value{Val: 9}, `0`},
		{&Uint64Value{Val: 9}, `0`},
		{&Float32Value{Val: 9}, `0`},
		{&Float64Value{Val: 9}, `0`},
		{&StringValue{Val: "old"}, `""`},
		{&BytesValue{Val: []byte("old")}, `null`},
		{&StringValues{Vals: []string{"old"}}, `null`},
		{&StringMap{Vals: map[string]string{"old": "value"}}, `null`},
	} {
		require.NoError(t, json.Unmarshal([]byte(`null`), tt.value))
		data, err := json.Marshal(tt.value)
		require.NoError(t, err)
		require.JSONEq(t, tt.want, string(data), "%T", tt.value)
	}
}

func TestStringMapJSONReplacesValue(t *testing.T) {
	values := &StringMap{Vals: map[string]string{"old": "value"}}
	for _, decode := range []func([]byte, any) error{json.Unmarshal, jsoniter.Unmarshal} {
		require.NoError(t, decode([]byte(`{"new":"value"}`), values))
		require.Equal(t, map[string]string{"new": "value"}, values.Vals)
		values.Vals["old"] = "value"
	}
	require.Error(t, values.UnmarshalJSON(nil))
	require.Error(t, (&StringsMap{}).UnmarshalJSON(nil))
}
