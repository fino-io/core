package core

import (
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestValues(t *testing.T) {
	vals := &Values{Vals: []*Value{
		NewBoolValue(true), NewStringValue("foo"), NewInt8Value(18),
	}}

	toString, err := jsoniter.MarshalToString(vals)
	require.NoError(t, err)
	require.JSONEq(t, `[true,"foo",18]`, toString)

	var decoded Values
	require.NoError(t, jsoniter.UnmarshalFromString(toString, &decoded))
	require.True(t, proto.Equal(vals, &decoded))
}
