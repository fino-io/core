package core

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToStringScalars(t *testing.T) {
	type count int32
	for _, tt := range []struct {
		value any
		want  string
	}{
		{int(12), "12"}, {int8(-12), "-12"}, {int16(12), "12"},
		{int32(12), "12"}, {int64(math.MinInt64), "-9223372036854775808"},
		{uint(12), "12"}, {uint8(12), "12"}, {uint16(12), "12"},
		{uint32(12), "12"}, {uint64(math.MaxUint64), "18446744073709551615"},
		{float32(1.25), "1.25"}, {float64(1.23456789012345), "1.23456789012345"},
		{count(12), "12"}, {true, "true"}, {[]byte{0xab}, "ab"}, {nil, ""},
	} {
		require.Equal(t, tt.want, ToString(tt.value), "%T", tt.value)
	}
}

func TestStringValuesMatching(t *testing.T) {
	values := NewStringValues("alpha", "beta", "alpha")
	require.True(t, values.Matched(`^al`))
	require.Equal(t, []string{"alpha", "alpha"}, values.Matches(`^al`).Vals)
	require.False(t, values.Matched("["))
	require.Empty(t, values.Matches("[").Vals)
	require.Equal(t, []string{"alpha", "beta"}, values.Unique().Vals)
	require.Equal(t, true, values.Contains("beta"))
	var empty *StringValues
	require.False(t, empty.Matched("."))
	require.Empty(t, empty.Matches(".").Vals)
	require.Empty(t, empty.Unique().Vals)
	require.Equal(t, false, empty.Contains("beta"))
}

func TestQuoteStringProducesJSON(t *testing.T) {
	for _, input := range []string{"a\"b", "a\\b", "a\nb", "\x01", `"`, `"unfinished`, `unfinished"`} {
		var decoded string
		require.NoError(t, json.Unmarshal([]byte(QuoteString(input)), &decoded))
		require.Equal(t, input, decoded)
	}
	require.Equal(t, `"already quoted"`, QuoteString(`"already quoted"`))
}

func TestOptionsZeroValue(t *testing.T) {
	var options Options
	options = options.SetValue("b", 2).SetValues("a", 1, 3, "ignored", "incomplete")
	require.Equal(t, []any{"a", 1, "b", 2}, options.KeyValues())
	var other Options
	other = other.Merge(options)
	require.Equal(t, options, other)
}

func TestCreatDirRejectsExistingFile(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, CreatDir(filepath.Join(dir, "nested", "dir")))
	file := filepath.Join(dir, "file")
	require.NoError(t, os.WriteFile(file, nil, 0600))
	require.Error(t, CreatDir(file))
}
