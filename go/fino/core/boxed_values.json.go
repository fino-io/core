package core

import (
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

// func init() {
// 	RegisterJSONValuesCodec[bool, BoolValues]("core.BoolValues", func(b *BoolValues) *[]bool { return &b.Vals })
// }

func RegisterJSONValuesCodec[T any, M any](typ string, fn func(*M) *[]T) {
	codec := &ValsCodec[T, M]{GetVals: fn}
	RegisterJSONTypeDecoder(typ, codec)
	RegisterJSONTypeEncoder(typ, codec)
}

type ValsCodec[T any, M any] struct {
	GetVals func(*M) *[]T
}

func (codec *ValsCodec[T, M]) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	msg := (*M)(ptr)
	vals := codec.GetVals(msg)
	a := iter.ReadAny()
	if a.ValueType() == jsoniter.ArrayValue {
		a.ToVal(&vals)
	}
}

func (codec *ValsCodec[T, M]) IsEmpty(ptr unsafe.Pointer) bool {
	msg := (*M)(ptr)
	vals := codec.GetVals(msg)
	return vals == nil || len(*vals) == 0
}

func (codec *ValsCodec[T, M]) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	msg := (*M)(ptr)
	vals := codec.GetVals(msg)
	stream.WriteVal(vals)
}
