package core

import jsoniter "github.com/json-iterator/go"

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
