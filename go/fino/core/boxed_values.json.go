package core

import (
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

func RegisterJSONValuesCodec[T any, M any](typ string, fn func(*M) *[]T) {
	codec := &ValsCodec[T, M]{GetVals: fn}
	RegisterJSONTypeDecoder(typ, codec)
	RegisterJSONTypeEncoder(typ, codec)
}

type ValsCodec[T any, M any] struct {
	GetVals func(*M) *[]T
}

func (codec *ValsCodec[T, M]) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	next := iter.WhatIsNext()
	if next != jsoniter.ArrayValue && next != jsoniter.NilValue {
		iter.ReportError("ValsCodec.Decode", "expected JSON array")
		return
	}
	destination := codec.GetVals((*M)(ptr))
	if destination == nil {
		iter.ReportError("ValsCodec.Decode", "nil slice destination")
		return
	}
	var values []T
	iter.ReadVal(&values)
	if iter.Error == nil {
		*destination = values
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
