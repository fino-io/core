package core

import (
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

type structCodecFields struct {
	URL      *Url `json:"url"`
	URLValue Url  `json:"url_value"`
}

type stringCodecFields struct {
	URL *Url `json:"url"`
}

func init() {
	RegisterJSONFieldDecoder("core.structCodecFields", "URL", &UrlStructCodec{isFieldPointer: true})
	RegisterJSONFieldDecoder("core.structCodecFields", "URLValue", &UrlStructCodec{})
	RegisterJSONFieldEncoder("core.stringCodecFields", "URL", &UrlStringCodec{isFieldPointer: true})
	RegisterJSONFieldDecoder("core.stringCodecFields", "URL", &UrlStringCodec{isFieldPointer: true})
}

func TestStringFieldCodecsRoundTrip(t *testing.T) {
	var fields stringCodecFields
	data, err := jsoniter.MarshalToString(&fields)
	require.NoError(t, err)
	require.JSONEq(t, `{"url":null}`, data)
	fields.URL, err = ParseUrl("https://example.com/path")
	require.NoError(t, err)
	data, err = jsoniter.MarshalToString(&fields)
	require.NoError(t, err)
	var decoded stringCodecFields
	require.NoError(t, jsoniter.UnmarshalFromString(data, &decoded))
	require.True(t, proto.Equal(fields.URL, decoded.URL))
	require.NoError(t, jsoniter.UnmarshalFromString(`{"url":null}`, &decoded))
	require.Nil(t, decoded.URL)
}

func TestStructCodecsDecode(t *testing.T) {
	var fields structCodecFields
	require.NoError(t, jsoniter.UnmarshalFromString(`{"url":{"scheme":"https","path":"/new"},"url_value":{"path":"/value"}}`, &fields))
	require.Equal(t, "https", fields.URL.GetScheme())
	require.Equal(t, "/new", fields.URL.GetPath())
	require.Equal(t, "/value", fields.URLValue.GetPath())
	// Struct decoding retains ordinary JSON partial-update behavior.
	require.NoError(t, jsoniter.UnmarshalFromString(`{"url":{"path":"/updated"}}`, &fields))
	require.Equal(t, "https", fields.URL.GetScheme())
	require.Equal(t, "/updated", fields.URL.GetPath())
	require.NoError(t, jsoniter.UnmarshalFromString(`{"url":null,"url_value":null}`, &fields))
	require.Nil(t, fields.URL)
	require.Empty(t, fields.URLValue.GetPath())
}

func TestStructCodecsRejectInvalidInput(t *testing.T) {
	for _, input := range []string{
		`{"url":1}`, `{"url_value":"invalid"}`,
		`{"url":{"scheme":1}}`, `{"url_value":{"authority":{"port":1}}}`,
	} {
		t.Run(input, func(t *testing.T) {
			var fields structCodecFields
			require.Error(t, jsoniter.UnmarshalFromString(input, &fields))
			require.Nil(t, fields.URL)
		})
	}
}

func TestURLStringCodecFailureKeepsOriginal(t *testing.T) {
	for _, input := range []string{`123`, `true`, `"unterminated`, `"https://host/?bad=%zz"`} {
		t.Run(input, func(t *testing.T) {
			url, err := ParseUrl("https://example.com/original")
			require.NoError(t, err)
			before := proto.Clone(url)
			require.Error(t, jsoniter.UnmarshalFromString(input, url))
			require.True(t, proto.Equal(before, url))
		})
	}
	url, err := ParseUrl("https://example.com/original")
	require.NoError(t, err)
	require.NoError(t, jsoniter.UnmarshalFromString(`null`, url))
	require.True(t, proto.Equal(&Url{}, url))
}

func TestScalarCodecsRejectMalformedJSON(t *testing.T) {
	for _, input := range []string{`01`, `1.`, `1.e2`, `+1`, `nul`, `"unterminated`} {
		t.Run(input, func(t *testing.T) {
			for _, target := range []proto.Message{
				NewIntValue(7), newTestDuration(t, 7.5), &Timestamp{Seconds: 7, Nanoseconds: 123},
			} {
				before := proto.Clone(target)
				require.Error(t, jsoniter.UnmarshalFromString(input, target), "%T", target)
				require.True(t, proto.Equal(before, target), "%T", target)
			}
		})
	}
}
