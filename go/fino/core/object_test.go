package core

import (
	"math"
	"testing"
	"time"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

type objectFromEncoded struct {
	Value string
}

type objectFromEncoder struct{}

func (*objectFromEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	stream.WriteVal(map[string]string{"codec_key": (*objectFromEncoded)(ptr).Value})
}

func (*objectFromEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return (*objectFromEncoded)(ptr).Value == ""
}

func init() {
	RegisterJSONTypeEncoder("core.objectFromEncoded", &objectFromEncoder{})
}

func TestObjectFromJSONSemantics(t *testing.T) {
	type Embedded struct {
		Label string `json:"label"`
	}
	type inputStruct struct {
		Embedded
		Name    string
		Tagged  string `json:"tagged_name"`
		Count   int    `json:"count"`
		Active  bool   `json:"active"`
		Omitted string `json:"omitted,omitempty"`
		Ignored string `json:"-"`
		hidden  string
	}
	input := inputStruct{Embedded: Embedded{Label: "embedded"}, Name: "alice", Tagged: "bob", Ignored: "ignored", hidden: "hidden"}
	fieldsJSON := `{"label":"embedded","Name":"alice","tagged_name":"bob","count":0,"active":false}`
	tests := []struct {
		name  string
		input any
		want  string
	}{
		{name: "struct", input: input, want: fieldsJSON},
		{name: "struct pointer", input: &input, want: fieldsJSON},
		{name: "typed map", input: map[string]int{"Count": 0}, want: `{"Count":0}`},
		{name: "map", input: map[string]any{"Name": "alice"}, want: `{"Name":"alice"}`},
		{name: "value map", input: map[string]*Value{"Name": NewStringValue("alice")}, want: `{"Name":"alice"}`},
		{name: "object", input: NewObject().SetString("Name", "alice"), want: `{"Name":"alice"}`},
		{name: "object value", input: NewObjectValue(NewObject().SetString("Name", "alice")), want: `{"Name":"alice"}`},
		{name: "integer precision", input: struct {
			ID uint64 `json:"id"`
		}{math.MaxUint64}, want: `{"id":18446744073709551615}`},
		{name: "type codec", input: &objectFromEncoded{Value: "custom"}, want: `{"codec_key":"custom"}`},
		{name: "error metadata", input: NewError(nil, "missing"), want: `{"message":"missing"}`},
		{name: "empty object", input: struct{}{}, want: `{}`},
		{name: "empty map", input: map[string]int{}, want: `{}`},
		{name: "nil", input: nil, want: `null`},
		{name: "nil object", input: (*Object)(nil), want: `null`},
		{name: "nil struct", input: (*inputStruct)(nil), want: `null`},
		{name: "nil map", input: map[string]any(nil), want: `null`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, entry := range []struct {
				name    string
				convert func(*Object, any) error
			}{
				{name: "From", convert: (*Object).From},
			} {
				t.Run(entry.name, func(t *testing.T) {
					object := NewObject().SetString("old", "stale")
					require.NoError(t, entry.convert(object, tt.input))
					data, err := jsoniter.MarshalToString(object)
					require.NoError(t, err)
					require.JSONEq(t, tt.want, data)
				})
			}
			t.Run("NewObjectFrom", func(t *testing.T) {
				object, err := NewObjectFrom(tt.input)
				require.NoError(t, err)
				data, err := jsoniter.MarshalToString(object)
				require.NoError(t, err)
				require.JSONEq(t, tt.want, data)
			})
		})
	}
}

