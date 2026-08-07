package core

import (
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

func init() {
	RegisterJSONTypeDecoder(ValuesTypeFullName, &ValuesCodec{})
	RegisterJSONTypeEncoder(ValuesTypeFullName, &ValuesCodec{})
}

type ValuesCodec struct{}

func (codec *ValuesCodec) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	a := iter.ReadAny()
	values := (*Values)(ptr)
	if a.ValueType() == jsoniter.ArrayValue {
		a.ToVal(&values.Vals)
	}
}

func (codec *ValuesCodec) IsEmpty(ptr unsafe.Pointer) bool {
	values := (*Values)(ptr)
	return values == nil || len(values.Vals) == 0
}

func (codec *ValuesCodec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	values := (*Values)(ptr)
	stream.WriteVal(&values.Vals)
}
