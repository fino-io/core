package core

import (
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

func init() {
	RegisterJSONFieldEncoder("core.Error", "Code", &ErrorCodeStringCodec{IsFieldPointer: true})
	RegisterJSONFieldDecoder("core.Error", "Code", &ErrorCodeStringCodec{IsFieldPointer: true})
}

// BareErrorCode will be jsonify to raw, without any codec
type BareErrorCode ErrorCode

type ErrorCodeStringCodec struct {
	IsFieldPointer bool
}

func (codec *ErrorCodeStringCodec) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	if iter.ReadNil() {
		if iter.Error == nil {
			if codec.IsFieldPointer {
				*(**ErrorCode)(ptr) = nil
			} else {
				(*ErrorCode)(ptr).Reset()
			}
		}
		return
	}
	s := iter.ReadString()
	if iter.Error != nil {
		return
	}
	errorCode := codec.errorCode(ptr)
	if errorCode == nil {
		errorCode = &ErrorCode{}
	}

	if err := errorCode.Parse(s); err != nil {
		iter.ReportError("ErrorCodeStringCodec", err.Error())
		return
	}
	if codec.IsFieldPointer {
		*(**ErrorCode)(ptr) = errorCode
	}
}

func (codec *ErrorCodeStringCodec) IsEmpty(ptr unsafe.Pointer) bool {
	return codec.errorCode(ptr) == nil
}

func (codec *ErrorCodeStringCodec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	errorCode := codec.errorCode(ptr)
	if errorCode == nil {
		stream.WriteNil()
		return
	}
	stream.WriteString(errorCode.Format())
}

func (codec *ErrorCodeStringCodec) errorCode(ptr unsafe.Pointer) *ErrorCode {
	if codec.IsFieldPointer {
		return *(**ErrorCode)(ptr)
	}
	return (*ErrorCode)(ptr)
}

type ErrorCodeStructCodec struct {
	IsFieldPointer bool
}

func (codec *ErrorCodeStructCodec) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	if iter.ReadNil() {
		if iter.Error == nil {
			if codec.IsFieldPointer {
				*(**ErrorCode)(ptr) = nil
			} else {
				(*ErrorCode)(ptr).Reset()
			}
		}
		return
	}
	if iter.WhatIsNext() != jsoniter.ObjectValue {
		iter.ReportError("ErrorCodeStructCodec.Decode", "expected JSON object")
		return
	}
	errorCode := codec.bareErrorCode(ptr)
	if errorCode == nil {
		errorCode = &BareErrorCode{}
	}
	iter.ReadVal(errorCode)
	if iter.Error == nil && codec.IsFieldPointer {
		*(**BareErrorCode)(ptr) = errorCode
	}
}

func (codec *ErrorCodeStructCodec) IsEmpty(ptr unsafe.Pointer) bool {
	return codec.bareErrorCode(ptr) == nil
}

func (codec *ErrorCodeStructCodec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	stream.WriteVal(codec.bareErrorCode(ptr))
}

func (codec *ErrorCodeStructCodec) bareErrorCode(ptr unsafe.Pointer) *BareErrorCode {
	if codec.IsFieldPointer {
		return *(**BareErrorCode)(ptr)
	}
	return (*BareErrorCode)(ptr)
}
