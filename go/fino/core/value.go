package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"unicode/utf8"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

const ValueTypeName = "Value"
const ValueTypeFullName = "core.Value"

// NewValue converts Go scalars, string-keyed maps, slices, arrays, and structs.
// Native collections preserve element kinds, integer precision, and binary data.
// Structs use jsoniter's registered codecs and JSON tags. Nil pointers become null.
func NewValue(val any) (*Value, error) {
	return newValue(val, make(map[valueVisit]bool))
}

type valueVisit struct {
	typ    reflect.Type
	ptr    unsafe.Pointer
	length int
}

func newValue(val any, visiting map[valueVisit]bool) (*Value, error) {
	rv := reflect.ValueOf(val)
	if !rv.IsValid() || (rv.Kind() == reflect.Ptr && rv.IsNil()) {
		return NewNullValue(), nil
	}
	switch v := val.(type) {
	case json.Number:
		return parseNumberValue(string(v))
	case *Value:
		return v, nil
	case *Object:
		return NewObjectValue(v), nil
	case *Values:
		return NewValuesValue(v), nil
	}

	// Track the current recursion path; repeated references outside it are valid.
	switch rv.Kind() {
	case reflect.Map, reflect.Slice, reflect.Ptr:
		if !rv.IsNil() {
			visit := valueVisit{typ: rv.Type(), ptr: rv.UnsafePointer()}
			if rv.Kind() == reflect.Slice {
				visit.length = rv.Len()
			}
			if visiting[visit] {
				return nil, fmt.Errorf("cyclic value: %T", val)
			}
			visiting[visit] = true
			defer delete(visiting, visit)
		}
	}

	switch rv.Kind() {
	case reflect.Bool:
		return NewBoolValue(rv.Bool()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return NewInt64Value(rv.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return NewUint64Value(rv.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return NewFloat64Value(rv.Float()), nil
	case reflect.String:
		text := rv.String()
		if !utf8.ValidString(text) {
			return nil, fmt.Errorf("invalid UTF-8 in string: %q", text)
		}
		return NewStringValue(text), nil
	case reflect.Map:
		if rv.Type().Key().Kind() != reflect.String {
			return nil, fmt.Errorf("object keys must be strings: %T", val)
		}
		var values map[string]*Value
		if !rv.IsNil() {
			values = make(map[string]*Value, rv.Len())
		}
		entries := rv.MapRange()
		for entries.Next() {
			key := entries.Key().String()
			if !utf8.ValidString(key) {
				return nil, fmt.Errorf("invalid UTF-8 in object key: %q", key)
			}
			value, err := newValue(entries.Value().Interface(), visiting)
			if err != nil {
				return nil, fmt.Errorf("object field %q: %w", key, err)
			}
			values[key] = value
		}
		return NewMapValue(values), nil
	case reflect.Slice, reflect.Array:
		if rv.Kind() == reflect.Slice && rv.Type().Elem().Kind() == reflect.Uint8 {
			return NewBytesValue(rv.Bytes()), nil
		}
		var values []*Value
		if rv.Kind() == reflect.Array || !rv.IsNil() {
			values = make([]*Value, rv.Len())
		}
		for i := 0; i < rv.Len(); i++ {
			value, err := newValue(rv.Index(i).Interface(), visiting)
			if err != nil {
				return nil, fmt.Errorf("array element %d: %w", i, err)
			}
			values[i] = value
		}
		return NewArrayValue(values...), nil
	case reflect.Ptr:
		if rv.Elem().Kind() != reflect.Struct {
			return newValue(rv.Elem().Interface(), visiting)
		}
		fallthrough
	case reflect.Struct:
		data, err := jsoniter.Marshal(val)
		if err != nil {
			return nil, err
		}
		return valueFromJSON(data)
	default:
		return nil, fmt.Errorf("invalid type: %T", val)
	}
}

// valueFromJSON reads ordinary Go JSON without interpreting Value wire prefixes.
func valueFromJSON(data []byte) (*Value, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		return nil, err
	}
	return NewValue(value)
}

// NewNullValue constructs a new null Value.
func NewNullValue() *Value {
	return &Value{Val: &Value_NullValue{NullValue: &Null{}}}
}

// NewBoolValue constructs a new boolean Value.
func NewBoolValue(v bool) *Value {
	return &Value{Val: &Value_BoolValue{BoolValue: v}}
}

// NewNumberValue constructs a new number Value.
func NewNumberValue(v float64) *Value {
	return &Value{Val: &Value_NumberValue{NumberValue: v}}
}

func NewIntValue(v int) *Value {
	return NewInt64Value(int64(v))
}

func NewInt8Value(v int8) *Value {
	return NewInt64Value(int64(v))
}

func NewInt16Value(v int16) *Value {
	return NewInt64Value(int64(v))
}

func NewInt32Value(v int32) *Value {
	return NewInt64Value(int64(v))
}

func NewInt64Value(v int64) *Value {
	if v >= 0 {
		return &Value{Val: &Value_PositiveValue{PositiveValue: uint64(v)}}
	}

	var val uint64
	if v == math.MinInt64 {
		// val = uint64(math.MinInt64) + 1
		val = uint64(1) << 63
	} else {
		val = uint64(-v)
	}

	return &Value{Val: &Value_NegativeValue{NegativeValue: val}}
}

func NewUintValue(v uint) *Value {
	return NewUint64Value(uint64(v))
}

func NewUint8Value(v uint8) *Value {
	return NewUint64Value(uint64(v))
}

func NewUint16Value(v uint16) *Value {
	return NewUint64Value(uint64(v))
}

func NewUint32Value(v uint32) *Value {
	return NewUint64Value(uint64(v))
}

func NewUint64Value(v uint64) *Value {
	return &Value{Val: &Value_PositiveValue{PositiveValue: v}}
}

func NewNegativeValue(v uint64) *Value {
	return &Value{Val: &Value_NegativeValue{NegativeValue: v}}
}

func NewPositiveValue(v uint64) *Value {
	return &Value{Val: &Value_PositiveValue{PositiveValue: v}}
}

func NewFloat32Value(v float32) *Value {
	return NewFloat64Value(float64(v))
}

func NewFloat64Value(v float64) *Value {
	return &Value{Val: &Value_NumberValue{NumberValue: v}}
}

// NewStringValue constructs a new string Value.
func NewStringValue(v string) *Value {
	return &Value{Val: &Value_StringValue{StringValue: v}}
}

func NewBytesValue(v []byte) *Value {
	return &Value{Val: &Value_BytesValue{BytesValue: v}}
}

func NewMapValue(v map[string]*Value) *Value {
	return &Value{Val: &Value_ObjectValue{ObjectValue: &Object{Vals: v}}}
}

// NewObjectValue constructs a new struct Value.
func NewObjectValue(obj *Object) *Value {
	return &Value{Val: &Value_ObjectValue{ObjectValue: obj}}
}

// NewValuesValue constructs a new list Value.
func NewValuesValue(vals *Values) *Value {
	return &Value{Val: &Value_ValuesValue{ValuesValue: vals}}
}

func NewArrayValue(vals ...*Value) *Value {
	return &Value{Val: &Value_ValuesValue{ValuesValue: &Values{Vals: vals}}}
}

// mapSlice preserves nil and empty slices while converting their elements.
func mapSlice[T, R any](values []T, convert func(T) R) []R {
	if values == nil {
		return nil
	}
	result := make([]R, len(values))
	for i, value := range values {
		result[i] = convert(value)
	}
	return result
}

func NewIntArrayValue(vals ...int) *Value {
	return NewArrayValue(mapSlice(vals, NewIntValue)...)
}

func NewInt32ArrayValue(vals ...int32) *Value {
	return NewArrayValue(mapSlice(vals, NewInt32Value)...)
}

func NewInt64ArrayValue(vals ...int64) *Value {
	return NewArrayValue(mapSlice(vals, NewInt64Value)...)
}

func NewUintArrayValue(vals ...uint) *Value {
	return NewArrayValue(mapSlice(vals, NewUintValue)...)
}

func NewUint32ArrayValue(vals ...uint32) *Value {
	return NewArrayValue(mapSlice(vals, NewUint32Value)...)
}

func NewUint64ArrayValue(vals ...uint64) *Value {
	return NewArrayValue(mapSlice(vals, NewUint64Value)...)
}

func NewFloat32ArrayValue(vals ...float32) *Value {
	return NewArrayValue(mapSlice(vals, NewFloat32Value)...)
}

func NewFloat64ArrayValue(vals ...float64) *Value {
	return NewArrayValue(mapSlice(vals, NewFloat64Value)...)
}

func NewStringArrayValue(vals ...string) *Value {
	return NewArrayValue(mapSlice(vals, NewStringValue)...)
}

func NewObjectArrayValue(vals ...*Object) *Value {
	return NewArrayValue(mapSlice(vals, NewObjectValue)...)
}

// AsInterface returns native Go values, including []byte and non-finite floats.
// Use the Value JSON codec when the wire representation is required.
func (x *Value) AsInterface() any {
	switch v := x.GetVal().(type) {
	case *Value_NullValue:
		return nil
	case *Value_BoolValue:
		return v.BoolValue
	case *Value_PositiveValue:
		return v.PositiveValue
	case *Value_NegativeValue:
		return negativeValueToInt64(v.NegativeValue)
	case *Value_NumberValue:
		return v.NumberValue
	case *Value_StringValue:
		return v.StringValue
	case *Value_BytesValue:
		return v.BytesValue
	case *Value_ObjectValue:
		if v.ObjectValue.GetVals() == nil {
			return nil
		}
		return v.ObjectValue.AsMap()
	case *Value_ValuesValue:
		if v.ValuesValue.GetVals() == nil {
			return nil
		}
		return v.ValuesValue.AsSlice()
	default:
		return v
	}
}

func (x *Value) GetKind() ValueKind {
	if x != nil {
		switch x.GetVal().(type) {
		case *Value_NullValue:
			return ValueKind_VALUE_KIND_NULL
		case *Value_BoolValue:
			return ValueKind_VALUE_KIND_BOOLEAN
		case *Value_NegativeValue, *Value_PositiveValue:
			return ValueKind_VALUE_KIND_INTEGER
		case *Value_NumberValue:
			return ValueKind_VALUE_KIND_NUMBER
		case *Value_StringValue:
			return ValueKind_VALUE_KIND_STRING
		case *Value_BytesValue:
			return ValueKind_VALUE_KIND_BYTES
		case *Value_ObjectValue:
			return ValueKind_VALUE_KIND_OBJECT
		case *Value_ValuesValue:
			return ValueKind_VALUE_KIND_ARRAY
		}
	}
	return ValueKind_VALUE_KIND_UNSPECIFIED
}

func (x *Value) GetBool() bool {
	return x.GetBoolValue()
}

func (x *Value) GetInt() int {
	return int(x.GetInt64())
}

func (x *Value) GetInt32() int32 {
	return int64ToInt32(x.GetInt64())
}

func (x *Value) GetInt64() int64 {
	if negative := x.GetNegativeValue(); negative > 0 {
		return negativeValueToInt64(negative)
	}
	return uint64ToInt64(x.GetPositiveValue())
}

func (x *Value) GetUint() uint {
	return uint(x.GetPositiveValue())
}

func (x *Value) GetUint32() uint32 {
	return uint64ToUint32(x.GetPositiveValue())
}

func negativeValueToInt64(value uint64) int64 {
	if value >= uint64(1)<<63 {
		return math.MinInt64
	}
	return -int64(value) // #nosec G115 -- value is strictly below 2^63 here.
}

func uint64ToInt64(value uint64) int64 {
	if value > math.MaxInt64 {
		return math.MaxInt64
	}
	return int64(value)
}

func int64ToInt32(value int64) int32 {
	if value > math.MaxInt32 {
		return math.MaxInt32
	}
	if value < math.MinInt32 {
		return math.MinInt32
	}
	return int32(value)
}

func uint64ToUint32(value uint64) uint32 {
	if value > math.MaxUint32 {
		return math.MaxUint32
	}
	return uint32(value)
}

func (x *Value) GetUint64() uint64 {
	return x.GetPositiveValue()
}

func (x *Value) GetFloat32() float32 {
	return float32(x.GetFloat64())
}

func (x *Value) GetFloat64() float64 {
	return x.GetNumberValue()
}

func (x *Value) GetDouble() float64 {
	return x.GetNumberValue()
}

func (x *Value) GetString() string {
	return x.GetStringValue()
}

func (x *Value) GetBytes() []byte {
	return x.GetBytesValue()
}

func (x *Value) GetObject() *Object {
	return x.GetObjectValue()
}

func (x *Value) GetValues() []*Value {
	if vals := x.GetValuesValue(); vals != nil {
		return vals.Vals
	}
	return nil
}

func (x *Value) GetBoolArray() []bool {
	return mapSlice(x.GetValues(), (*Value).GetBool)
}

func (x *Value) GetIntArray() []int {
	return mapSlice(x.GetValues(), (*Value).GetInt)
}

func (x *Value) GetInt64Array() []int64 {
	return mapSlice(x.GetValues(), (*Value).GetInt64)
}

func (x *Value) GetUintArray() []uint {
	return mapSlice(x.GetValues(), (*Value).GetUint)
}

func (x *Value) GetFloat64Array() []float64 {
	return mapSlice(x.GetValues(), (*Value).GetFloat64)
}

func (x *Value) GetStringArray() []string {
	return mapSlice(x.GetValues(), (*Value).GetString)
}

func (x *Value) GetObjectArray() []*Object {
	return mapSlice(x.GetValues(), (*Value).GetObject)
}
