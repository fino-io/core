package core

import (
	"encoding/json"
	"math"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/require"
)

func TestNewValueCompositeAndNil(t *testing.T) {
	var pointer *int
	value, err := NewValue(pointer)
	require.NoError(t, err)
	require.Equal(t, ValueKind_VALUE_KIND_NULL, value.GetKind())

	for _, input := range []any{
		struct {
			ID uint64 `json:"id"`
		}{math.MaxUint64},
		map[string]uint64{"id": math.MaxUint64},
	} {
		value, err := NewValue(input)
		require.NoError(t, err)
		require.Equal(t, uint64(math.MaxUint64), value.GetObject().GetUint64("id"))
	}
	value, err = NewValue([]int{1, 2})
	require.NoError(t, err)
	require.Equal(t, []int{1, 2}, value.GetIntArray())
	value, err = NewValue(json.Number("18446744073709551615"))
	require.NoError(t, err)
	require.Equal(t, uint64(math.MaxUint64), value.GetUint64())
}

func TestValueCodecInvalidInput(t *testing.T) {
	for _, input := range []string{`"b64.!"`, `{"nested":"b64.!"}`, `["b64.!"]`, `1e999`, `{"nested":`} {
		t.Run(input, func(t *testing.T) {
			var value Value
			require.Error(t, jsoniter.UnmarshalFromString(input, &value))
		})
	}
}

func TestValueCodecRoundTrip(t *testing.T) {
	for _, value := range []*Value{
		NewFloat64Value(1.23456789012345), NewFloat64Value(math.NaN()),
		NewFloat64Value(math.Inf(1)), NewFloat64Value(math.Inf(-1)),
		NewBytesValue([]byte("hello")), NewObjectValue(nil), NewValuesValue(nil), NewNullValue(),
		NewObjectValue(&Object{}), NewValuesValue(&Values{}),
		NewObjectValue(NewObject()), NewValuesValue(&Values{Vals: []*Value{}}),
	} {
		data, err := jsoniter.Marshal(value)
		require.NoError(t, err)
		var decoded Value
		require.NoError(t, jsoniter.Unmarshal(data, &decoded))
		if math.IsNaN(value.GetNumberValue()) {
			require.True(t, math.IsNaN(decoded.GetNumberValue()))
		} else {
			require.Equal(t, value.AsInterface(), decoded.AsInterface())
		}
	}
}

func TestCompositeCodecsRejectWrongTypes(t *testing.T) {
	for _, target := range []any{&Object{}, &Values{}} {
		require.Error(t, jsoniter.UnmarshalFromString(`true`, target))
		require.Error(t, jsoniter.UnmarshalFromString(`{"bad":"b64.!"}`, target))
	}
	values := &Values{Vals: []*Value{NewStringValue("old")}}
	require.NoError(t, jsoniter.UnmarshalFromString(`null`, values))
	require.Nil(t, values.Vals)
	values.Vals = []*Value{NewIntValue(1)}
	data, err := json.Marshal(values)
	require.NoError(t, err)
	require.Equal(t, `[1]`, string(data))
	require.NoError(t, json.Unmarshal([]byte(`[2]`), values))
	require.Equal(t, int64(2), values.Vals[0].GetInt64())
}
