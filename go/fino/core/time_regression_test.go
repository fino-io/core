package core

import (
	"math"
	"strconv"
	"testing"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/require"
)

func TestDurationBoundaryConversions(t *testing.T) {
	for _, seconds := range []float64{1.9999999996, -1.9999999996} {
		duration := newTestDuration(t, seconds)
		require.Equal(t, int32(0), duration.Nanoseconds)
		require.Equal(t, math.Round(seconds), duration.ToSeconds())
	}
	zero, err := (*Duration)(nil).ToDuration()
	require.NoError(t, err)
	require.Zero(t, zero)
	large := &Duration{Seconds: 10000000000}
	require.Equal(t, float64(10000000000), large.ToSeconds())
	require.Equal(t, float64(10000000000)/3600, large.ToHours())
}

func TestTimestampParseKeepsTimezoneInput(t *testing.T) {
	for _, input := range []string{
		"2019-10-13 11:52:25 UTC", "2019-10-13 11:52:25 +0000",
		"2019-10-13T19:52:25+08:00",
	} {
		t.Run(input, func(t *testing.T) {
			value, err := ParseTimestamp(input)
			require.NoError(t, err)
			require.Equal(t, Timestamp1, value.Seconds)
		})
	}
	value := &Timestamp{Seconds: 1, Nanoseconds: 2}
	require.Error(t, value.Parse("invalid time"))
	require.Equal(t, int64(1), value.Seconds)
	require.Equal(t, int32(2), value.Nanoseconds)
	require.Error(t, (*Timestamp)(nil).Parse(TimestampString1))
}

func TestTimeSQLScan(t *testing.T) {
	duration := newTestDuration(t, 2.5)
	for _, input := range []any{float64(1.5), "1.5", []byte("1.5")} {
		require.NoError(t, duration.Scan(input))
		require.Equal(t, 1.5, duration.ToSeconds())
	}
	require.NoError(t, duration.Scan(int64(math.MaxInt64)))
	require.Equal(t, int64(math.MaxInt64), duration.Seconds)
	require.NoError(t, duration.Scan(nil))
	require.Zero(t, duration.ToSeconds())
	require.Error(t, duration.Scan(math.Inf(1)))
	require.Error(t, duration.Scan(true))
	require.Error(t, (*Duration)(nil).Scan(1.0))

	now := time.Date(2500, 1, 2, 3, 4, 5, 123456789, time.UTC)
	timestamp := FromTime(now)
	require.Equal(t, int32(now.Nanosecond()), timestamp.Nanoseconds)
	for _, input := range []any{now, now.Format(time.RFC3339Nano), []byte(now.Format(time.RFC3339Nano))} {
		require.NoError(t, timestamp.Scan(input))
		require.True(t, now.Equal(timestamp.ToTime()))
	}
	require.NoError(t, timestamp.Scan(nil))
	require.Zero(t, timestamp.Seconds)
	require.Zero(t, timestamp.Nanoseconds)
	require.Error(t, timestamp.Scan(true))
	require.Error(t, (*Timestamp)(nil).Scan(now))
}

func TestTimeJSONDecodeValidation(t *testing.T) {
	timestamp := &Timestamp{Seconds: 1, Nanoseconds: 123}
	require.NoError(t, jsoniter.UnmarshalFromString(`2`, timestamp))
	require.Equal(t, int64(2), timestamp.Seconds)
	require.Zero(t, timestamp.Nanoseconds)
	for _, input := range []string{`true`, `{}`, `1.5`, `9223372036854775808`} {
		require.Error(t, jsoniter.UnmarshalFromString(input, timestamp))
	}
	require.NoError(t, jsoniter.UnmarshalFromString(`null`, timestamp))
	require.Zero(t, timestamp.Seconds)

	duration := newTestDuration(t, 1)
	for _, input := range []string{`true`, `{}`, `1e999`} {
		require.Error(t, jsoniter.UnmarshalFromString(input, duration))
	}
	require.NoError(t, jsoniter.UnmarshalFromString(`null`, duration))
	require.Zero(t, duration.ToSeconds())
}

func TestDurationRejectsOverflow(t *testing.T) {
	for _, input := range []string{`1e30`, `-1e30`, `9223372036854775808`, `-9223372036854775809`} {
		t.Run(input, func(t *testing.T) {
			duration := newTestDuration(t, 2.5)
			require.Error(t, jsoniter.UnmarshalFromString(input, duration))
			require.Equal(t, 2.5, duration.ToSeconds())
			require.Error(t, duration.Scan(input))
			require.Equal(t, 2.5, duration.ToSeconds())
		})
	}
	for _, input := range []float64{math.NaN(), math.Inf(1), math.Inf(-1), 1e30, -1e30} {
		duration := newTestDuration(t, 2.5)
		require.Error(t, duration.Scan(input))
		require.Equal(t, 2.5, duration.ToSeconds())
	}
}

func TestDurationIntegerPrecision(t *testing.T) {
	for _, input := range []string{`9223372036854775807`, `-9223372036854775808`} {
		t.Run(input, func(t *testing.T) {
			want, err := strconv.ParseInt(input, 10, 64)
			require.NoError(t, err)
			duration := newTestDuration(t, 2.5)
			require.NoError(t, jsoniter.UnmarshalFromString(input, duration))
			require.Equal(t, want, duration.Seconds)
			require.Zero(t, duration.Nanoseconds)
			require.NoError(t, duration.Scan(input))
			require.Equal(t, want, duration.Seconds)
		})
	}
}
