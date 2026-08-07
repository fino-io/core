package core

import (
	"encoding/base64"
	"errors"
	"math"
	"strings"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

const Base64Prefix = "b64."

func init() {
	RegisterJSONTypeDecoder(ValueTypeFullName, &ValueCodec{})
	RegisterJSONTypeEncoder(ValueTypeFullName, &ValueCodec{})
}

type ValueCodec struct{}

func NewValueCodec() *ValueCodec {
	return &ValueCodec{}
}

func (codec *ValueCodec) DecodeAny(a jsoniter.Any) (*Value, error) {
	switch a.ValueType() {
	case jsoniter.NilValue:
		return &Value{}, nil
	case jsoniter.BoolValue:
		return NewBoolValue(a.ToBool()), nil
	case jsoniter.NumberValue:
		floatVal := a.ToFloat64()
		intVal := a.ToInt64()
		uintVal := a.ToUint64()

		if intVal == 0 && uintVal > 0 { // > int64 max
			return NewUint64Value(uintVal), nil
		} else if floatVal != float64(intVal) {
			return NewFloat64Value(floatVal), nil
		}
		return NewInt64Value(intVal), nil
	case jsoniter.StringValue:
		str := a.ToString()
		if strings.HasPrefix(str, Base64Prefix) {
			ds, err := base64.StdEncoding.DecodeString(str[len(Base64Prefix):])
			if err != nil {
				return nil, err
			}
			return NewBytesValue(ds), nil
		}

		switch str {
		case "NaN":
			return NewFloat64Value(math.NaN()), nil
		case "Infinity":
			return NewFloat64Value(math.Inf(1)), nil
		case "-Infinity":
			return NewFloat64Value(math.Inf(-1)), nil
		default:
			return NewStringValue(str), nil
		}
	case jsoniter.ObjectValue:
		val := make(map[string]*Value)
		a.ToVal(&val)
		return NewMapValue(val), nil
	case jsoniter.ArrayValue:
		val := make([]*Value, 0)
		a.ToVal(&val)
		return NewArrayValue(val...), nil
	default:
		return nil, errors.New("type is invalid")
	}
}

func (codec *ValueCodec) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	a := iter.ReadAny()
	v, _ := codec.DecodeAny(a)
	(*Value)(ptr).Val = v.Val
}

func (codec *ValueCodec) IsEmpty(ptr unsafe.Pointer) bool {
	v := (*Value)(ptr)
	return v == nil || v.Val == nil
}

func (codec *ValueCodec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	val := (*Value)(ptr)
	switch v := val.Val.(type) {
	case *Value_BoolValue:
		stream.WriteBool(v.BoolValue)
	case *Value_PositiveValue:
		stream.WriteUint64(v.PositiveValue)
	case *Value_NegativeValue:
		stream.WriteInt64(val.GetInt64())
	case *Value_NumberValue:
		stream.WriteFloat64Lossy(v.NumberValue)
	case *Value_StringValue:
		stream.WriteString(v.StringValue)
	case *Value_BytesValue:
		stream.WriteString(Base64Prefix + base64.StdEncoding.EncodeToString(v.BytesValue))
	case *Value_ValuesValue:
		stream.WriteVal(v.ValuesValue.Vals)
	case *Value_ObjectValue:
		stream.WriteVal(v.ObjectValue.Vals)
	default:
		stream.WriteNil()
	}
}
