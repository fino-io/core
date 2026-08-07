package core

import (
	"sync"

	jsoniter "github.com/json-iterator/go"
)

var (
	registerJSONEncoderTypes     map[string]jsoniter.ValEncoder
	registerJSONEncoderTypesOnce = &sync.Once{}

	registerJSONEncoderTypeFields     map[string]jsoniter.ValEncoder
	registerJSONEncoderTypeFieldsOnce = &sync.Once{}
)

func RegisterJSONTypeEncoder(typ string, encoder jsoniter.ValEncoder) {
	jsoniter.RegisterTypeEncoder(typ, encoder)

	registerJSONEncoderTypesOnce.Do(func() {
		registerJSONEncoderTypes = make(map[string]jsoniter.ValEncoder)
	})
	registerJSONEncoderTypes[typ] = encoder
}

func RegisterJSONTypeDecoder(typ string, decoder jsoniter.ValDecoder) {
	jsoniter.RegisterTypeDecoder(typ, decoder)
}

func RegisterJSONFieldEncoder(typ, field string, encoder jsoniter.ValEncoder) {
	jsoniter.RegisterFieldEncoder(typ, field, encoder)

	registerJSONEncoderTypeFieldsOnce.Do(func() {
		registerJSONEncoderTypeFields = make(map[string]jsoniter.ValEncoder)
	})
	registerJSONEncoderTypeFields[typ+"."+field] = encoder
}

func RegisterJSONFieldDecoder(typ, field string, decoder jsoniter.ValDecoder) {
	jsoniter.RegisterFieldDecoder(typ, field, decoder)
}
