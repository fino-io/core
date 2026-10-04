package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

// readJSON validates the whole document and preserves numbers in dynamic data.
func readJSON(data []byte, destination any) error {
	if !json.Valid(data) {
		return errors.New("invalid JSON")
	}
	decoder := jsoniter.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	return decoder.Decode(destination)
}

// Standard JSON hooks replace a field only after the whole input succeeds.
func decodeJSON[T any](data []byte, destination *T) error {
	var value T
	if err := jsoniter.Unmarshal(data, &value); err != nil {
		return err
	}
	*destination = value
	return nil
}

// Standard JSON hooks reuse the wire codecs without depending on jsoniter's
// choice between a registered codec and a MarshalJSON method.
func marshalJSONCodec[T any](value *T, codec jsoniter.ValEncoder) ([]byte, error) {
	if value == nil {
		return []byte("null"), nil
	}
	stream := jsoniter.NewStream(jsoniter.ConfigDefault, nil, 256)
	codec.Encode(unsafe.Pointer(value), stream)
	return stream.Buffer(), stream.Error
}

func unmarshalJSONCodec[T any](data []byte, codec jsoniter.ValDecoder) (*T, error) {
	if !json.Valid(data) {
		return nil, errors.New("invalid JSON")
	}
	value := new(T)
	iter := jsoniter.ParseBytes(jsoniter.ConfigDefault, data)
	codec.Decode(unsafe.Pointer(value), iter)
	if iter.Error != nil && iter.Error != io.EOF {
		return nil, iter.Error
	}
	return value, nil
}

// RegisterJSONTypeEncoder registers a type encoder during package initialization.
// jsoniter caches codecs, so runtime registration is not supported.
func RegisterJSONTypeEncoder(typ string, encoder jsoniter.ValEncoder) {
	jsoniter.RegisterTypeEncoder(typ, encoder)
}

// RegisterJSONTypeDecoder registers a type decoder during package initialization.
func RegisterJSONTypeDecoder(typ string, decoder jsoniter.ValDecoder) {
	jsoniter.RegisterTypeDecoder(typ, decoder)
}

// RegisterJSONFieldEncoder registers a field encoder during package initialization.
func RegisterJSONFieldEncoder(typ, field string, encoder jsoniter.ValEncoder) {
	jsoniter.RegisterFieldEncoder(typ, field, encoder)
}

// RegisterJSONFieldDecoder registers a field decoder during package initialization.
func RegisterJSONFieldDecoder(typ, field string, decoder jsoniter.ValDecoder) {
	jsoniter.RegisterFieldDecoder(typ, field, decoder)
}

// RegisterJSONValuesCodec registers a wrapper's slice field as its JSON value.
func RegisterJSONValuesCodec[T any, M any](typ string, field func(*M) *[]T) {
	registerJSONFieldCodec(typ, field)
}

// Wrapper messages serialize their field on the caller's iterator or stream.
func registerJSONFieldCodec[T any, M any](typ string, field func(*M) *T) {
	codec := &jsonFieldCodec[T, M]{field: field}
	RegisterJSONTypeDecoder(typ, codec)
	RegisterJSONTypeEncoder(typ, codec)
}

type jsonFieldCodec[T any, M any] struct {
	field func(*M) *T
}

func (codec *jsonFieldCodec[T, M]) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	destination := codec.field((*M)(ptr))
	if destination == nil {
		iter.ReportError("jsonFieldCodec.Decode", "nil field destination")
		return
	}
	// Success replaces the field; decoding errors leave it intact.
	var value T
	iter.ReadVal(&value)
	if iter.Error == nil || iter.Error == io.EOF {
		*destination = value
	}
}

func (codec *jsonFieldCodec[T, M]) IsEmpty(ptr unsafe.Pointer) bool {
	if ptr == nil {
		return true
	}
	field := codec.field((*M)(ptr))
	if field == nil {
		return true
	}
	value := reflect.ValueOf(*field)
	switch value.Kind() {
	case reflect.Slice, reflect.Map:
		return value.Len() == 0
	default:
		return value.IsZero()
	}
}

func (codec *jsonFieldCodec[T, M]) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	if ptr == nil {
		stream.WriteNil()
		return
	}
	stream.WriteVal(codec.field((*M)(ptr)))
}
