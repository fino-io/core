package core

import (
	"bytes"
	"encoding/json"
	"math"
	"net/url"
	"reflect"
	"testing"
	"unicode/utf8"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestQueryLiteralListRoundTrip(t *testing.T) {
	for _, input := range [][]string{{"a,b"}, {" a "}, {""}, {"[literal]"}, {`"quoted"`}, {"a,b", "", " c "}} {
		query := newTestQuery(t, "tags", input)
		encoded := (&Url{Path: "/search", Query: query}).Format()
		parsed, err := ParseUrl(encoded)
		require.NoError(t, err)
		var output []string
		require.NoError(t, parsed.Query.Unmarshal("tags", &output))
		require.Equal(t, input, output)
	}
	query := newTestQuery(t, "codes", []*ErrorCode{NotFound, InvalidArgument})
	var codes []*ErrorCode
	require.NoError(t, query.Unmarshal("codes", &codes))
	require.Equal(t, int32(404), codes[0].Code)
	require.Equal(t, int32(3), codes[1].Code)
	var indirect *[]string
	require.NoError(t, newTestQuery(t, "tags", []string{"a,b"}).Unmarshal("tags", &indirect))
	require.Equal(t, []string{"a,b"}, *indirect)
}

func TestParamFailurePreservesDestination(t *testing.T) {
	value := 9
	require.Error(t, UnmarshalParam("01", &value))
	require.Equal(t, 9, value)
	values := []int{8, 9}
	require.Error(t, UnmarshalParam(`[1,"bad"]`, &values))
	require.Equal(t, []int{8, 9}, values)
	query := NewUrlQueryFrom(url.Values{"values": {"1", "bad"}})
	require.Error(t, query.Unmarshal("values", &values))
	require.Equal(t, []int{8, 9}, values)
	array := [2]string{"old", "value"}
	require.Error(t, newTestQuery(t, "values", []string{"one"}).Unmarshal("values", &array))
	require.Equal(t, [2]string{"old", "value"}, array)
	indirect := &values
	require.Error(t, query.Unmarshal("values", &indirect))
	require.Same(t, &values, indirect)
	require.Equal(t, []int{8, 9}, *indirect)
}

func TestJSONEncodingRejectsValueCycles(t *testing.T) {
	object := NewObject()
	object.SetValue("self", NewObjectValue(object))
	values := &Values{}
	value := NewValuesValue(values)
	values.Vals = []*Value{value}
	for _, input := range []proto.Message{object, NewObjectValue(object), values, value} {
		_, err := json.Marshal(input)
		require.ErrorContains(t, err, "cyclic")
		_, err = jsoniter.Marshal(input)
		require.ErrorContains(t, err, "cyclic")
		var output bytes.Buffer
		require.ErrorContains(t, jsoniter.NewEncoder(&output).Encode(input), "cyclic")
	}

	// Repeated references are valid, and stream state belongs to the caller.
	shared := NewObjectValue(NewObject().SetString("key", "value"))
	stream := jsoniter.NewStream(jsoniter.ConfigCompatibleWithStandardLibrary, nil, 256)
	attachment := &struct{ name string }{"caller"}
	stream.Attachment = attachment
	stream.WriteVal(value)
	require.ErrorContains(t, stream.Error, "cyclic")
	require.Same(t, attachment, stream.Attachment)
	stream.Reset(nil)
	stream.Error = nil
	stream.WriteVal(NewArrayValue(shared, shared))
	require.NoError(t, stream.Error)
	require.JSONEq(t, `[{"key":"value"},{"key":"value"}]`, string(stream.Buffer()))
	require.Same(t, attachment, stream.Attachment)
}

func TestCompactJSONMatchesStandardJSON(t *testing.T) {
	u, err := ParseUrl("https://example.com/path?tag=a%2Cb")
	require.NoError(t, err)
	for _, original := range []proto.Message{
		&Timestamp{Seconds: 1, Nanoseconds: 123},
		&Duration{Seconds: -2, Nanoseconds: -123}, u,
		NewStringValue("NaN"), NewBytesValue([]byte{0, 1, 255}),
		NewInt64Value(math.MinInt64), NewUint64Value(math.MaxUint64),
		NewObject().SetValue("nested", NewArrayValue(NewStringValue("str.test"), NewFloat64Value(1))),
	} {
		t.Run(string(original.ProtoReflect().Descriptor().Name()), func(t *testing.T) {
			standard, err := json.Marshal(original)
			require.NoError(t, err)
			compact, err := jsoniter.Marshal(original)
			require.NoError(t, err)
			require.JSONEq(t, string(compact), string(standard))
			decoded := original.ProtoReflect().New().Interface()
			require.NoError(t, json.Unmarshal(standard, decoded))
			require.True(t, proto.Equal(original, decoded))
			before := proto.Clone(decoded)
			require.Error(t, json.Unmarshal(append(standard, []byte(` true`)...), decoded))
			require.True(t, proto.Equal(before, decoded))
		})
	}
}

func TestInvalidValueBoundariesAreRejected(t *testing.T) {
	for _, value := range []proto.Message{
		&Timestamp{Seconds: 253402300800}, &Timestamp{Seconds: -62135596801},
		&Timestamp{Nanoseconds: -1}, &Timestamp{Nanoseconds: 1000000000},
		&Duration{Seconds: 1, Nanoseconds: -1}, &Duration{Nanoseconds: 1000000000},
		NewNegativeValue(math.MaxUint64),
		NewObject().SetValue("invalid", NewNegativeValue(math.MaxUint64)),
	} {
		_, err := jsoniter.Marshal(value)
		require.Error(t, err, "%T", value)
		_, err = json.Marshal(value)
		require.Error(t, err, "%T", value)
	}
	ts := &Timestamp{Seconds: 7}
	for _, input := range []string{`253402300800`, `-62135596801`} {
		require.Error(t, jsoniter.UnmarshalFromString(input, ts))
		require.Equal(t, int64(7), ts.Seconds)
	}
	for _, input := range []string{`18446744073709551616`, `-9223372036854775809`} {
		value := NewIntValue(7)
		require.Error(t, jsoniter.UnmarshalFromString(input, value))
		require.Equal(t, int64(7), value.GetInt64())
	}
}

func FuzzQueryStringRoundTrip(f *testing.F) {
	for _, value := range []string{"", "a,b", " a ", "[literal]", `"quoted"`, "中文+/?"} {
		f.Add(value)
	}
	f.Fuzz(func(t *testing.T, input string) {
		query, err := NewUrlQuery("value", []string{input})
		require.NoError(t, err)
		parsed, err := ParseUrl((&Url{Path: "/", Query: query}).Format())
		require.NoError(t, err)
		var output []string
		require.NoError(t, parsed.Query.Unmarshal("value", &output))
		require.True(t, reflect.DeepEqual([]string{input}, output))
	})
}

func FuzzValueJSONStringRoundTrip(f *testing.F) {
	for _, value := range []string{"", "NaN", "b64.test", "str.test", "中文", "\x00"} {
		f.Add(value)
	}
	f.Fuzz(func(t *testing.T, input string) {
		if !utf8.ValidString(input) {
			return
		}
		encoded, err := json.Marshal(NewStringValue(input))
		require.NoError(t, err)
		var output Value
		require.NoError(t, json.Unmarshal(encoded, &output))
		require.Equal(t, ValueKind_VALUE_KIND_STRING, output.GetKind())
		require.Equal(t, input, output.GetString())
	})
}
