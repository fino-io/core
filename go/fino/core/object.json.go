package core

import (
	"fmt"
	"unicode/utf8"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

func init() {
	RegisterJSONTypeDecoder(ObjectTypeFullName, &ObjectCodec{})
	RegisterJSONTypeEncoder(ObjectTypeFullName, &ObjectCodec{})
}

type ObjectCodec struct{}

func validateJSONKeys(values map[string]*Value) error {
	for key := range values {
		if !utf8.ValidString(key) {
			return fmt.Errorf("invalid UTF-8 in object key: %q", key)
		}
	}
	return nil
}

func readJSONObject(iter *jsoniter.Iterator) map[string]*Value {
	var values map[string]*Value
	iter.ReadVal(&values)
	if iter.Error == nil {
		if err := validateJSONKeys(values); err != nil {
			iter.ReportError("readJSONObject", err.Error())
		}
	}
	return values
}

func (codec *ObjectCodec) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	next := iter.WhatIsNext()
	if next != jsoniter.ObjectValue && next != jsoniter.NilValue {
		iter.ReportError("objectCodec.Decode", "expected JSON object")
		return
	}

	values := readJSONObject(iter)
	if iter.Error == nil {
		(*Object)(ptr).Vals = values
	}
}

func (codec *ObjectCodec) IsEmpty(ptr unsafe.Pointer) bool {
	return (*Object)(ptr).IsEmpty()
}

func (codec *ObjectCodec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	values := (*Object)(ptr).GetVals()
	if err := validateJSONKeys(values); err != nil {
		stream.Error = err
		return
	}
	stream.WriteVal(values)
}

func (x *Object) MarshalJSON() ([]byte, error) {
	if err := x.CheckValid(); err != nil {
		return nil, err
	}
	return marshalJSONCodec(x, &ObjectCodec{})
}

func (x *Object) UnmarshalJSON(data []byte) error {
	if x == nil {
		return fmt.Errorf("object.UnmarshalJSON: nil receiver")
	}
	value, err := unmarshalJSONCodec[Object](data, &ObjectCodec{})
	if err == nil {
		x.Vals = value.Vals
	}
	return err
}
