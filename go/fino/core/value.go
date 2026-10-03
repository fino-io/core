package core

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"unicode/utf8"

	jsoniter "github.com/json-iterator/go"
)

const ValueTypeName = "Value"
const ValueTypeFullName = "core.Value"

// NewValue converts Go scalars, maps, slices, arrays, and structs into Value.
// Integers retain int64/uint64 precision; decimal JSON numbers use float64.
// Structs and typed collections use jsoniter's registered codecs and JSON tags.
// Nil pointers become null. For compatibility, []byte becomes a base64 string;
// use NewBytesValue to retain the binary kind and its b64. JSON representation.
func NewValue(val any) (*Value, error) {
	switch v := val.(type) {
	case nil:
		return NewNullValue(), nil
	case bool:
		return NewBoolValue(v), nil
	case int:
		return NewIntValue(v), nil
	case int8:
		return NewInt8Value(v), nil
	case int16:
		return NewInt16Value(v), nil
	case int32:
		return NewInt32Value(v), nil
	case int64:
		return NewInt64Value(v), nil
	case uint:
		return NewUintValue(v), nil
	case uint8:
		return NewUint8Value(v), nil
	case uint16:
		return NewUint16Value(v), nil
	case uint32:
		return NewUint32Value(v), nil
	case uint64:
		return NewUint64Value(v), nil
	case float32:
		return NewFloat32Value(v), nil
	case float64:
		return NewFloat64Value(v), nil
	case json.Number:
		if !json.Valid([]byte(v)) {
			return nil, fmt.Errorf("invalid JSON number: %q", v)
		}
		return parseNumberValue(string(v))
	case *Value:
		if v == nil {
			return NewNullValue(), nil
		}
		return v, nil
	case *Object:
		if v == nil {
			return NewNullValue(), nil
		}
		return NewObjectValue(v), nil
	case *Values:
		if v == nil {
			return NewNullValue(), nil
		}
		return NewValuesValue(v), nil
	case string:
		if !utf8.ValidString(v) {
			return nil, fmt.Errorf("invalid UTF-8 in string: %q", v)
		}
		return NewStringValue(v), nil
	case []byte:
		s := base64.StdEncoding.EncodeToString(v)
		return NewStringValue(s), nil
	case map[string]any:
		v2, err := NewObjectFromMap(v)
		if err != nil {
			return nil, err
		}
		return NewObjectValue(v2), nil
	case []any:
		v2, err := NewValues(v)
		if err != nil {
			return nil, err
		}
		return NewValuesValue(v2), nil
	default:
		rv := reflect.ValueOf(val)
		for rv.Kind() == reflect.Ptr {
			if rv.IsNil() {
				return NewNullValue(), nil
			}
			rv = rv.Elem()
		}
		switch rv.Kind() {
		case reflect.Struct, reflect.Map, reflect.Slice, reflect.Array:
			data, err := jsoniter.Marshal(val)
			if err != nil {
				return nil, err
			}
			var value Value
			if err := jsoniter.Unmarshal(data, &value); err != nil {
				return nil, err
			}
			return &value, nil
		case reflect.Bool:
			return NewBoolValue(rv.Bool()), nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return NewInt64Value(rv.Int()), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return NewUint64Value(rv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return NewFloat64Value(rv.Float()), nil
		case reflect.String:
			return NewValue(rv.String())
		default:
			return nil, fmt.Errorf("invalid type: %T", v)
		}
	}
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

func NewIntArrayValue(vals ...int) *Value {
	_vals := make([]*Value, 0, len(vals))
	for _, v := range vals {
		_vals = append(_vals, NewIntValue(v))
	}
	return &Value{Val: &Value_ValuesValue{ValuesValue: &Values{Vals: _vals}}}
}

func NewInt32ArrayValue(vals ...int32) *Value {
	_vals := make([]*Value, 0, len(vals))
	for _, v := range vals {
		_vals = append(_vals, NewInt32Value(v))
	}
	return &Value{Val: &Value_ValuesValue{ValuesValue: &Values{Vals: _vals}}}
}

func NewInt64ArrayValue(vals ...int64) *Value {
	_vals := make([]*Value, 0, len(vals))
	for _, v := range vals {
		_vals = append(_vals, NewInt64Value(v))
	}
	return &Value{Val: &Value_ValuesValue{ValuesValue: &Values{Vals: _vals}}}
}

func NewUintArrayValue(vals ...uint) *Value {
	_vals := make([]*Value, 0, len(vals))
	for _, v := range vals {
		_vals = append(_vals, NewUintValue(v))
	}
	return &Value{Val: &Value_ValuesValue{ValuesValue: &Values{Vals: _vals}}}
}

func NewUint32ArrayValue(vals ...uint32) *Value {
	_vals := make([]*Value, 0, len(vals))
	for _, v := range vals {
		_vals = append(_vals, NewUint32Value(v))
	}
	return &Value{Val: &Value_ValuesValue{ValuesValue: &Values{Vals: _vals}}}
}

func NewUint64ArrayValue(vals ...uint64) *Value {
	_vals := make([]*Value, 0, len(vals))
	for _, v := range vals {
		_vals = append(_vals, NewUint64Value(v))
	}
	return &Value{Val: &Value_ValuesValue{ValuesValue: &Values{Vals: _vals}}}
}

func NewFloat32ArrayValue(vals ...float32) *Value {
	_vals := make([]*Value, 0, len(vals))
	for _, v := range vals {
		_vals = append(_vals, NewFloat32Value(v))
	}
	return &Value{Val: &Value_ValuesValue{ValuesValue: &Values{Vals: _vals}}}
}

func NewFloat64ArrayValue(vals ...float64) *Value {
	_vals := make([]*Value, 0, len(vals))
	for _, v := range vals {
		_vals = append(_vals, NewFloat64Value(v))
	}
	return &Value{Val: &Value_ValuesValue{ValuesValue: &Values{Vals: _vals}}}
}

func NewStringArrayValue(vals ...string) *Value {
	_vals := make([]*Value, 0, len(vals))
	for _, v := range vals {
		_vals = append(_vals, NewStringValue(v))
	}
	return &Value{Val: &Value_ValuesValue{ValuesValue: &Values{Vals: _vals}}}
}

func NewObjectArrayValue(vals ...*Object) *Value {
	_vals := make([]*Value, 0, len(vals))
	for _, v := range vals {
		_vals = append(_vals, NewObjectValue(v))
	}
	return &Value{Val: &Value_ValuesValue{ValuesValue: &Values{Vals: _vals}}}
}

// AsInterface converts x to a general-purpose Go interface.
//
// jsoniter.Marshal(x) and encoding/json.Marshal(x.AsInterface()) produce
// semantically equivalent JSON (assuming no errors occur).
//
// Floating-point values (i.e., "NaN", "Infinity", and "-Infinity") are
// converted as strings to remain compatible with MarshalJSON.
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
		switch {
		case math.IsNaN(v.NumberValue):
			return "NaN"
		case math.IsInf(v.NumberValue, +1):
			return "Infinity"
		case math.IsInf(v.NumberValue, -1):
			return "-Infinity"
		default:
			return v.NumberValue
		}
	case *Value_StringValue:
		return v.StringValue
	case *Value_BytesValue:
		return Base64Prefix + base64.StdEncoding.EncodeToString(v.BytesValue)
	case *Value_ObjectValue:
		if v.ObjectValue == nil {
			return nil
		}
		return v.ObjectValue.AsMap()
	case *Value_ValuesValue:
		if v.ValuesValue == nil {
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
	vals := x.GetValues()
	array := make([]bool, 0, len(vals))
	for _, v := range vals {
		array = append(array, v.GetBoolValue())
	}
	return array
}

func (x *Value) GetIntArray() []int {
	vals := x.GetValues()
	array := make([]int, 0, len(vals))
	for _, v := range vals {
		array = append(array, v.GetInt())
	}
	return array
}

func (x *Value) GetInt64Array() []int64 {
	vals := x.GetValues()
	array := make([]int64, 0, len(vals))
	for _, v := range vals {
		array = append(array, v.GetInt64())
	}
	return array
}

func (x *Value) GetUintArray() []uint {
	vals := x.GetValues()
	array := make([]uint, 0, len(vals))
	for _, v := range vals {
		array = append(array, v.GetUint())
	}
	return array
}

func (x *Value) GetFloat64Array() []float64 {
	vals := x.GetValues()
	array := make([]float64, 0, len(vals))
	for _, v := range vals {
		array = append(array, v.GetFloat64())
	}
	return array
}

func (x *Value) GetStringArray() []string {
	vals := x.GetValues()
	array := make([]string, 0, len(vals))
	for _, v := range vals {
		array = append(array, v.GetStringValue())
	}
	return array
}

func (x *Value) GetValueArray() []*Value {
	if values := x.GetValuesValue(); values != nil {
		return values.Vals
	}
	return nil
}

func (x *Value) GetObjectArray() []*Object {
	if values := x.GetValueArray(); len(values) > 0 {
		objs := make([]*Object, 0, len(values))
		for _, v := range values {
			objs = append(objs, v.GetObject())
		}
		return objs
	}
	return nil
}
