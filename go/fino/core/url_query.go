package core

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"slices"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

// NewUrlQuery ignores non-string keys and an incomplete trailing pair.
func NewUrlQuery(kvs ...any) *Url_Query {
	query := &Url_Query{
		Vals: make(map[string]*StringValues),
	}

	for i := 0; i < len(kvs)-1; i += 2 {
		if key, ok := kvs[i].(string); ok {
			query.Add(key, kvs[i+1])
		}
	}
	return query
}

func NewUrlQueryFrom(values url.Values) *Url_Query {
	query := &Url_Query{
		Vals: make(map[string]*StringValues),
	}
	return query.FromUrlValues(values)
}

func (x *Url_Query) FromUrlValues(values url.Values) *Url_Query {
	if x != nil {
		if x.Vals == nil {
			x.Vals = make(map[string]*StringValues)
		}

		for k, v := range values {
			x.Vals[k] = &StringValues{Vals: slices.Clone(v)}
		}
	}
	return x
}

func (x *Url_Query) Has(k string) bool {
	if x != nil {
		_, ok := x.Vals[k]
		return ok
	}
	return false
}

func (x *Url_Query) Add(k string, v any) *Url_Query {
	k = strings.TrimSpace(k)
	if x != nil && k != "" {
		if x.Vals == nil {
			x.Vals = make(map[string]*StringValues)
		}
		if x.Vals[k] == nil {
			x.Vals[k] = &StringValues{}
		}
		x.Vals[k].Vals = append(x.Vals[k].Vals, queryValueFormat(v)...)
	}
	return x
}

func (x *Url_Query) Set(k string, v any) *Url_Query {
	k = strings.TrimSpace(k)
	if x != nil && k != "" {
		if x.Vals == nil {
			x.Vals = make(map[string]*StringValues)
		}
		x.Vals[k] = &StringValues{
			Vals: queryValueFormat(v),
		}
	}
	return x
}

func (x *Url_Query) Del(k string) *Url_Query {
	k = strings.TrimSpace(k)
	if x != nil && k != "" {
		delete(x.Vals, k)
	}
	return x
}

func queryValueFormat(val any) []string {
	switch v := val.(type) {
	case []string:
		return slices.Clone(v)
	case bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64,
		float32, float64, string, ToStringConverter, Formatter:
		return []string{ToString(val)}
	default:
		return []string{}
	}
}

// Unmarshal
// list type
//
//	foo=bar&foo=baz
//	foo=bar,baz
//	foo="bar","baz"
//	foo=["bar","baz"]
//
// map
//
//	foo=key1,bar,key2,baz <not support>
//
// object
//
//	foo={"key1":"bar","key2","baz"}
func (x *Url_Query) Unmarshal(name string, value any) error {
	param := x.GetVals()[name]
	if len(param.GetVals()) == 0 {
		return nil
	}

	v, err := paramDestination(value)
	if err != nil {
		return err
	}
	if (v.Kind() != reflect.Slice && v.Kind() != reflect.Array) || len(param.Vals) == 1 {
		return UnmarshalParam(param.Vals[0], value)
	}
	if isStringParamType(v.Type().Elem()) {
		data, err := json.Marshal(param.Vals)
		if err != nil {
			return err
		}
		return UnmarshalParam(string(data), value)
	}
	return UnmarshalParam("["+strings.Join(param.Vals, ",")+"]", value)
}

func UnmarshalParam(str string, value any) error {
	v, err := paramDestination(value)
	if err != nil {
		return err
	}
	if len(str) == 0 {
		return nil
	}

	switch v.Kind() {
	case reflect.String:
		v.SetString(str)
	case reflect.Slice, reflect.Array:
		str = strings.TrimSpace(str)
		if str == "" {
			return nil
		}
		if str[0] != '[' {
			if isStringParamType(v.Type().Elem()) {
				vals := splitQuotedString(str)
				for i, val := range vals {
					vals[i] = QuoteString(strings.TrimSpace(val))
				}
				str = strings.Join(vals, ",")
			}
			str = "[" + str + "]"
		}
		return jsoniter.Unmarshal([]byte(str), value)
	default:
		if reflect.PointerTo(v.Type()).Implements(reflect.TypeFor[StringLike]()) {
			str = QuoteString(str)
		}

		err := jsoniter.ConfigFastest.Unmarshal([]byte(str), value)
		if err != nil {
			return fmt.Errorf("couldn't decode value from %v, error: %w", str, err)
		}
	}

	return nil
}

func paramDestination(value any) (reflect.Value, error) {
	v := reflect.ValueOf(value)
	if !v.IsValid() || v.Kind() != reflect.Ptr || v.IsNil() {
		return reflect.Value{}, fmt.Errorf("expected a non-nil pointer, got %T", value)
	}
	return v.Elem(), nil
}

func splitQuotedString(str string) []string {
	if IsQuotedString(str, `"`) {
		var raw []json.RawMessage
		if json.Unmarshal([]byte("["+str+"]"), &raw) == nil {
			vals := make([]string, len(raw))
			for i, value := range raw {
				vals[i] = string(value)
			}
			return vals
		}
	}
	return strings.Split(str, ",")
}

func isStringParamType(t reflect.Type) bool {
	return t.Kind() == reflect.String ||
		(t.Kind() == reflect.Ptr && reflect.PointerTo(t.Elem()).Implements(reflect.TypeFor[StringLike]()))
}
