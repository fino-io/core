package core

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"unicode/utf8"
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

// Track the active Value path on the stream so jsoniter can reject cycles
// without repeatedly validating every subtree.
type valueEncodePath map[*Value]bool

func decodeScalarValue(a jsoniter.Any) (*Value, error) {
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
		if !utf8.ValidString(str) {
			return nil, fmt.Errorf("invalid UTF-8 in string: %q", str)
		}
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
	default:
		return nil, fmt.Errorf("expected JSON scalar")
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
	if !strings.ContainsAny(raw, ".eE") {
		return nil, fmt.Errorf("json integer out of range: %q", raw)
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
		values := readJSONObject(iter)
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
	v, err := decodeScalarValue(a)
	if err != nil {
		iter.ReportError("valueCodec.Decode", err.Error())
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
	if err := val.checkScalar(); err != nil {
		stream.Error = err
		return
	}
	switch val.Val.(type) {
	case *Value_ValuesValue, *Value_ObjectValue:
		previous := stream.Attachment
		path, ok := previous.(valueEncodePath)
		if !ok {
			path = make(valueEncodePath)
			stream.Attachment = path
			defer func() { stream.Attachment = previous }()
		}
		if path[val] {
			stream.Error = fmt.Errorf("cyclic value")
			return
		}
		path[val] = true
		defer delete(path, val)
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
			number := strconv.FormatFloat(v.NumberValue, 'g', -1, 64)
			if !strings.ContainsAny(number, ".eE") {
				number += ".0"
			}
			stream.WriteRaw(number)
		}
	case *Value_StringValue:
		stream.WriteVal(jsonString(v.StringValue))
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

func (x *Value) MarshalJSON() ([]byte, error) {
	if err := x.CheckValid(); err != nil {
		return nil, err
	}
	return marshalJSONCodec(x, &ValueCodec{})
}

func (x *Value) UnmarshalJSON(data []byte) error {
	if x == nil {
		return fmt.Errorf("value.UnmarshalJSON: nil receiver")
	}
	value, err := unmarshalJSONCodec[Value](data, &ValueCodec{})
	if err == nil {
		x.Val = value.Val
	}
	return err
}
