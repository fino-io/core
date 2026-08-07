package core

import (
	"testing"

	jsoniter "github.com/json-iterator/go"
)

func TestValues(t *testing.T) {
	vals := &Values{Vals: []*Value{
		NewBoolValue(true), NewStringValue("foo"), NewInt8Value(18),
	}}

	toString, err := jsoniter.MarshalToString(vals)
	if err != nil {
		t.Fatal(err)
	}

	// fmt.Printf("%s\n", toString)
	// fmt.Printf("%+v\n", toString)
	// fmt.Printf("%#v\n", toString)

	vals2 := &Values{Vals: []*Value{}}
	err = jsoniter.UnmarshalFromString(toString, vals2)
	if err != nil {
		t.Fatal(err)
	}
	// fmt.Printf("%v\n", vals2)
	// fmt.Printf("%+v\n", vals2)
	// fmt.Printf("%#v\n", vals2)
}
