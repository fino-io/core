package core

import (
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

func NewUrlQuery(kvs ...any) *Url_Query {
	query := &Url_Query{
		Vals: make(map[string]*StringValues),
	}

	for i := 0; i < len(kvs)-1; i += 2 {
		key := kvs[i].(string)
		query.Add(key, kvs[i+1])
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
			x.Vals[k] = &StringValues{Vals: v}
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
	case bool:
		return []string{strconv.FormatBool(v)}
	case int8:
		return []string{strconv.FormatInt(int64(v), 10)}
	case int16:
		return []string{strconv.FormatInt(int64(v), 10)}
	case int32:
		return []string{strconv.FormatInt(int64(v), 10)}
	case int64:
		return []string{strconv.FormatInt(v, 10)}
	case int:
		return []string{strconv.FormatInt(int64(v), 10)}
	case uint8:
		return []string{strconv.FormatUint(uint64(v), 10)}
	case uint16:
		return []string{strconv.FormatUint(uint64(v), 10)}
	case uint32:
		return []string{strconv.FormatUint(uint64(v), 10)}
	case uint64:
		return []string{strconv.FormatUint(v, 10)}
	case uint:
		return []string{strconv.FormatUint(uint64(v), 10)}
	case float32:
		return []string{strconv.FormatFloat(float64(v), 'g', -1, 32)}
	case float64:
		return []string{strconv.FormatFloat(v, 'g', -1, 64)}
	case string:
		return []string{v}
	case []string:
		return v
	}

	if c, ok := val.(ToStringConverter); ok {
		return []string{c.ToString()}
	}

	if f, ok := val.(Formatter); ok {
		return []string{f.Format()}
	}

	return []string{}
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
func (x *Url_Query) Unmarshal(name string, value interface{}) error {
	if x == nil || x.Vals == nil {
		return nil
	}

	param, ok := x.Vals[name]
	if !ok || param == nil || len(param.Vals) == 0 {
		return nil
	}

	v := reflect.Indirect(reflect.ValueOf(value))
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		elem := v.Type().Elem()
		isStringType := elem.Kind() == reflect.String
		if elem.Kind() == reflect.Ptr {
			if _, ok := reflect.New(elem.Elem()).Interface().(StringLike); ok {
				isStringType = true
			}
		}

		var sliceValue string
		if len(param.Vals) == 1 {
			sliceValue = param.Vals[0]
		} else {
			if isStringType {
				for i := 0; i < len(param.Vals); i++ {
					param.Vals[i] = QuoteString(param.Vals[i])
				}
			}
			sliceValue = "[" + strings.Join(param.Vals, ",") + "]"
		}
		return UnmarshalParam(sliceValue, value)
	default:
		if len(param.Vals) > 0 {
			return UnmarshalParam(param.Vals[0], value)
		}
	}
	return nil
}

func UnmarshalParam(str string, value interface{}) error {
	if len(str) == 0 {
		return nil
	}

	v := reflect.Indirect(reflect.ValueOf(value))
	switch v.Kind() {
	case reflect.String:
		v.SetString(str)
	case reflect.Slice, reflect.Array:
		elem := v.Type().Elem()
		isStringType := elem.Kind() == reflect.String
		if elem.Kind() == reflect.Ptr {
			if _, ok := reflect.New(elem.Elem()).Interface().(StringLike); ok {
				isStringType = true
			}
		}

		str = strings.TrimSpace(str)
		if len(str) > 0 {
			if str[0] != '[' {
				if isStringType {
					vals := splitQuotedString(str)
					for i := 0; i < len(vals); i++ {
						vals[i] = strings.TrimSpace(vals[i])
						vals[i] = QuoteString(vals[i])
					}
					str = "[" + strings.Join(vals, ",") + "]"
				} else {
					str = "[" + str + "]"
				}
			}
			return jsoniter.Unmarshal([]byte(str), value)
		}
		return nil
	default:
		if _, ok := reflect.New(v.Type()).Interface().(StringLike); ok {
			str = QuoteString(str)
		}

		err := jsoniter.ConfigFastest.Unmarshal([]byte(str), value)
		if err != nil {
			return fmt.Errorf("couldn't decode value from %v, error: %w", str, err)
		}
	}

	return nil
}

var separator = regexp.MustCompile(`" *, *"`)

func splitQuotedString(str string) []string {
	if !IsQuotedString(str, `"`) {
		return strings.Split(str, ",")
	}

	vals := separator.Split(str, -1)
	if len(vals) > 1 {
		vals[0] = vals[0] + `"`
		vals[len(vals)-1] = `"` + vals[len(vals)-1]

		for i := 1; i < len(vals)-1; i++ {
			vals[i] = QuoteString(vals[i])
		}
	}
	return vals
}

func isStringSlice(v any) bool {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() == reflect.Slice {
		t = t.Elem()
		if t.Kind() == reflect.String {
			return true
		}
	}
	return false
}