func TestObjectFromFailureKeepsOriginal(t *testing.T) {
	tests := []struct {
		name  string
		input any
	}{
		{name: "scalar", input: 1},
		{name: "string", input: "value"},
		{name: "array", input: []int{1, 2}},
		{name: "scalar codec", input: FromTime(time.Unix(1, 0))},
		{name: "unsupported field", input: map[string]any{"channel": make(chan int)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, convert := range []func(*Object, any) error{(*Object).From} {
				object := NewObject().SetObject("nested", NewObject().SetString("name", "original"))
				before := object.Clone()
				require.Error(t, convert(object, tt.input))
				require.True(t, proto.Equal(before, object))
				empty := &Object{}
				require.Error(t, convert(empty, tt.input))
				require.Nil(t, empty.Vals)
			}
			object, err := NewObjectFrom(tt.input)
			require.Error(t, err)
			require.Nil(t, object)
		})
	}
	for _, convert := range []func(*Object, any) error{(*Object).From} {
		var object *Object
		require.NoError(t, convert(object, make(chan int)))
	}
}

func TestObjectFromSeparatesNestedValues(t *testing.T) {
	source := NewObject().SetObject("nested", NewObject().SetString("name", "original")).
		SetBytes("data", []byte("bytes")).SetStringArray("items", "original").SetFloat64("number", 1)
	source.ProtoReflect().SetUnknown([]byte{0xa0, 0x06, 0x01})
	for _, input := range []any{source, source.Vals} {
		object, err := NewObjectFrom(input)
		require.NoError(t, err)
		object.GetObject("nested").SetString("name", "changed")
		object.GetBytes("data")[0] = 'x'
		object.GetValueArray("items")[0].Val = NewStringValue("changed").Val
		require.Equal(t, "original", source.GetObject("nested").GetString("name"))
		require.Equal(t, []byte("bytes"), source.GetBytes("data"))
		require.Equal(t, []string{"original"}, source.GetStringArray("items"))
		require.Equal(t, ValueKind_VALUE_KIND_NUMBER, object.GetValue("number").GetKind())
		require.Empty(t, object.ProtoReflect().GetUnknown())
	}
	clone := source.Clone()
	require.True(t, proto.Equal(source, clone))
	require.Equal(t, ValueKind_VALUE_KIND_NUMBER, clone.GetValue("number").GetKind())
	merged := NewObject().SetString("keep", "present").Merge(source)
	require.Equal(t, "present", merged.GetString("keep"))
	require.Same(t, source.GetValue("nested"), merged.GetValue("nested"))
}

func TestObjectFromIntegerPrecision(t *testing.T) {
	input := struct {
		Max uint64 `json:"max"`
		Min int64  `json:"min"`
	}{Max: math.MaxUint64, Min: math.MinInt64}
	object, err := NewObjectFrom(input)
	require.NoError(t, err)
	require.Equal(t, input.Max, object.GetUint64("max"))
	require.Equal(t, input.Min, object.GetInt64("min"))
}

func TestObjectMergeAndClone(t *testing.T) {
	object := &Object{}
	object.Merge(nil).Merge(NewObject().SetObject("nested", NewObject().SetString("name", "original")))
	clone := object.Clone()
	clone.GetObject("nested").SetString("name", "changed")
	require.Equal(t, "original", object.GetObject("nested").GetString("name"))
	require.Equal(t, "changed", clone.GetObject("nested").GetString("name"))
}

func TestObjectToJSONRoundTrip(t *testing.T) {
	type payload struct {
		Name  string
		Count int `json:"count"`
	}
	input := payload{Name: "alice"}
	object, err := NewObjectFrom(input)
	require.NoError(t, err)
	var output payload
	require.NoError(t, object.To(&output))
	require.Equal(t, input, output)
	require.Error(t, object.To((*int)(nil)))
	require.Error(t, object.To(1))
}

func TestNewObject(t *testing.T) {
	v := map[string]any{"name": "apple", "age": 20}
	got, err := NewObjectFromMap(v)
	assert.NoError(t, err)
	m := got.AsMap()
	assert.Equal(t, "apple", m["name"])
	assert.EqualValues(t, 20, m["age"])
}

func TestObject_AsMap(t *testing.T) {
	v := map[string]any{
		"map": map[string]any{
			"inmap": struct{ a int }{a: 11},
		},
	}
	got, err := NewObjectFromMap(v)
	assert.NoError(t, err)
	bytes, err := jsoniter.Marshal(got)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"map":{"inmap":{}}}`, string(bytes))
}
