package core

import (
	"fmt"
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
	if iter.ReadNil() {
		if iter.Error == nil {
			if codec.isFieldPointer {
				*(**Url)(ptr) = nil
			} else {
				(*Url)(ptr).Reset()
			}
		}
		return
	}
	s := iter.ReadString()
	if iter.Error != nil {
		return
	}
	url := codec.url(ptr)
	if url == nil {
		url = &Url{}
	}

	if err := url.Parse(s); err != nil {
		iter.ReportError(UrlTypeFullName, err.Error())
		return
	}
	if codec.isFieldPointer {
		*(**Url)(ptr) = url
	}
}

func (codec *UrlStringCodec) IsEmpty(ptr unsafe.Pointer) bool {
	return codec.url(ptr) == nil
}

func (codec *UrlStringCodec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	url := codec.url(ptr)
	if url == nil {
		stream.WriteNil()
		return
	}
	stream.WriteVal(url.Format())
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
	if iter.ReadNil() {
		if iter.Error == nil {
			if codec.isFieldPointer {
				*(**Url)(ptr) = nil
			} else {
				(*Url)(ptr).Reset()
			}
		}
		return
	}
	if iter.WhatIsNext() != jsoniter.ObjectValue {
		iter.ReportError("urlStructCodec.Decode", "expected JSON object")
		return
	}
	url := codec.bareUrl(ptr)
	if url == nil {
		url = &BareUrl{}
	}
	iter.ReadVal(url)
	if iter.Error == nil && codec.isFieldPointer {
		*(**BareUrl)(ptr) = url
	}
}

func (codec *UrlStructCodec) IsEmpty(ptr unsafe.Pointer) bool {
	return codec.bareUrl(ptr) == nil
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

func (x *Url) MarshalJSON() ([]byte, error) {
	return marshalJSONCodec(x, &UrlStringCodec{})
}

func (x *Url) UnmarshalJSON(data []byte) error {
	if x == nil {
		return fmt.Errorf("url.UnmarshalJSON: nil receiver")
	}
	value, err := unmarshalJSONCodec[Url](data, &UrlStringCodec{})
	if err == nil {
		x.Scheme, x.Authority, x.Path = value.Scheme, value.Authority, value.Path
		x.Query, x.Fragment = value.Query, value.Fragment
	}
	return err
}
