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

// NewValue constructs a Value from a general-purpose Go interface.
//
//	╔═══════════════════════════════════════╤════════════════════════════════════════════╗
//	║ Go type                               │ Conversion                                 ║
//	╠═══════════════════════════════════════╪════════════════════════════════════════════╣
//	║ nil                                   │ stored as NullValue                        ║
//	║ bool                                  │ stored as BoolValue                        ║
//	║ int, int8, int16, int32, int64        │ stored as NegativeVal/PositiveVal          ║
//	║ uint, uint8, uint16, uint32, uint64   │ stored as NegativeVal/PositiveVal          ║
//	║ float32, float64                      │ stored as NumberValue                      ║
//	║ json.Number                           │ stored as NumberValue                      ║
//	║ string                                │ stored as StringValue; must be valid UTF-8 ║
//	║ []byte                                │ stored as StringValue; base64-encoded      ║
//	║ map[string]any                        │ stored as ObjectValue                      ║
//	║ []any                                 │ stored as ValuesValue                      ║
//	╚═══════════════════════════════════════╧════════════════════════════════════════════╝
//
// When converting an int64 or uint64 to a NumberValue, numeric precision loss
// is possible since they are stored as a float64.
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
		n, err := v.Float64()
		if err != nil {
			return nil, fmt.Errorf("invalid number format %q, expected a float64: %v", v, err)
		}
		return NewNumberValue(n), nil
	case string:
		if !utf8.ValidString(v) {
			return nil, fmt.Errorf("invalid UTF-8 in string: %q", v)
		}
		return NewStringValue(v), nil
	case []byte:
		s := base64.StdEncoding.EncodeToString(v)
		return NewStringValue(s), nil
		// return NewBytesValue(v), nil
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
		_val := reflect.ValueOf(val)
		if _val.Kind() == reflect.Ptr {
			_val = _val.Elem()
		}
		typ := reflect.Indirect(_val).Type()

		switch typ.Kind() {
		case reflect.Struct:
			m, err := structToMap(val)
			if err != nil {
				return nil, err
			}

			obj, err := NewObjectFromMap(m)
			if err != nil {
				return nil, err
			}
			return NewObjectValue(obj), nil
			// return nil, fmt.Errorf("struct type %T must be explicitly converted", val)
		default:
			return nil, fmt.Errorf("invalid type: %T", v)
		}
	}
}

func structToMap(val any) (map[string]any, error) {
	bytes, err := jsoniter.Marshal(val)
	if err != nil {
		return nil, err
	}

	var m map[string]any
	if err := jsoniter.Unmarshal(bytes, &m); err != nil {
		return nil, err
	}

	return m, nil
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
// Calling Value.MarshalJSON and "encoding/json".Marshal on this output produce
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
		if v.NegativeValue == uint64(1)<<63 {
			return math.MinInt64
		}
		return -int64(v.NegativeValue)
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
	case *Value_ObjectValue:
		return v.ObjectValue.AsMap()
	case *Value_ValuesValue:
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
	return int32(x.GetInt64())
}

func (x *Value) GetInt64() int64 {
	if negative := x.GetNegativeValue(); negative > 0 {
		if negative == uint64(1)<<63 {
			return math.MinInt64
		}
		return -int64(negative)
	}
	return int64(x.GetPositiveValue())
}

func (x *Value) GetUint() uint {
	return uint(x.GetPositiveValue())
}

func (x *Value) GetUint32() uint32 {
	return uint32(x.GetPositiveValue())
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
