package core

import (
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

func init() {
	RegisterJSONTypeDecoder(UrlTypeFullName, &UrlStringCodec{})
	RegisterJSONTypeEncoder(UrlTypeFullName, &UrlStringCodec{})
}

type UrlStringCodec struct {
	isFieldPointer bool
}

func (codec *UrlStringCodec) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	url := codec.url(ptr)
	if url == nil {
		url = &Url{}
		*(**Url)(ptr) = url
	}

	if err := url.Parse(iter.ReadString()); err != nil {
		iter.ReportError(UrlTypeFullName, err.Error())
	}
}

func (codec *UrlStringCodec) IsEmpty(ptr unsafe.Pointer) bool {
	url := codec.url(ptr)
	if url != nil {
		if checker, ok := interface{}(url).(EmptyChecker); ok {
			return checker.IsEmpty()
		}
		return false
	}
	return true
}

func (codec *UrlStringCodec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	url := codec.url(ptr)
	stream.WriteString(url.Format())
}

func (codec *UrlStringCodec) url(ptr unsafe.Pointer) *Url {
	if codec.isFieldPointer {
		return *(**Url)(ptr)
	}
	return (*Url)(ptr)
}

// BareUrl will be jsonify to raw, without any codec
type BareUrl Url

type UrlStructCodec struct {
	isFieldPointer bool
}

func (codec *UrlStructCodec) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	url := codec.bareUrl(ptr)
	a := iter.ReadAny()
	if a.ValueType() == jsoniter.ObjectValue {
		if url == nil {
			url = &BareUrl{}
			*(**BareUrl)(ptr) = url
		}
		a.ToVal(url)
	}
}

func (codec *UrlStructCodec) IsEmpty(ptr unsafe.Pointer) bool {
	url := (*Url)(codec.bareUrl(ptr))
	if url != nil {
		if checker, ok := interface{}(url).(EmptyChecker); ok {
			return checker.IsEmpty()
		}
		return false
	}
	return true
}

func (codec *UrlStructCodec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	stream.WriteVal(codec.bareUrl(ptr))
}

func (codec *UrlStructCodec) bareUrl(ptr unsafe.Pointer) *BareUrl {
	if codec.isFieldPointer {
		return *(**BareUrl)(ptr)
	}
	return (*BareUrl)(ptr)
}
