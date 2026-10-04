package core

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	jsoniter "github.com/json-iterator/go"
)

func TestQueryUnmarshalKeepsRawValues(t *testing.T) {
	raw := []string{`a"b`, `a\b`, "a\nb", `"quoted"`}
	query := newTestQuery(t, "foo", raw)
	for i := 0; i < 2; i++ {
		var values []string
		require.NoError(t, query.Unmarshal("foo", &values))
		require.Equal(t, raw, values)
		require.Equal(t, raw, query.Vals["foo"].Vals)
	}
}

func TestQueryZeroValue(t *testing.T) {
	query := &Url_Query{}
	require.NoError(t, query.Add("a", 1))
	require.NoError(t, query.Set("b", "two"))
	require.Equal(t, []string{"1"}, query.Vals["a"].Vals)
	require.Equal(t, []string{"two"}, query.Vals["b"].Vals)
}

func TestUnmarshalParamInvalidDestination(t *testing.T) {
	for _, value := range []any{nil, (*string)(nil), "not a pointer", 1} {
		require.Error(t, UnmarshalParam("value", value))
		require.Error(t, newTestQuery(t, "key", "value").Unmarshal("key", value))
	}
}

func TestUnmarshalParamEscapedList(t *testing.T) {
	var values []string
	require.NoError(t, UnmarshalParam(`"a,\"b", "c\\d"`, &values))
	require.Equal(t, []string{`a,"b`, `c\d`}, values)
}

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

func TestQueryFromURLValuesReplacesAndCopies(t *testing.T) {
	query := newTestQuery(t, "old", "value")
	source := url.Values{"new": {"value"}, "empty": {}}
	require.Same(t, query, query.FromUrlValues(source))
	require.False(t, query.Has("old"))
	source["new"][0] = "changed"
	require.Equal(t, []string{"value"}, query.Vals["new"].Vals)
	require.NotNil(t, query.Vals["empty"].Vals)
	query.FromUrlValues(nil)
	require.Empty(t, query.Vals)
}

type queryParserValue struct{ text string }

func (x *queryParserValue) Parse(value string) error {
	x.text = value
	return nil
}

type queryStringValue struct{ text string }

func (x *queryStringValue) ToString() string { return x.text }

func TestQueryParserAndValueLists(t *testing.T) {
	var custom queryParserValue
	require.NoError(t, UnmarshalParam("raw text", &custom))
	require.Equal(t, "raw text", custom.text)
	var code ErrorCode
	require.NoError(t, UnmarshalParam("404", &code))
	require.Equal(t, int32(404), code.Code)
	query := newTestQuery(t, "months", []Month{Month_MONTH_JANUARY, Month_MONTH_FEBRUARY})
	var months []Month
	require.NoError(t, query.Unmarshal("months", &months))
	require.Equal(t, []Month{Month_MONTH_JANUARY, Month_MONTH_FEBRUARY}, months)
	var timestamps []Timestamp
	require.NoError(t, UnmarshalParam(TimestampString1+","+TimestampString2, &timestamps))
	require.Equal(t, []Timestamp{{Seconds: Timestamp1}, {Seconds: Timestamp1 + 60}}, timestamps)
}

func TestQueryTypedNilFormatter(t *testing.T) {
	var value *queryStringValue
	query := &Url_Query{}
	require.NotPanics(t, func() {
		require.NoError(t, query.Set("empty", value))
		require.Empty(t, ToString(value))
	})
	require.Empty(t, query.Vals["empty"].Vals)
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
			require.NoError(t, x.Add(tt.args.k, tt.args.v))
			require.Equal(t, tt.want, x)
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
		{name: "set-empty-value", Vals: map[string]*StringValues{"key1": {Vals: []string{"v1"}}}, args: args{k: "key1"}, want: &Url_Query{Vals: map[string]*StringValues{"key1": {}}}},
		{name: "set-nonempty-value", Vals: map[string]*StringValues{"key1": {Vals: []string{"v1"}}}, args: args{k: "key1", v: "v11"}, want: &Url_Query{Vals: map[string]*StringValues{"key1": {Vals: []string{"v11"}}}}},
		{name: "set-other-empty-key", Vals: map[string]*StringValues{"key1": {Vals: []string{"v1"}}}, args: args{k: "key2"}, want: &Url_Query{Vals: map[string]*StringValues{"key1": {Vals: []string{"v1"}}, "key2": {}}}},
		{name: "set-other-nonempty-key", Vals: map[string]*StringValues{"key1": {Vals: []string{"v1"}}}, args: args{k: "key2", v: "v2"}, want: &Url_Query{Vals: map[string]*StringValues{"key1": {Vals: []string{"v1"}}, "key2": {Vals: []string{"v2"}}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &Url_Query{
				Vals: tt.Vals,
			}
			require.NoError(t, x.Set(tt.args.k, tt.args.v))
			require.Equal(t, tt.want, x)
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
		{name: "string-array", query: newTestQuery(t, "foo", []string{"bar", "bba"}), k: "foo", typ: reflect.TypeOf([2]string{}), want: [2]string{"bar", "bba"}},
		{name: "timestamp", query: newTestQuery(t, "foo", TimestampString1), k: "foo", typ: reflect.TypeOf(Timestamp{}), want: Timestamp{Seconds: Timestamp1}},
		{name: "timestamp-slice", query: newTestQuery(t, "foo", []string{TimestampString1, TimestampString2}), k: "foo", typ: reflect.TypeOf([]*Timestamp{}), want: []*Timestamp{{Seconds: Timestamp1}, {Seconds: Timestamp1 + 60}}},
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
		{name: "timestamp-slice", str: TimestampString1 + "," + TimestampString2, typ: reflect.TypeOf([]*Timestamp{}), want: []*Timestamp{{Seconds: Timestamp1}, {Seconds: Timestamp1 + 60}}},
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

func Test_unmarshal(t *testing.T) {
	arrayStr := `["bar", "bba"]`
	var v []string
	err := jsoniter.ConfigFastest.UnmarshalFromString(arrayStr, &v)
	fmt.Println(err)
	fmt.Println(v)
}
