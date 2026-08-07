package core

import (
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
)

func TestObjectCodec_Decode(t *testing.T) {
	val := `{"key":"value"}`
	obj := &Object{}
	err := jsoniter.ConfigFastest.UnmarshalFromString(val, obj)
	assert.NoError(t, err)
	assert.Equal(t, "value", obj.GetString("key"))
}

type Foo struct {
	Name   string  `json:"name"`
	Object *Object `json:"object,omitempty"`
}

func TestObjectCodec_Decode_2(t *testing.T) {
	foo := &Foo{}
	err := jsoniter.ConfigFastest.UnmarshalFromString(`{"name":"foo","object":{"f1":"v1","f2":100}}`, foo)
	assert.NoError(t, err)
	assert.Equal(t, `v1`, foo.Object.GetString("f1"))
}

func TestObjectCodec_Encode(t *testing.T) {
	str, err := jsoniter.ConfigFastest.MarshalToString(&Foo{Name: "foo"})
	assert.NoError(t, err)
	assert.Equal(t, `{"name":"foo"}`, str)
}

func TestObjectCodec_IsEmpty(t *testing.T) {
	Val := "{}"
	object := &Object{}
	err := jsoniter.ConfigFastest.UnmarshalFromString(Val, object)
	assert.NoError(t, err)
}
