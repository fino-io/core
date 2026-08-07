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
	s := iter.ReadString()
	errorCode := codec.errorCode(ptr)
	if errorCode == nil {
		errorCode = &ErrorCode{}
		*(**ErrorCode)(ptr) = errorCode
	}

	if err := errorCode.Parse(s); err != nil {
		iter.ReportError("ErrorCodeStringCodec", err.Error())
	}
}

func (codec *ErrorCodeStringCodec) IsEmpty(ptr unsafe.Pointer) bool {
	errorCode := codec.errorCode(ptr)
	if errorCode != nil {
		if checker, ok := interface{}(errorCode).(EmptyChecker); ok {
			return checker.IsEmpty()
		}
		return false
	}
	return true
}

func (codec *ErrorCodeStringCodec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	errorCode := codec.errorCode(ptr)
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
	errorCode := codec.bareErrorCode(ptr)
	if a := iter.ReadAny(); a.ValueType() == jsoniter.ObjectValue {
		if errorCode == nil {
			errorCode = &BareErrorCode{}
			*(**BareErrorCode)(ptr) = errorCode
		}
		a.ToVal(errorCode)
	}
}

func (codec *ErrorCodeStructCodec) IsEmpty(ptr unsafe.Pointer) bool {
	errorCode := (*ErrorCode)(codec.bareErrorCode(ptr))
	if errorCode != nil {
		if checker, ok := interface{}(errorCode).(EmptyChecker); ok {
			return checker.IsEmpty()
		}
		return false
	}
	return true
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
