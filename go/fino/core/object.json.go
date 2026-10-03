package core

import (
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

func init() {
	RegisterJSONTypeDecoder(ObjectTypeFullName, &ObjectCodec{})
	RegisterJSONTypeEncoder(ObjectTypeFullName, &ObjectCodec{})
}

type ObjectCodec struct{}

func (codec *ObjectCodec) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	next := iter.WhatIsNext()
	if next != jsoniter.ObjectValue && next != jsoniter.NilValue {
		iter.ReportError("ObjectCodec.Decode", "expected JSON object")
		return
	}

	var values map[string]*Value
	iter.ReadVal(&values)
	if iter.Error == nil {
		(*Object)(ptr).Vals = values
	}
}

func (codec *ObjectCodec) IsEmpty(ptr unsafe.Pointer) bool {
	return (*Object)(ptr).IsEmpty()
}

func (codec *ObjectCodec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	stream.WriteVal((*Object)(ptr).GetVals())
}
