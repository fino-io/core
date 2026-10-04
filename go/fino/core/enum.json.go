package core

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

// RegisterJSONEnum shares strict enum encoding across generated helpers.
func RegisterJSONEnum[T ~int32](typ string, names map[int32]string, values map[string]int32) {
	codec := &enumCodec[T]{names: names, values: values}
	RegisterJSONTypeDecoder(typ, codec)
	RegisterJSONTypeEncoder(typ, codec)
}

type enumCodec[T ~int32] struct {
	names  map[int32]string
	values map[string]int32
}

func (codec *enumCodec[T]) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	a := iter.ReadAny()
	if iter.Error != nil && iter.Error != io.EOF {
		return
	}
	var value int32
	switch a.ValueType() {
	case jsoniter.NilValue:
	case jsoniter.StringValue:
		var ok bool
		value, ok = codec.values[a.ToString()]
		if !ok {
			iter.ReportError("decode enum", "unknown enum name")
			return
		}
	case jsoniter.NumberValue:
		raw := a.ToString()
		number, err := strconv.ParseInt(raw, 10, 32)
		if err != nil || !json.Valid([]byte(raw)) {
			iter.ReportError("decode enum", "expected a JSON integer")
			return
		}
		value = int32(number)
		if _, ok := codec.names[value]; !ok {
			iter.ReportError("decode enum", "unknown enum value")
			return
		}
	default:
		iter.ReportError("decode enum", "expected an enum name or integer")
		return
	}
	*(*T)(ptr) = T(value)
}

func (codec *enumCodec[T]) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	value := int32(*(*T)(ptr))
	name, ok := codec.names[value]
	if !ok {
		stream.Error = fmt.Errorf("unknown enum value: %d", value)
		return
	}
	stream.WriteVal(name)
}

func (codec *enumCodec[T]) IsEmpty(ptr unsafe.Pointer) bool {
	return *(*T)(ptr) == 0
}
