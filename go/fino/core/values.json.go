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
	next := iter.WhatIsNext()
	if next != jsoniter.ArrayValue && next != jsoniter.NilValue {
		iter.ReportError("ValuesCodec.Decode", "expected JSON array")
		return
	}
	var values []*Value
	iter.ReadVal(&values)
	if iter.Error == nil {
		(*Values)(ptr).Vals = values
	}
}

func (codec *ValuesCodec) IsEmpty(ptr unsafe.Pointer) bool {
	return len((*Values)(ptr).GetVals()) == 0
}

func (codec *ValuesCodec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	stream.WriteVal((*Values)(ptr).GetVals())
}
