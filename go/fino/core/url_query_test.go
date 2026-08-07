package core

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"

	jsoniter "github.com/json-iterator/go"
)

func TestNewUrlQueryFrom(t *testing.T) {
	tests := []struct {
		name   string
		values url.Values
		want   *Url_Query
	}{
		{name: "empty", values: url.Values{}, want: &Url_Query{Vals: make(map[string]*StringValues)}},
		{name: "v1", values: url.Values{"key1": {"v1"}}, want: &Url_Query{Vals: map[string]*StringValues{"key1": {Vals: []string{"v1"}}}}},
		{name: "v2", values: url.Values{"key1": {"v1"}, "key2": {"v2"}}, want: &Url_Query{Vals: map[string]*StringValues{"key1": {Vals: []string{"v1"}}, "key2": {Vals: []string{"v2"}}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUrlQueryFrom(tt.values); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUrlQueryFrom() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

type MyToStringStruct struct{}

func (x *MyToStringStruct) ToString() string {
	return "this is my custom struct"
}

type MyFormatterStruct struct{}

func (x *MyFormatterStruct) Format() string {
	return "this is a formatter"
}

func TestQuery_Add(t *testing.T) {
	type args struct {
		k string
		v any
	}
	tests := []struct {
		name string
		Vals map[string]*StringValues
		args args
		want *Url_Query
	}{
		{name: "empty", Vals: map[string]*StringValues{}, args: args{}, want: &Url_Query{Vals: map[string]*StringValues{}}},
		{name: "add-nil", Vals: map[string]*StringValues{"key1": {Vals: []string{"v1"}}}, args: args{}, want: &Url_Query{Vals: map[string]*StringValues{"key1": {Vals: []string{"v1"}}}}},
		{name: "add-empty-value", Vals: map[string]*StringValues{"key1": {Vals: []string{"v1"}}}, args: args{k: "key1"}, want: &Url_Query{Vals: map[string]*StringValues{"key1": {Vals: []string{"v1"}}}}},
		{name: "add-nonempty-value", Vals: map[string]*StringValues{"key1": {Vals: []string{"v1"}}}, args: args{k: "key1", v: "v11"}, want: &Url_Query{Vals: map[string]*StringValues{"key1": {Vals: []string{"v1", "v11"}}}}},
		{name: "add-other-empty-key", Vals: map[string]*StringValues{"key1": {Vals: []string{"v1"}}}, args: args{k: "key2"}, want: &Url_Query{Vals: map[string]*StringValues{"key1": {Vals: []string{"v1"}}, "key2": {}}}},
		{name: "add-other-nonempty-key", Vals: map[string]*StringValues{"key1": {Vals: []string{"v1"}}}, args: args{k: "key2", v: "v2"}, want: &Url_Query{Vals: map[string]*StringValues{"key1": {Vals: []string{"v1"}}, "key2": {Vals: []string{"v2"}}}}},
		{name: "add-key-bool-true", Vals: map[string]*StringValues{}, args: args{k: "key2", v: true}, want: &Url_Query{Vals: map[string]*StringValues{"key2": {Vals: []string{"true"}}}}},
		{name: "add-key-bool-false", Vals: map[string]*StringValues{}, args: args{k: "key2", v: false}, want: &Url_Query{Vals: map[string]*StringValues{"key2": {Vals: []string{"false"}}}}},
		{name: "add-key-int8", Vals: map[string]*StringValues{}, args: args{k: "key2", v: int8(13)}, want: &Url_Query{Vals: map[string]*StringValues{"key2": {Vals: []string{"13"}}}}},
		{name: "add-key-int16", Vals: map[string]*StringValues{}, args: args{k: "key2", v: int16(13)}, want: &Url_Query{Vals: map[string]*StringValues{"key2": {Vals: []string{"13"}}}}},
		{name: "add-key-int32", Vals: map[string]*StringValues{}, args: args{k: "key2", v: int32(13)}, want: &Url_Query{Vals: map[string]*StringValues{"key2": {Vals: []string{"13"}}}}},
		{name: "add-key-int64", Vals: map[string]*StringValues{}, args: args{k: "key2", v: int64(13)}, want: &Url_Query{Vals: map[string]*StringValues{"key2": {Vals: []string{"13"}}}}},
		{name: "add-key-int", Vals: map[string]*StringValues{}, args: args{k: "key2", v: int(13)}, want: &Url_Query{Vals: map[string]*StringValues{"key2": {Vals: []string{"13"}}}}},
		{name: "add-key-uint8", Vals: map[string]*StringValues{}, args: args{k: "key2", v: uint8(13)}, want: &Url_Query{Vals: map[string]*StringValues{"key2": {Vals: []string{"13"}}}}},
		{name: "add-key-uint16", Vals: map[string]*StringValues{}, args: args{k: "key2", v: uint16(13)}, want: &Url_Query{Vals: map[string]*StringValues{"key2": {Vals: []string{"13"}}}}},
		{name: "add-key-uint32", Vals: map[string]*StringValues{}, args: args{k: "key2", v: uint32(13)}, want: &Url_Query{Vals: map[string]*StringValues{"key2": {Vals: []string{"13"}}}}},
		{name: "add-key-uint64", Vals: map[string]*StringValues{}, args: args{k: "key2", v: uint64(13)}, want: &Url_Query{Vals: map[string]*StringValues{"key2": {Vals: []string{"13"}}}}},
		{name: "add-key-uint", Vals: map[string]*StringValues{}, args: args{k: "key2", v: uint(13)}, want: &Url_Query{Vals: map[string]*StringValues{"key2": {Vals: []string{"13"}}}}},
		{name: "add-key-float32", Vals: map[string]*StringValues{}, args: args{k: "key2", v: float32(13.0)}, want: &Url_Query{Vals: map[string]*StringValues{"key2": {Vals: []string{"13"}}}}},
		{name: "add-key-float64", Vals: map[string]*StringValues{}, args: args{k: "key2", v: float64(13.0)}, want: &Url_Query{Vals: map[string]*StringValues{"key2": {Vals: []string{"13"}}}}},
		{name: "add-key-string-array", Vals: map[string]*StringValues{}, args: args{k: "key2", v: []string{"v1", "v2"}}, want: &Url_Query{Vals: map[string]*StringValues{"key2": {Vals: []string{"v1", "v2"}}}}},
		{name: "add-key-ToStringConverter", Vals: map[string]*StringValues{}, args: args{k: "key2", v: &MyToStringStruct{}}, want: &Url_Query{Vals: map[string]*StringValues{"key2": {Vals: []string{"this is my custom struct"}}}}},
		{name: "add-key-Formatter", Vals: map[string]*StringValues{}, args: args{k: "key2", v: &MyFormatterStruct{}}, want: &Url_Query{Vals: map[string]*StringValues{"key2": {Vals: []string{"this is a formatter"}}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &Url_Query{
				Vals: tt.Vals,
			}
			if got := x.Add(tt.args.k, tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuery_Set(t *testing.T) {
	type args struct {
		k string
		v any
	}
	tests := []struct {
		name string
		Vals map[string]*StringValues
		args args
		want *Url_Query
	}{
		{name: "empty", Vals: map[string]*StringValues{}, args: args{}, want: &Url_Query{Vals: map[string]*StringValues{}}},
		{name: "set-nil", Vals: map[string]*StringValues{"key1": {Vals: []string{"v1"}}}, args: args{}, want: &Url_Query{Vals: map[string]*StringValues{"key1": {Vals: []string{"v1"}}}}},
		{name: "set-empty-value", Vals: map[string]*StringValues{"key1": {Vals: []string{"v1"}}}, args: args{k: "key1"}, want: &Url_Query{Vals: map[string]*StringValues{"key1": {Vals: []string{}}}}},
		{name: "set-nonempty-value", Vals: map[string]*StringValues{"key1": {Vals: []string{"v1"}}}, args: args{k: "key1", v: "v11"}, want: &Url_Query{Vals: map[string]*StringValues{"key1": {Vals: []string{"v11"}}}}},
		{name: "set-other-empty-key", Vals: map[string]*StringValues{"key1": {Vals: []string{"v1"}}}, args: args{k: "key2"}, want: &Url_Query{Vals: map[string]*StringValues{"key1": {Vals: []string{"v1"}}, "key2": {Vals: []string{}}}}},
		{name: "set-other-nonempty-key", Vals: map[string]*StringValues{"key1": {Vals: []string{"v1"}}}, args: args{k: "key2", v: "v2"}, want: &Url_Query{Vals: map[string]*StringValues{"key1": {Vals: []string{"v1"}}, "key2": {Vals: []string{"v2"}}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &Url_Query{
				Vals: tt.Vals,
			}
			if got := x.Set(tt.args.k, tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Set() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuery_Unmarshal(t *testing.T) {
	tests := []struct {
		name    string
		query   *Url_Query
		k       string
		typ     reflect.Type
		want    any
		wantErr bool
	}{
		{name: "string-slice-1", query: &Url_Query{Vals: map[string]*StringValues{"foo": {Vals: []string{"bar", "bba"}}}}, k: "foo", typ: reflect.TypeOf([]string{}), want: []string{"bar", "bba"}, wantErr: false},
		{name: "string-slice-2", query: &Url_Query{Vals: map[string]*StringValues{"foo": {Vals: []string{"bar,bba"}}}}, k: "foo", typ: reflect.TypeOf([]string{}), want: []string{"bar", "bba"}, wantErr: false},
		{name: "string-slice-3", query: &Url_Query{Vals: map[string]*StringValues{"foo": {Vals: []string{`"bar", "bba"`}}}}, typ: reflect.TypeOf([]string{}), k: "foo", want: []string{"bar", "bba"}, wantErr: false},
		{name: "string-slice-4", query: &Url_Query{Vals: map[string]*StringValues{"foo": {Vals: []string{`"bar,bba`}}}}, typ: reflect.TypeOf([]string{}), k: "foo", want: []string{"\"bar", "bba"}, wantErr: false},
		{name: "string-slice-5", query: &Url_Query{Vals: map[string]*StringValues{"foo": {Vals: []string{`["bar", "bba"]`}}}}, typ: reflect.TypeOf([]string{}), k: "foo", want: []string{"bar", "bba"}, wantErr: false},
		{name: "int-slice-1", query: &Url_Query{Vals: map[string]*StringValues{"foo": {Vals: []string{"123", "234"}}}}, k: "foo", typ: reflect.TypeOf([]int32{}), want: []int32{123, 234}, wantErr: false},
		{name: "int-slice-2", query: &Url_Query{Vals: map[string]*StringValues{"foo": {Vals: []string{"123,234"}}}}, k: "foo", typ: reflect.TypeOf([]int32{}), want: []int32{123, 234}, wantErr: false},
		// {name: "int32Value-slice-1", query: &Url_Query{Vals: map[string]*StringValues{"foo": {Vals: []string{"123,234"}}}}, k: "foo", typ: reflect.TypeOf([]*Int32Value{}), want: []*Int32Value{{Value: 123}, {Value: 234}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ptr := reflect.New(tt.typ)
			if err := tt.query.Unmarshal(tt.k, ptr.Interface()); (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			vals := ptr.Elem().Interface()
			if !reflect.DeepEqual(vals, tt.want) {
				t.Errorf("Unmarshal() = %v, want %v", vals, tt.want)
			}
		})
	}
}

func TestUnmarshalParam(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		typ     reflect.Type
		want    any
		wantErr bool
	}{
		{name: "string-slice-1", str: "bar,bba", typ: reflect.TypeOf([]string{}), want: []string{"bar", "bba"}, wantErr: false},
		{name: "string-slice-2", str: `"bar", "bba"`, typ: reflect.TypeOf([]string{}), want: []string{"bar", "bba"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ptr := reflect.New(tt.typ)
			if err := UnmarshalParam(tt.str, ptr.Interface()); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalParam() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := ptr.Elem().Interface(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnmarshalParam() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_splitQuoteString(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want []string
	}{
		{name: "empty", s: "", want: []string{""}},
		{name: "1 value", s: `"a"`, want: []string{`"a"`}},
		{name: "1 values-2", s: `"a,b"`, want: []string{`"a,b"`}},
		{name: "2 values", s: `"a","b"`, want: []string{`"a"`, `"b"`}},
		{name: "2 values-2", s: `"a,c","b"`, want: []string{`"a,c"`, `"b"`}},
		{name: "2 values-3", s: `"a,c" , "b"`, want: []string{`"a,c"`, `"b"`}},
		{name: "3 values", s: "a,b,c", want: []string{"a", "b", "c"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitQuotedString(tt.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitQuoteString() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func Test_isStringSlice(t *testing.T) {
	tests := []struct {
		name string
		v    any
		want bool
	}{
		{name: "strings", v: []string{}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isStringSlice(tt.v); got != tt.want {
				t.Errorf("isStringSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_unmarshal(t *testing.T) {
	arrayStr := `["bar", "bba"]`
	var v []string
	err := jsoniter.ConfigFastest.UnmarshalFromString(arrayStr, &v)
	fmt.Println(err)
	fmt.Println(v)
}
