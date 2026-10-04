package core

import (
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

type structCodecFields struct {
	Code      *ErrorCode `json:"code"`
	CodeValue ErrorCode  `json:"code_value"`
	URL       *Url       `json:"url"`
	URLValue  Url        `json:"url_value"`
}

type stringCodecFields struct {
	Code *ErrorCode `json:"code"`
	URL  *Url       `json:"url"`
}

func init() {
	RegisterJSONFieldDecoder("core.structCodecFields", "Code", &ErrorCodeStructCodec{IsFieldPointer: true})
	RegisterJSONFieldDecoder("core.structCodecFields", "CodeValue", &ErrorCodeStructCodec{})
	RegisterJSONFieldDecoder("core.structCodecFields", "URL", &UrlStructCodec{isFieldPointer: true})
	RegisterJSONFieldDecoder("core.structCodecFields", "URLValue", &UrlStructCodec{})
	RegisterJSONFieldEncoder("core.stringCodecFields", "Code", &ErrorCodeStringCodec{IsFieldPointer: true})
	RegisterJSONFieldDecoder("core.stringCodecFields", "Code", &ErrorCodeStringCodec{IsFieldPointer: true})
	RegisterJSONFieldEncoder("core.stringCodecFields", "URL", &UrlStringCodec{isFieldPointer: true})
	RegisterJSONFieldDecoder("core.stringCodecFields", "URL", &UrlStringCodec{isFieldPointer: true})
}

func TestStringFieldCodecsRoundTrip(t *testing.T) {
	var fields stringCodecFields
	data, err := jsoniter.MarshalToString(&fields)
	require.NoError(t, err)
	require.JSONEq(t, `{"code":null,"url":null}`, data)
	fields.Code = NewErrorCode(404)
	fields.URL, err = ParseUrl("https://example.com/path")
	require.NoError(t, err)
	data, err = jsoniter.MarshalToString(&fields)
	require.NoError(t, err)
	var decoded stringCodecFields
	require.NoError(t, jsoniter.UnmarshalFromString(data, &decoded))
	require.True(t, proto.Equal(fields.Code, decoded.Code))
	require.True(t, proto.Equal(fields.URL, decoded.URL))
	require.NoError(t, jsoniter.UnmarshalFromString(`{"code":null,"url":null}`, &decoded))
	require.Nil(t, decoded.Code)
	require.Nil(t, decoded.URL)
}

func TestStructCodecsDecode(t *testing.T) {
	var fields structCodecFields
	require.NoError(t, jsoniter.UnmarshalFromString(`{"code":{"code":404},"code_value":{"code":400},"url":{"scheme":"https","path":"/new"},"url_value":{"path":"/value"}}`, &fields))
	require.Equal(t, int32(404), fields.Code.GetCode())
	require.Equal(t, int32(400), fields.CodeValue.GetCode())
	require.Equal(t, "https", fields.URL.GetScheme())
	require.Equal(t, "/new", fields.URL.GetPath())
	require.Equal(t, "/value", fields.URLValue.GetPath())
	// Struct decoding retains ordinary JSON partial-update behavior.
	require.NoError(t, jsoniter.UnmarshalFromString(`{"url":{"path":"/updated"}}`, &fields))
	require.Equal(t, "https", fields.URL.GetScheme())
	require.Equal(t, "/updated", fields.URL.GetPath())
	require.NoError(t, jsoniter.UnmarshalFromString(`{"code":null,"code_value":null,"url":null,"url_value":null}`, &fields))
	require.Nil(t, fields.Code)
	require.Nil(t, fields.URL)
	require.Zero(t, fields.CodeValue.GetCode())
	require.Empty(t, fields.URLValue.GetPath())
}

func TestStructCodecsRejectInvalidInput(t *testing.T) {
	for _, input := range []string{
		`{"code":true}`, `{"code_value":[]}`, `{"url":1}`, `{"url_value":"invalid"}`,
		`{"code":{"code":"bad"}}`, `{"code_value":{"httpStatusCode":"bad"}}`,
		`{"url":{"scheme":1}}`, `{"url_value":{"authority":{"port":1}}}`,
	} {
		t.Run(input, func(t *testing.T) {
			var fields structCodecFields
			require.Error(t, jsoniter.UnmarshalFromString(input, &fields))
			require.Nil(t, fields.Code)
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

func TestErrorCodeStringCodecNullAndInvalidInput(t *testing.T) {
	err := NewErrorFrom(404, "missing")
	require.NoError(t, jsoniter.UnmarshalFromString(`{"code":null}`, err))
	require.Nil(t, err.Code)
	for _, input := range []string{`{"code":true}`, `{"code":"bad"}`, `{"code":nul}`} {
		var target Error
		require.Error(t, jsoniter.UnmarshalFromString(input, &target))
		require.Nil(t, target.Code)
	}
}

func TestScalarCodecsRejectMalformedJSON(t *testing.T) {
	for _, input := range []string{`01`, `1.`, `1.e2`, `+1`, `nul`, `"unterminated`} {
		t.Run(input, func(t *testing.T) {
			for _, target := range []proto.Message{
				NewIntValue(7), NewDuration(7.5), &Timestamp{Seconds: 7, Nanoseconds: 123},
			} {
				before := proto.Clone(target)
				require.Error(t, jsoniter.UnmarshalFromString(input, target), "%T", target)
				require.True(t, proto.Equal(before, target), "%T", target)
			}
		})
	}
}
