package core

import (
	"encoding/json"
	"math"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestNewValueNativeCollectionsPreserveKinds(t *testing.T) {
	type bytesAlias []byte
	type valuesAlias map[string]*Value
	data := []byte("hello")
	value := NewStringValue("NaN")
	for _, tt := range []struct {
		name  string
		input any
		want  *Value
	}{
		{"bytes map", map[string][]byte{"data": data}, NewMapValue(map[string]*Value{"data": NewBytesValue(data)})},
		{"bytes list", [][]byte{data}, NewArrayValue(NewBytesValue(data))},
		{"named bytes", bytesAlias(data), NewBytesValue(data)},
		{"bytes pointer", &data, NewBytesValue(data)},
		{"value pointer", &value, value},
		{"named value map", valuesAlias{"text": value}, NewMapValue(map[string]*Value{"text": value})},
		{"floats", []float64{1, math.NaN(), math.Inf(1)}, NewArrayValue(NewFloat64Value(1), NewFloat64Value(math.NaN()), NewFloat64Value(math.Inf(1)))},
		{"float map", map[string]float64{"number": 1}, NewMapValue(map[string]*Value{"number": NewFloat64Value(1)})},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewValue(tt.input)
			require.NoError(t, err)
			require.True(t, proto.Equal(tt.want, got), "want %v, got %v", tt.want, got)
		})
	}
}

func TestNewValueRejectsInvalidNativeData(t *testing.T) {
	for _, input := range []any{
		map[int]string{1: "value"},
		map[string]string{"\xff": "value"},
		map[string]*Value{"\xff": NewStringValue("value")},
		map[string]string{"key": "\xff"},
		[]string{"\xff"},
	} {
		value, err := NewValue(input)
		require.Error(t, err, "%T", input)
		require.Nil(t, value)
	}
}

func TestNewValueRejectsNativeCycles(t *testing.T) {
	object := map[string]any{}
	object["self"] = object
	values := make([]any, 1)
	values[0] = values
	var pointer any
	pointer = &pointer
	for _, input := range []any{object, values, pointer} {
		value, err := NewValue(input)
		require.ErrorContains(t, err, "cyclic value")
		require.Nil(t, value)
	}
	_, err := NewObjectFromMap(object)
	require.ErrorContains(t, err, "cyclic value")
	_, err = NewValues(values)
	require.ErrorContains(t, err, "cyclic value")
	destination := NewObject().SetString("old", "value")
	require.Error(t, destination.From(object))
	require.Equal(t, "value", destination.GetString("old"))
}

func TestNewValueAllowsSharedAndOverlappingCollections(t *testing.T) {
	shared := map[string]any{"name": "shared"}
	value, err := NewValue(map[string]any{"a": shared, "b": shared})
	require.NoError(t, err)
	require.Equal(t, "shared", value.GetObject().GetObject("a").GetString("name"))
	require.Equal(t, "shared", value.GetObject().GetObject("b").GetString("name"))
	values := []any{"first", nil}
	values[1] = values[:1]
	value, err = NewValue(values)
	require.NoError(t, err)
	require.Equal(t, []any{"first", []any{"first"}}, value.AsInterface())
}

func TestValueArraysPreserveNilAndEmpty(t *testing.T) {
	for _, pair := range [][2]*Value{
		{NewIntArrayValue(), NewIntArrayValue([]int{}...)},
		{NewInt32ArrayValue(), NewInt32ArrayValue([]int32{}...)},
		{NewInt64ArrayValue(), NewInt64ArrayValue([]int64{}...)},
		{NewUintArrayValue(), NewUintArrayValue([]uint{}...)},
		{NewUint32ArrayValue(), NewUint32ArrayValue([]uint32{}...)},
		{NewUint64ArrayValue(), NewUint64ArrayValue([]uint64{}...)},
		{NewFloat32ArrayValue(), NewFloat32ArrayValue([]float32{}...)},
		{NewFloat64ArrayValue(), NewFloat64ArrayValue([]float64{}...)},
		{NewStringArrayValue(), NewStringArrayValue([]string{}...)},
		{NewObjectArrayValue(), NewObjectArrayValue([]*Object{}...)},
	} {
		require.Nil(t, pair[0].GetValues())
		require.NotNil(t, pair[1].GetValues())
		for i, want := range []string{`null`, `[]`} {
			data, err := jsoniter.Marshal(pair[i])
			require.NoError(t, err)
			require.Equal(t, want, string(data))
		}
	}
	for _, value := range []*Value{NewArrayValue(), NewArrayValue([]*Value{}...)} {
		for _, view := range []any{value.GetBoolArray(), value.GetIntArray(), value.GetInt64Array(), value.GetUintArray(), value.GetFloat64Array(), value.GetStringArray(), value.GetObjectArray()} {
			if value.GetValues() == nil {
				require.Nil(t, view)
			} else {
				require.NotNil(t, view)
				require.Empty(t, view)
			}
		}
	}
}

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
			value := NewStringValue("original")
			require.Error(t, jsoniter.UnmarshalFromString(input, value))
			require.Equal(t, "original", value.GetString())
		})
	}
}

func TestValueCodecRoundTrip(t *testing.T) {
	for _, value := range []*Value{
		NewFloat64Value(0), NewFloat64Value(1), NewFloat64Value(-1),
		NewFloat64Value(math.Copysign(0, -1)), NewFloat64Value(math.MaxFloat64),
		NewFloat64Value(1.23456789012345), NewFloat64Value(math.NaN()),
		NewFloat64Value(math.Inf(1)), NewFloat64Value(math.Inf(-1)),
		NewBytesValue([]byte("hello")), NewObjectValue(nil), NewValuesValue(nil), NewNullValue(),
		NewObjectValue(&Object{}), NewValuesValue(&Values{}),
		NewObjectValue(NewObject()), NewValuesValue(&Values{Vals: []*Value{}}),
	} {
		for _, codec := range []jsoniter.API{jsoniter.ConfigDefault, jsoniter.ConfigFastest, jsoniter.ConfigCompatibleWithStandardLibrary} {
			data, err := codec.Marshal(value)
			require.NoError(t, err)
			var decoded Value
			require.NoError(t, codec.Unmarshal(data, &decoded))
			if value.GetKind() == ValueKind_VALUE_KIND_NUMBER {
				require.Equal(t, value.GetKind(), decoded.GetKind())
				require.Equal(t, math.Signbit(value.GetNumberValue()), math.Signbit(decoded.GetNumberValue()))
			}
			if math.IsNaN(value.GetNumberValue()) {
				require.True(t, math.IsNaN(decoded.GetNumberValue()))
			} else {
				require.Equal(t, value.AsInterface(), decoded.AsInterface())
			}
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
