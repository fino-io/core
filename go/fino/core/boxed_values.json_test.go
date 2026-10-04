package core

import (
	"testing"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/require"
)

type codecIntValues struct {
	Vals []int
}

func init() {
	RegisterJSONValuesCodec[int, codecIntValues]("core.codecIntValues", func(x *codecIntValues) *[]int { return &x.Vals })
}

func TestRegisteredValuesCodec(t *testing.T) {
	values := &codecIntValues{Vals: []int{9}}
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

func TestRegisteredValuesCodecInvalidInput(t *testing.T) {
	for _, input := range []string{`true`, `{}`, `1`, `[1,"bad"]`, `[1,`, `nul`} {
		t.Run(input, func(t *testing.T) {
			values := &codecIntValues{Vals: []int{9}}
			require.Error(t, jsoniter.UnmarshalFromString(input, values))
			require.Equal(t, []int{9}, values.Vals)
		})
	}
}

func TestValuesCodecNilDestination(t *testing.T) {
	codec := &ValsCodec[int, codecIntValues]{GetVals: func(*codecIntValues) *[]int { return nil }}
	iter := jsoniter.ParseString(jsoniter.ConfigFastest, `[1]`)
	var values codecIntValues
	require.NotPanics(t, func() { codec.Decode(unsafe.Pointer(&values), iter) })
	require.Error(t, iter.Error)
}

func TestValuesCodecNilMessage(t *testing.T) {
	codec := &ValsCodec[int, codecIntValues]{GetVals: func(x *codecIntValues) *[]int { return &x.Vals }}
	require.True(t, codec.IsEmpty(nil))
	stream := jsoniter.NewStream(jsoniter.ConfigFastest, nil, 16)
	codec.Encode(nil, stream)
	require.NoError(t, stream.Error)
	require.Equal(t, `null`, string(stream.Buffer()))
}
