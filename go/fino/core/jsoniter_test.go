package core

import (
	"testing"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/require"
)

type jsonIntValues struct {
	Vals []int
}

func init() {
	RegisterJSONValuesCodec("core.jsonIntValues", func(x *jsonIntValues) *[]int { return &x.Vals })
}

func TestRegisterJSONValuesCodec(t *testing.T) {
	values := &jsonIntValues{Vals: []int{9}}
	require.NoError(t, jsoniter.UnmarshalFromString(`[1,2]`, values))
	require.Equal(t, []int{1, 2}, values.Vals)
	data, err := jsoniter.MarshalToString(values)
	require.NoError(t, err)
	require.Equal(t, `[1,2]`, data)
	require.NoError(t, jsoniter.UnmarshalFromString(`null`, values))
	require.Nil(t, values.Vals)
	require.NoError(t, jsoniter.UnmarshalFromString(`[]`, values))
	require.NotNil(t, values.Vals)
	require.Empty(t, values.Vals)
}

func TestRegisterJSONValuesCodecRejectsInvalidInput(t *testing.T) {
	for _, input := range []string{`true`, `{}`, `1`, `[1,"bad"]`, `[1,`, `nul`} {
		t.Run(input, func(t *testing.T) {
			values := &jsonIntValues{Vals: []int{9}}
			require.Error(t, jsoniter.UnmarshalFromString(input, values))
			require.Equal(t, []int{9}, values.Vals)
		})
	}
}

func TestJSONFieldCodecNilDestination(t *testing.T) {
	codec := &jsonFieldCodec[[]int, jsonIntValues]{field: func(*jsonIntValues) *[]int { return nil }}
	iter := jsoniter.ParseString(jsoniter.ConfigDefault, `[1]`)
	var values jsonIntValues
	require.NotPanics(t, func() { codec.Decode(unsafe.Pointer(&values), iter) })
	require.Error(t, iter.Error)
}

func TestJSONFieldCodecNilMessage(t *testing.T) {
	codec := &jsonFieldCodec[[]int, jsonIntValues]{field: func(x *jsonIntValues) *[]int { return &x.Vals }}
	require.True(t, codec.IsEmpty(nil))
	stream := jsoniter.NewStream(jsoniter.ConfigDefault, nil, 16)
	codec.Encode(nil, stream)
	require.NoError(t, stream.Error)
	require.Equal(t, `null`, string(stream.Buffer()))
}
