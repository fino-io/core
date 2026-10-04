package core

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestCalendarParsingAndValidation(t *testing.T) {
	date, err := ParseDate("2024-02-29")
	require.NoError(t, err)
	require.Equal(t, "2024-02-29", date.Format())
	before := proto.Clone(date)
	for _, value := range []string{"2023-02-29", "2024-13-01", "0000-01-01"} {
		require.Error(t, date.Parse(value))
		require.True(t, proto.Equal(before, date))
	}
	location := time.FixedZone("", 8*3600)
	timestamp, err := date.ToTime(location)
	require.NoError(t, err)
	require.Equal(t, location, timestamp.Location())
	require.Equal(t, "2024-02-29", (&Date{}).FromTime(timestamp).Format())
	tod, err := ParseTimeOfDay("23:59:59.123456789")
	require.NoError(t, err)
	require.Equal(t, "23:59:59.123456789", tod.Format())
	duration, err := tod.ToDuration()
	require.NoError(t, err)
	require.Equal(t, 24*time.Hour-time.Second+123456789, duration)
	before = proto.Clone(tod)
	require.Error(t, tod.Parse("24:00:00"))
	require.True(t, proto.Equal(before, tod))
	for _, value := range []*TimeOfDay{{Hours: -1}, {Minutes: 60}, {Seconds: 60}, {Nanoseconds: 1000000000}} {
		require.Error(t, value.CheckValid())
	}
	require.Error(t, (*Date)(nil).CheckValid())
	require.Error(t, (*TimeOfDay)(nil).CheckValid())
}

func TestDateTimeTimezoneRules(t *testing.T) {
	for _, input := range []string{"2024-02-29T12:34:56.123456789+08:00", "2024-01-01T00:00:00Z", "0001-01-01T00:00:00Z", "9999-12-31T23:59:59Z"} {
		value, err := ParseDateTime(input)
		require.NoError(t, err)
		require.Equal(t, input, value.Format())
		converted, err := value.ToTime()
		require.NoError(t, err)
		require.Equal(t, input, converted.Format(time.RFC3339Nano))
	}
	for _, month := range []int32{1, 7} {
		value := &DateTime{Year: 2024, Month: month, Day: 1, Hour: 12, TimeZone: &TimeZone{Name: "America/New_York", Offset: 8 * 3600}}
		converted, err := value.ToTime()
		require.NoError(t, err)
		_, offset := converted.Zone()
		want := -5 * 3600
		if month == 7 {
			want = -4 * 3600
		}
		require.Equal(t, want, offset)
	}
	gap := &DateTime{Year: 2024, Month: 3, Day: 10, Hour: 2, Minute: 30, TimeZone: &TimeZone{Name: "America/New_York"}}
	require.Error(t, gap.CheckValid())
	apia, err := time.LoadLocation("Pacific/Apia")
	require.NoError(t, err)
	_, err = (&Date{Year: 2011, Month: 12, Day: 30}).ToTime(apia)
	require.ErrorContains(t, err, "nonexistent local date")
	require.Error(t, (&DateTime{Year: 2011, Month: 12, Day: 30, TimeZone: &TimeZone{Name: "Pacific/Apia"}}).CheckValid())
	for _, zone := range []*TimeZone{{Name: "not/a/timezone"}, {Offset: 86400}, {Offset: -86400}, {Offset: 1}} {
		require.Error(t, zone.CheckValid())
	}
	location, err := (*TimeZone)(nil).Location()
	require.NoError(t, err)
	require.Equal(t, time.UTC, location)
	custom := time.Date(2024, 1, 1, 12, 0, 0, 0, time.FixedZone("business", 8*3600))
	value := (&DateTime{}).FromTime(custom)
	converted, err := value.ToTime()
	require.NoError(t, err)
	require.True(t, custom.Equal(converted))
	before := proto.Clone(value)
	require.Error(t, value.Parse("2023-02-29T00:00:00Z"))
	require.True(t, proto.Equal(before, value))
	newYork, err := time.LoadLocation("America/New_York")
	require.NoError(t, err)
	fold := time.Date(2024, 11, 3, 6, 30, 0, 0, time.UTC).In(newYork)
	for _, original := range []time.Time{fold, time.Date(2024, 1, 1, 12, 0, 0, 0, time.FixedZone("UTC", 8*3600))} {
		converted, err := (&DateTime{}).FromTime(original).ToTime()
		require.NoError(t, err)
		require.True(t, original.Equal(converted))
	}
}

func TestValueGraphValidation(t *testing.T) {
	shared := NewStringValue("shared")
	valid := NewObject().SetValue("a", shared).SetValue("b", shared)
	require.NoError(t, valid.CheckValid())
	cycle := NewObject()
	cycle.SetValue("self", NewObjectValue(cycle))
	require.ErrorContains(t, cycle.CheckValid(), "cyclic")
	for _, input := range []any{cycle, NewNegativeValue(math.MaxUint64), NewObject().SetValue("bad", NewNegativeValue(0)), &Value{Val: (*Value_BoolValue)(nil)}} {
		converted, err := NewValue(input)
		require.Error(t, err)
		require.Nil(t, converted)
	}
	require.NoError(t, NewInt64Value(math.MinInt64).CheckValid())
	require.NoError(t, NewUint64Value(math.MaxUint64).CheckValid())
	require.NoError(t, (*Value)(nil).CheckValid())
}

func TestTimeSQLAndArithmeticValidateBoundaries(t *testing.T) {
	ts := &Timestamp{Seconds: 1, Nanoseconds: 2}
	require.Error(t, ts.Scan(time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC)))
	require.Equal(t, int64(1), ts.Seconds)
	require.Equal(t, int32(2), ts.Nanoseconds)
	_, err := (&Timestamp{Nanoseconds: -1}).Value()
	require.Error(t, err)
	_, err = (&Duration{Seconds: 1, Nanoseconds: -1}).Value()
	require.Error(t, err)
	_, err = (&Timestamp{Seconds: 253402300799, Nanoseconds: 999999999}).Add(&Duration{Nanoseconds: 1})
	require.Error(t, err)
	query := newTestQuery(t, "time", "original")
	for _, invalid := range []any{&Timestamp{Seconds: 253402300800}, []*Duration{{Seconds: 1, Nanoseconds: -1}}} {
		require.Error(t, query.Set("time", invalid))
		require.Equal(t, []string{"original"}, query.Vals["time"].Vals)
	}
}

func FuzzTimestampJSONRoundTrip(f *testing.F) {
	f.Add(int64(0), int32(0))
	f.Add(int64(-62135596800), int32(1))
	f.Add(int64(253402300799), int32(999999999))
	f.Fuzz(func(t *testing.T, seconds int64, nanos int32) {
		original := &Timestamp{Seconds: seconds, Nanoseconds: nanos}
		if original.CheckValid() != nil {
			return
		}
		encoded, err := original.MarshalJSON()
		require.NoError(t, err)
		var decoded Timestamp
		require.NoError(t, decoded.UnmarshalJSON(encoded))
		require.True(t, proto.Equal(original, &decoded))
	})
}
