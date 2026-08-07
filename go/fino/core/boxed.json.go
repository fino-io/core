package core

import jsoniter "github.com/json-iterator/go"

func (x *BoolValue) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(x.Val)
}

func (x *BoolValue) UnmarshalJSON(data []byte) error {
	return jsoniter.Unmarshal(data, &x.Val)
}

func (x *BoolValues) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(x.Vals)
}

func (x *BoolValues) UnmarshalJSON(data []byte) error {
	return jsoniter.Unmarshal(data, &x.Vals)
}

func (x *Int32Value) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(x.Val)
}

func (x *Int32Value) UnmarshalJSON(data []byte) error {
	return jsoniter.Unmarshal(data, &x.Val)
}

func (x *Int64Value) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(x.Val)
}

func (x *Int64Value) UnmarshalJSON(data []byte) error {
	return jsoniter.Unmarshal(data, &x.Val)
}

func (x *Uint32Value) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(x.Val)
}

func (x *Uint32Value) UnmarshalJSON(data []byte) error {
	return jsoniter.Unmarshal(data, &x.Val)
}

func (x *Uint64Value) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(x.Val)
}

func (x *Uint64Value) UnmarshalJSON(data []byte) error {
	return jsoniter.Unmarshal(data, &x.Val)
}

func (x *Float32Value) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(x.Val)
}

func (x *Float32Value) UnmarshalJSON(data []byte) error {
	return jsoniter.Unmarshal(data, &x.Val)
}

func (x *Float64Value) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(x.Val)
}

func (x *Float64Value) UnmarshalJSON(data []byte) error {
	return jsoniter.Unmarshal(data, &x.Val)
}

func (x *StringValue) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(x.Val)
}

func (x *StringValue) UnmarshalJSON(data []byte) error {
	return jsoniter.Unmarshal(data, &x.Val)
}

func (x *BytesValue) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(x.Val)
}

func (x *BytesValue) UnmarshalJSON(data []byte) error {
	return jsoniter.Unmarshal(data, &x.Val)
}

func (x *Int32Values) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(x.Vals)
}

func (x *Int32Values) UnmarshalJSON(data []byte) error {
	return jsoniter.Unmarshal(data, &x.Vals)
}

func (x *Int64Values) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(x.Vals)
}

func (x *Int64Values) UnmarshalJSON(data []byte) error {
	return jsoniter.Unmarshal(data, &x.Vals)
}

func (x *Uint32Values) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(x.Vals)
}

func (x *Uint32Values) UnmarshalJSON(data []byte) error {
	return jsoniter.Unmarshal(data, &x.Vals)
}

func (x *Uint64Values) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(x.Vals)
}

func (x *Uint64Values) UnmarshalJSON(data []byte) error {
	return jsoniter.Unmarshal(data, &x.Vals)
}

func (x *Float32Values) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(x.Vals)
}

func (x *Float32Values) UnmarshalJSON(data []byte) error {
	return jsoniter.Unmarshal(data, &x.Vals)
}

func (x *Float64Values) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(x.Vals)
}

func (x *Float64Values) UnmarshalJSON(data []byte) error {
	return jsoniter.Unmarshal(data, &x.Vals)
}

func (x *StringValues) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(x.Vals)
}

func (x *StringValues) UnmarshalJSON(data []byte) error {
	return jsoniter.Unmarshal(data, &x.Vals)
}

func (x *StringMap) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(x.Vals)
}

func (x *StringMap) UnmarshalJSON(data []byte) error {
	return jsoniter.Unmarshal(data, &x.Vals)
}

func (x *StringsMap) MarshalJSON() ([]byte, error) {
	if x == nil {
		return []byte("null"), nil
	}

	tmp := make(map[string]*StringValues, len(x.Vals))
	for k, v := range x.Vals {
		if v != nil {
			tmp[k] = v
		}
	}

	return jsoniter.Marshal(tmp)
}

func (x *StringsMap) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	tmp := make(map[string]*StringValues)
	if err := jsoniter.Unmarshal(data, &tmp); err != nil {
		return err
	}

	if x.Vals == nil {
		x.Vals = make(map[string]*StringValues, len(tmp))
	}

	for k, v := range tmp {
		x.Vals[k] = v
	}

	return nil
}
