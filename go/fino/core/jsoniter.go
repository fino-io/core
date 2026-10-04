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

// Register codecs during package initialization, before any encoding or decoding.
// jsoniter caches codecs, so runtime registration is not supported.
func RegisterJSONTypeEncoder(typ string, encoder jsoniter.ValEncoder) {
	jsoniter.RegisterTypeEncoder(typ, encoder)
}

func RegisterJSONTypeDecoder(typ string, decoder jsoniter.ValDecoder) {
	jsoniter.RegisterTypeDecoder(typ, decoder)
}

func RegisterJSONFieldEncoder(typ, field string, encoder jsoniter.ValEncoder) {
	jsoniter.RegisterFieldEncoder(typ, field, encoder)
}

func RegisterJSONFieldDecoder(typ, field string, decoder jsoniter.ValDecoder) {
	jsoniter.RegisterFieldDecoder(typ, field, decoder)
}

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
