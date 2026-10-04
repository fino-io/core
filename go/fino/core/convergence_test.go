package core

import (
	"errors"
	"fmt"
	"math"
	"testing"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func newTestDuration(t *testing.T, seconds float64) *Duration {
	t.Helper()
	value, err := NewDuration(seconds)
	require.NoError(t, err)
	return value
}

func newTestQuery(t *testing.T, pairs ...any) *Url_Query {
	t.Helper()
	value, err := NewUrlQuery(pairs...)
	require.NoError(t, err)
	return value
}

func TestDurationExactRoundTrip(t *testing.T) {
	for _, original := range []*Duration{
		{Seconds: 10000000000, Nanoseconds: 123456789},
		{Seconds: -10000000000, Nanoseconds: -123456789},
		{Seconds: math.MaxInt64, Nanoseconds: 999999999},
		{Seconds: math.MinInt64, Nanoseconds: -999999999},
		{Nanoseconds: -1}, {},
	} {
		t.Run(original.Format(), func(t *testing.T) {
			data, err := jsoniter.Marshal(original)
			require.NoError(t, err)
			var decoded Duration
			require.NoError(t, jsoniter.Unmarshal(data, &decoded))
			require.True(t, proto.Equal(original, &decoded))
			value, err := original.Value()
			require.NoError(t, err)
			require.NoError(t, decoded.Scan(value))
			require.True(t, proto.Equal(original, &decoded))
		})
	}
	for _, input := range []time.Duration{math.MinInt64, math.MaxInt64} {
		got, err := FromDuration(input).ToDuration()
		require.NoError(t, err)
		require.Equal(t, input, got)
	}
	for _, duration := range []*Duration{{Seconds: 10000000000}, {Seconds: -10000000000}} {
		_, err := duration.ToDuration()
		require.Error(t, err)
		_, err = duration.ToNanoSeconds()
		require.Error(t, err)
	}
	require.Zero(t, (&Duration{Seconds: 1, Nanoseconds: -1000000000}).Compare(&Duration{}))
}

func TestDurationConstructorRejectsInvalidSeconds(t *testing.T) {
	minimum, err := NewDuration(float64(math.MinInt64))
	require.NoError(t, err)
	require.Equal(t, int64(math.MinInt64), minimum.Seconds)
	for _, input := range []float64{1.0 / 1024, -1.0 / 1024} {
		value, err := NewDuration(input)
		require.NoError(t, err)
		require.Equal(t, int32(math.Round(input*float64(time.Second))), value.Nanoseconds)
	}
	for _, input := range []float64{math.NaN(), math.Inf(1), math.Inf(-1), 1e30, -1e30} {
		value, err := NewDuration(input)
		require.Error(t, err)
		require.Nil(t, value)
		original := &Duration{Seconds: 1, Nanoseconds: 2}
		require.Error(t, original.FromSeconds(input))
		require.Equal(t, int64(1), original.Seconds)
		require.Equal(t, int32(2), original.Nanoseconds)
	}
	var nilDuration *Duration
	require.Error(t, nilDuration.FromSeconds(1))
	for _, input := range []string{"0.0000000005", "-0.0000000005", "1.9999999996"} {
		var value Duration
		require.NoError(t, value.Scan(input))
		data, err := value.Value()
		require.NoError(t, err)
		var decoded Duration
		require.NoError(t, decoded.Scan(data))
		require.True(t, proto.Equal(&value, &decoded))
	}
}

func TestTimestampPrecisionAndAddOverflow(t *testing.T) {
	original := FromTime(time.Date(2025, 1, 2, 3, 4, 5, 123456789, time.FixedZone("test", 8*3600)))
	data, err := jsoniter.Marshal(original)
	require.NoError(t, err)
	require.Equal(t, `"2025-01-01T19:04:05.123456789Z"`, string(data))
	var decoded Timestamp
	require.NoError(t, jsoniter.Unmarshal(data, &decoded))
	require.True(t, proto.Equal(original, &decoded))
	_, err = original.Add(&Duration{Seconds: 10000000000})
	require.Error(t, err)
	added, err := original.Add(FromDuration(time.Second))
	require.NoError(t, err)
	require.True(t, added.ToTime().Equal(original.ToTime().Add(time.Second)))
}

func TestErrorMetadataAndDetails(t *testing.T) {
	document, err := ParseUrl("https://example.com/errors/conflict")
	require.NoError(t, err)
	code := &ErrorCode{Code: 600001, Name: "BUSINESS", Domain: "billing", Description: "conflict", HttpStatusCode: 409, Document: document}
	original := NewError(code, "conflict")
	code.Document.Path = "/changed"
	require.Equal(t, "/errors/conflict", original.Code.Document.Path)
	data, err := jsoniter.Marshal(original)
	require.NoError(t, err)
	var decoded Error
	require.NoError(t, jsoniter.Unmarshal(data, &decoded))
	require.True(t, proto.Equal(original, &decoded))
	require.Equal(t, 409, decoded.StatusCode())
	require.Error(t, original.AddDetail(make(chan int)))
	require.Empty(t, original.Details)
	require.NoError(t, original.AddDetail("details"))
	require.Equal(t, "details", original.Details[0].GetString())
	require.NoError(t, original.Code.Parse("404"))
	require.Empty(t, original.Code.Domain)
	require.Nil(t, original.Code.Document)
	require.True(t, proto.Equal(NewErrorCode(404), original.Code))
	before := proto.Clone(original.Code)
	require.Error(t, original.Code.Parse("bad"))
	require.True(t, proto.Equal(before, original.Code))
}

func TestErrorClassificationAndTypedNil(t *testing.T) {
	var nilError *Error
	require.Nil(t, AsError(nilError))
	require.False(t, IsNotFoundError(nilError))
	require.False(t, errors.Is(nilError, &Error{}))
	require.Empty(t, nilError.Error())
	err := fmt.Errorf("lookup: %w", NewNotFoundError("missing"))
	require.True(t, IsNotFoundError(err))
	require.False(t, IsBadRequestError(err))
	require.True(t, errors.Is(err, &Error{Code: NewErrorCode(404)}))
	require.False(t, errors.Is(err, &Error{Code: &ErrorCode{Code: 404, Domain: "other"}}))
}

func TestReservedStringsRemainStrings(t *testing.T) {
	for _, input := range []string{"NaN", "Infinity", "-Infinity", "b64.aGVsbG8=", "b64.!", "str.NaN", "str.str.NaN"} {
		t.Run(input, func(t *testing.T) {
			original := NewStringValue(input)
			data, err := jsoniter.Marshal(original)
			require.NoError(t, err)
			var decoded Value
			require.NoError(t, jsoniter.Unmarshal(data, &decoded))
			require.True(t, proto.Equal(original, &decoded))
			require.Equal(t, input, decoded.AsInterface())
			object, err := NewObjectFrom(struct {
				Text string `json:"text"`
			}{input})
			require.NoError(t, err)
			require.Equal(t, input, object.GetString("text"))
			var output struct {
				Text string `json:"text"`
			}
			require.NoError(t, object.To(&output))
			require.Equal(t, input, output.Text)
		})
	}
	bytes, err := NewValue([]byte("hello"))
	require.NoError(t, err)
	require.Equal(t, ValueKind_VALUE_KIND_BYTES, bytes.GetKind())
	require.Equal(t, []byte("hello"), bytes.AsInterface())
}

func TestMapNullEntriesAndReplacement(t *testing.T) {
	values := &StringsMap{Vals: map[string]*StringValues{"old": {Vals: []string{"stale"}}}}
	require.NoError(t, jsoniter.UnmarshalFromString(`{"x":null}`, values))
	require.Len(t, values.Vals, 1)
	data, err := jsoniter.Marshal(values)
	require.NoError(t, err)
	require.JSONEq(t, `{"x":null}`, string(data))
	require.Error(t, jsoniter.UnmarshalFromString(`{"bad":123}`, values))
	require.Len(t, values.Vals, 1)
}

func TestObjectClonePreservesEmptyCollections(t *testing.T) {
	source := NewObject().SetObject("empty", NewObject()).
		SetValue("array", NewArrayValue([]*Value{}...)).SetObject("null", &Object{}).
		SetValue("nested", NewArrayValue(NewObjectValue(NewObject())))
	clone := source.Clone()
	want, err := jsoniter.Marshal(source)
	require.NoError(t, err)
	got, err := jsoniter.Marshal(clone)
	require.NoError(t, err)
	require.JSONEq(t, string(want), string(got))
	clone.GetObject("empty").SetString("changed", "value")
	require.Empty(t, source.GetObject("empty").Vals)
	for _, input := range []any{map[string]any(nil), map[string]int(nil), []any(nil), []int(nil)} {
		value, err := NewValue(input)
		require.NoError(t, err)
		data, err := jsoniter.Marshal(value)
		require.NoError(t, err)
		require.Equal(t, "null", string(data))
	}
}

func TestObjectKeyCollisions(t *testing.T) {
	object := NewObject().SetString("foo_bar", "A").SetString("fooBar", "B")
	for _, convert := range []func() (*Object, error){object.ToLowerCamelKeys, object.ToSnakeKeys} {
		result, err := convert()
		require.Error(t, err)
		require.Nil(t, result)
	}
	require.Len(t, object.Vals, 2)
}

func TestQueryNamedScalarsAndInvalidValues(t *testing.T) {
	type Count int32
	type Names []string
	query := newTestQuery(t, "count", Count(3), "names", Names{"alice", "bob"})
	require.Equal(t, []string{"3"}, query.Vals["count"].Vals)
	require.Equal(t, []string{"alice", "bob"}, query.Vals["names"].Vals)
	before := proto.Clone(query)
	require.Error(t, query.Set("count", make(chan int)))
	require.Error(t, query.Add("count", []any{4, make(chan int)}))
	require.True(t, proto.Equal(before, query))
	for _, pairs := range [][]any{{"key"}, {1, "value"}, {"key", make(chan int)}} {
		value, err := NewUrlQuery(pairs...)
		require.Error(t, err)
		require.Nil(t, value)
	}
}

func TestEnumsRejectInvalidJSON(t *testing.T) {
	for _, input := range []string{`01`, `+1`, `1.0`, `1e0`, `1e999`, `2147483648`, `true`, `"bad"`, `"unterminated`} {
		month := Month(1)
		require.Error(t, jsoniter.UnmarshalFromString(input, &month), input)
		require.Equal(t, Month(1), month, input)
	}
	month := Month(1)
	require.NoError(t, jsoniter.UnmarshalFromString(`null`, &month))
	require.Zero(t, month)
	require.NoError(t, jsoniter.UnmarshalFromString(`1`, &month))
	data, err := jsoniter.Marshal(month)
	require.NoError(t, err)
	var decoded Month
	require.NoError(t, jsoniter.Unmarshal(data, &decoded))
	require.Equal(t, month, decoded)
	_, err = jsoniter.Marshal(Month(99))
	require.Error(t, err)
}
