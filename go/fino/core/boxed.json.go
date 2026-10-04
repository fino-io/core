package core

import jsoniter "github.com/json-iterator/go"

func init() {
	registerJSONFieldCodec("core.BoolValue", func(x *BoolValue) *bool { return &x.Val })
	registerJSONFieldCodec("core.BoolValues", func(x *BoolValues) *[]bool { return &x.Vals })
	registerJSONFieldCodec("core.Int32Value", func(x *Int32Value) *int32 { return &x.Val })
	registerJSONFieldCodec("core.Int64Value", func(x *Int64Value) *int64 { return &x.Val })
	registerJSONFieldCodec("core.Uint32Value", func(x *Uint32Value) *uint32 { return &x.Val })
	registerJSONFieldCodec("core.Uint64Value", func(x *Uint64Value) *uint64 { return &x.Val })
	registerJSONFieldCodec("core.Float32Value", func(x *Float32Value) *float32 { return &x.Val })
	registerJSONFieldCodec("core.Float64Value", func(x *Float64Value) *float64 { return &x.Val })
	registerJSONFieldCodec("core.StringValue", func(x *StringValue) *string { return &x.Val })
	registerJSONFieldCodec("core.BytesValue", func(x *BytesValue) *[]byte { return &x.Val })
	registerJSONFieldCodec("core.Int32Values", func(x *Int32Values) *[]int32 { return &x.Vals })
	registerJSONFieldCodec("core.Int64Values", func(x *Int64Values) *[]int64 { return &x.Vals })
	registerJSONFieldCodec("core.Uint32Values", func(x *Uint32Values) *[]uint32 { return &x.Vals })
	registerJSONFieldCodec("core.Uint64Values", func(x *Uint64Values) *[]uint64 { return &x.Vals })
	registerJSONFieldCodec("core.Float32Values", func(x *Float32Values) *[]float32 { return &x.Vals })
	registerJSONFieldCodec("core.Float64Values", func(x *Float64Values) *[]float64 { return &x.Vals })
	registerJSONFieldCodec("core.StringValues", func(x *StringValues) *[]string { return &x.Vals })
	registerJSONFieldCodec("core.StringMap", func(x *StringMap) *map[string]string { return &x.Vals })
	registerJSONFieldCodec("core.StringsMap", func(x *StringsMap) *map[string]*StringValues { return &x.Vals })
}

// Standard JSON hooks encode the field directly; jsoniter uses registered codecs.
func (x *BoolValue) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(x.Val)
}

func (x *BoolValue) UnmarshalJSON(data []byte) error {
	return decodeJSON(data, &x.Val)
}

func (x *BoolValues) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(x.Vals)
}

func (x *BoolValues) UnmarshalJSON(data []byte) error {
	return decodeJSON(data, &x.Vals)
}

func (x *Int32Value) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(x.Val)
}

func (x *Int32Value) UnmarshalJSON(data []byte) error {
	return decodeJSON(data, &x.Val)
}

func (x *Int64Value) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(x.Val)
}

func (x *Int64Value) UnmarshalJSON(data []byte) error {
	return decodeJSON(data, &x.Val)
}

func (x *Uint32Value) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(x.Val)
}

func (x *Uint32Value) UnmarshalJSON(data []byte) error {
	return decodeJSON(data, &x.Val)
}

func (x *Uint64Value) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(x.Val)
}

func (x *Uint64Value) UnmarshalJSON(data []byte) error {
	return decodeJSON(data, &x.Val)
}

func (x *Float32Value) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(x.Val)
}

func (x *Float32Value) UnmarshalJSON(data []byte) error {
	return decodeJSON(data, &x.Val)
}

func (x *Float64Value) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(x.Val)
}

func (x *Float64Value) UnmarshalJSON(data []byte) error {
	return decodeJSON(data, &x.Val)
}

func (x *StringValue) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(x.Val)
}

func (x *StringValue) UnmarshalJSON(data []byte) error {
	return decodeJSON(data, &x.Val)
}

func (x *BytesValue) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(x.Val)
}

func (x *BytesValue) UnmarshalJSON(data []byte) error {
	return decodeJSON(data, &x.Val)
}

func (x *Int32Values) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(x.Vals)
}

func (x *Int32Values) UnmarshalJSON(data []byte) error {
	return decodeJSON(data, &x.Vals)
}

func (x *Int64Values) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(x.Vals)
}

func (x *Int64Values) UnmarshalJSON(data []byte) error {
	return decodeJSON(data, &x.Vals)
}

func (x *Uint32Values) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(x.Vals)
}

func (x *Uint32Values) UnmarshalJSON(data []byte) error {
	return decodeJSON(data, &x.Vals)
}

func (x *Uint64Values) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(x.Vals)
}

func (x *Uint64Values) UnmarshalJSON(data []byte) error {
	return decodeJSON(data, &x.Vals)
}

func (x *Float32Values) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(x.Vals)
}

func (x *Float32Values) UnmarshalJSON(data []byte) error {
	return decodeJSON(data, &x.Vals)
}

func (x *Float64Values) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(x.Vals)
}

func (x *Float64Values) UnmarshalJSON(data []byte) error {
	return decodeJSON(data, &x.Vals)
}

func (x *StringValues) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(x.Vals)
}

func (x *StringValues) UnmarshalJSON(data []byte) error {
	return decodeJSON(data, &x.Vals)
}

func (x *StringMap) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(x.Vals)
}

func (x *StringMap) UnmarshalJSON(data []byte) error {
	return decodeJSON(data, &x.Vals)
}

func (x *StringsMap) MarshalJSON() ([]byte, error) {
	if x == nil {
		return []byte("null"), nil
	}
	return jsoniter.Marshal(x.Vals)
}

func (x *StringsMap) UnmarshalJSON(data []byte) error {
	return decodeJSON(data, &x.Vals)
}
