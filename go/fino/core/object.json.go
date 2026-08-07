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
	obj := (*Object)(ptr)
	if iter.WhatIsNext() == jsoniter.NilValue {
		iter.ReadNil()
		obj.Vals = nil
		return
	}

	a := iter.ReadAny()
	if a.ValueType() != jsoniter.ObjectValue {
		iter.ReportError("ObjectCodec.Decode", "expected JSON object")
		return
	}

	obj.Vals = make(map[string]*Value, a.Size())
	for _, k := range a.Keys() {
		val, err := NewValueCodec().DecodeAny(a.Get(k))
		if err != nil {
			iter.ReportError("ObjectCodec.Decode", err.Error())
		}
		obj.Vals[k] = val
	}
}

func (codec *ObjectCodec) IsEmpty(ptr unsafe.Pointer) bool {
	obj := (*Object)(ptr)
	return obj == nil || len(obj.Vals) == 0
}

func (codec *ObjectCodec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	obj := (*Object)(ptr)
	if obj == nil || obj.Vals == nil {
		stream.WriteNil()
		return
	}
	stream.WriteVal(obj.Vals)
}
