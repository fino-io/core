package core

import jsoniter "github.com/json-iterator/go"

const ValuesTypeName = "Values"
const ValuesTypeFullName = "core.Values"

// NewValues constructs a ListValue from a general-purpose Go slice.
// The slice elements are converted using NewValue.
func NewValues(v []any) (*Values, error) {
	if v == nil {
		return &Values{}, nil
	}
	x := &Values{Vals: make([]*Value, len(v))}
	for i, v := range v {
		var err error
		x.Vals[i], err = NewValue(v)
		if err != nil {
			return nil, err
		}
	}
	return x, nil
}

// AsSlice converts x to a general-purpose Go slice.
// The slice elements are converted by calling Value.AsInterface.
func (x *Values) AsSlice() []any {
	if x == nil || x.Vals == nil {
		return nil
	}
	vals := x.GetVals()
	vs := make([]any, len(vals))
	for i, v := range vals {
		vs[i] = v.AsInterface()
	}
	return vs
}

func (x *Values) MarshalJSON() ([]byte, error) {
	if x == nil {
		return []byte("null"), nil
	}
	return jsoniter.Marshal(x.Vals)
}

func (x *Values) UnmarshalJSON(b []byte) error {
	return decodeJSON(b, &x.Vals)
}
