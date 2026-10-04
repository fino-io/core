package core

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"slices"
	"strings"
)

func NewUrlQuery(kvs ...any) (*Url_Query, error) {
	if len(kvs)%2 != 0 {
		return nil, fmt.Errorf("query requires key/value pairs")
	}
	query := &Url_Query{
		Vals: make(map[string]*StringValues),
	}

	for i := 0; i < len(kvs); i += 2 {
		key, ok := kvs[i].(string)
		if !ok {
			return nil, fmt.Errorf("query key must be a string: %T", kvs[i])
		}
		if err := query.Add(key, kvs[i+1]); err != nil {
			return nil, err
		}
	}
	return query, nil
}

func NewUrlQueryFrom(values url.Values) *Url_Query {
	return (&Url_Query{}).FromUrlValues(values)
}

// FromUrlValues replaces the query with an independent copy of values.
func (x *Url_Query) FromUrlValues(values url.Values) *Url_Query {
	if x != nil {
		x.Vals = make(map[string]*StringValues, len(values))
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

func (x *Url_Query) Add(k string, v any) error {
	values, err := queryValueFormat(v)
	if err != nil {
		return err
	}
	k = strings.TrimSpace(k)
	if x != nil && k != "" {
		if x.Vals == nil {
			x.Vals = make(map[string]*StringValues)
		}
		if x.Vals[k] == nil {
			x.Vals[k] = &StringValues{}
		}
		x.Vals[k].Vals = append(x.Vals[k].Vals, values...)
	}
	return nil
}

func (x *Url_Query) Set(k string, v any) error {
	values, err := queryValueFormat(v)
	if err != nil {
		return err
	}
	k = strings.TrimSpace(k)
	if x != nil && k != "" {
		if x.Vals == nil {
			x.Vals = make(map[string]*StringValues)
		}
		x.Vals[k] = &StringValues{
			Vals: values,
		}
	}
	return nil
}

func (x *Url_Query) Del(k string) *Url_Query {
	k = strings.TrimSpace(k)
	if x != nil && k != "" {
		delete(x.Vals, k)
	}
	return x
}

func queryValueFormat(val any) ([]string, error) {
	v := reflect.ValueOf(val)
	if !v.IsValid() || (v.Kind() == reflect.Ptr && v.IsNil()) {
		return nil, nil
	}
	if checked, ok := val.(interface{ CheckValid() error }); ok {
		if err := checked.CheckValid(); err != nil {
			return nil, err
		}
	}
	if values, ok := val.([]string); ok {
		return slices.Clone(values), nil
	}
	if text, ok := formatScalar(val); ok {
		return []string{text}, nil
	}
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return nil, fmt.Errorf("unsupported query value: %T", val)
	}
	var values []string
	if v.Kind() == reflect.Array || !v.IsNil() {
		values = make([]string, 0, v.Len())
	}
	for i := 0; i < v.Len(); i++ {
		element := v.Index(i).Interface()
		rv := reflect.ValueOf(element)
		if !rv.IsValid() || (rv.Kind() == reflect.Ptr && rv.IsNil()) {
			continue
		}
		if checked, ok := element.(interface{ CheckValid() error }); ok {
			if err := checked.CheckValid(); err != nil {
				return nil, fmt.Errorf("query element %d: %w", i, err)
			}
		}
		text, ok := formatScalar(element)
		if !ok {
			return nil, fmt.Errorf("unsupported query element %d: %T", i, element)
		}
		values = append(values, text)
	}
	return values, nil
}

// Unmarshal decodes repeated parameters as literal list elements. Use
// UnmarshalParam explicitly to parse a comma-separated or JSON-encoded list.
// Missing parameters leave the destination unchanged, as do decoding errors.
func (x *Url_Query) Unmarshal(name string, value any) error {
	param := x.GetVals()[name]
	if len(param.GetVals()) == 0 {
		return nil
	}

	v, err := paramDestination(value)
	if err != nil {
		return err
	}
	typ := v.Type()
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Slice && typ.Kind() != reflect.Array {
		return UnmarshalParam(param.Vals[0], value)
	}
	decoded := reflect.New(v.Type()).Elem()
	list := decoded
	for list.Kind() == reflect.Ptr {
		list.Set(reflect.New(list.Type().Elem()))
		list = list.Elem()
	}
	if list.Kind() == reflect.Array && list.Len() != len(param.Vals) {
		return fmt.Errorf("query %q: expected %d elements, got %d", name, list.Len(), len(param.Vals))
	}
	if list.Kind() == reflect.Slice {
		list.Set(reflect.MakeSlice(list.Type(), len(param.Vals), len(param.Vals)))
	}
	for i, raw := range param.Vals {
		if err := decodeParam(raw, list.Index(i)); err != nil {
			return fmt.Errorf("query %q element %d: %w", name, i, err)
		}
	}
	v.Set(decoded)
	return nil
}

// UnmarshalParam parses one scalar, a comma-separated list, or a JSON value.
// It replaces the destination only after the whole input succeeds.
func UnmarshalParam(str string, value any) error {
	v, err := paramDestination(value)
	if err != nil {
		return err
	}
	decoded := reflect.New(v.Type()).Elem()
	if err := decodeParam(str, decoded); err != nil {
		return err
	}
	v.Set(decoded)
	return nil
}

func decodeParam(str string, v reflect.Value) error {
	if v.Kind() == reflect.Ptr {
		v.Set(reflect.New(v.Type().Elem()))
		return decodeParam(str, v.Elem())
	}
	value := v.Addr().Interface()
	if parser, ok := value.(Parser); ok {
		return parser.Parse(str)
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
		return readJSON([]byte(str), value)
	default:
		if isStringParamType(v.Type()) {
			str = QuoteString(str)
		}
		err := readJSON([]byte(str), value)
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
		t.Implements(reflect.TypeFor[ToStringConverter]()) ||
		reflect.PointerTo(t).Implements(reflect.TypeFor[ToStringConverter]()) ||
		t.Implements(reflect.TypeFor[Formatter]()) ||
		reflect.PointerTo(t).Implements(reflect.TypeFor[Formatter]())
}
