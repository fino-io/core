package core

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

const Base64Prefix = "b64."
const StringPrefix = "str."

func jsonString(value string) string {
	if strings.HasPrefix(value, Base64Prefix) || strings.HasPrefix(value, StringPrefix) ||
		value == "NaN" || value == "Infinity" || value == "-Infinity" {
		return StringPrefix + value
	}
	return value
}

func init() {
	RegisterJSONTypeDecoder(ValueTypeFullName, &ValueCodec{})
	RegisterJSONTypeEncoder(ValueTypeFullName, &ValueCodec{})
}

type ValueCodec struct{}

func NewValueCodec() *ValueCodec {
	return &ValueCodec{}
}

func (codec *ValueCodec) DecodeAny(a jsoniter.Any) (*Value, error) {
	if err := a.LastError(); err != nil {
		return nil, err
	}
	switch a.ValueType() {
	case jsoniter.NilValue:
		return NewNullValue(), nil
	case jsoniter.BoolValue:
		return NewBoolValue(a.ToBool()), nil
	case jsoniter.NumberValue:
		return parseNumberValue(a.ToString())
	case jsoniter.StringValue:
		str := a.ToString()
		if strings.HasPrefix(str, StringPrefix) {
			return NewStringValue(strings.TrimPrefix(str, StringPrefix)), nil
		}
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
		if err := jsoniter.UnmarshalFromString(a.ToString(), &val); err != nil {
			return nil, err
		}
		return NewMapValue(val), nil
	case jsoniter.ArrayValue:
		val := make([]*Value, 0)
		if err := jsoniter.UnmarshalFromString(a.ToString(), &val); err != nil {
			return nil, err
		}
		return NewArrayValue(val...), nil
	default:
		return nil, errors.New("type is invalid")
	}
}

func parseNumberValue(raw string) (*Value, error) {
	if !json.Valid([]byte(raw)) {
		return nil, fmt.Errorf("invalid JSON number: %q", raw)
	}
	if value, err := strconv.ParseInt(raw, 10, 64); err == nil {
		return NewInt64Value(value), nil
	}
	if value, err := strconv.ParseUint(raw, 10, 64); err == nil {
		return NewUint64Value(value), nil
	}
	value, err := strconv.ParseFloat(raw, 64)
	if err != nil || math.IsNaN(value) || math.IsInf(value, 0) {
		return nil, fmt.Errorf("invalid JSON number: %q", raw)
	}
	return NewFloat64Value(value), nil
}

func (codec *ValueCodec) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	// Decode collections on the same iterator so nested errors and depth limits
	// propagate without reparsing every enclosing object or array.
	switch iter.WhatIsNext() {
	case jsoniter.ObjectValue:
		var values map[string]*Value
		iter.ReadVal(&values)
		if iter.Error == nil {
			(*Value)(ptr).Val = NewMapValue(values).Val
		}
		return
	case jsoniter.ArrayValue:
		var values []*Value
		iter.ReadVal(&values)
		if iter.Error == nil {
			(*Value)(ptr).Val = NewArrayValue(values...).Val
		}
		return
	}
	a := iter.ReadAny()
	if iter.Error != nil && iter.Error != io.EOF {
		return
	}
	v, err := codec.DecodeAny(a)
	if err != nil {
		iter.ReportError("ValueCodec.Decode", err.Error())
		return
	}
	(*Value)(ptr).Val = v.Val
}

func (codec *ValueCodec) IsEmpty(ptr unsafe.Pointer) bool {
	v := (*Value)(ptr)
	return v == nil || v.Val == nil
}

func (codec *ValueCodec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	val := (*Value)(ptr)
	if val == nil {
		stream.WriteNil()
		return
	}
	switch v := val.Val.(type) {
	case *Value_BoolValue:
		stream.WriteBool(v.BoolValue)
	case *Value_PositiveValue:
		stream.WriteUint64(v.PositiveValue)
	case *Value_NegativeValue:
		stream.WriteInt64(val.GetInt64())
	case *Value_NumberValue:
		switch {
		case math.IsNaN(v.NumberValue):
			stream.WriteString("NaN")
		case math.IsInf(v.NumberValue, 1):
			stream.WriteString("Infinity")
		case math.IsInf(v.NumberValue, -1):
			stream.WriteString("-Infinity")
		default:
			stream.WriteFloat64(v.NumberValue)
		}
	case *Value_StringValue:
		stream.WriteString(jsonString(v.StringValue))
	case *Value_BytesValue:
		stream.WriteString(Base64Prefix + base64.StdEncoding.EncodeToString(v.BytesValue))
	case *Value_ValuesValue:
		stream.WriteVal(v.ValuesValue)
	case *Value_ObjectValue:
		stream.WriteVal(v.ObjectValue)
	default:
		stream.WriteNil()
	}
}
