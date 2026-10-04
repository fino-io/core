package core

import (
	"fmt"
	"reflect"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// CheckValid checks the standard Protobuf timestamp range (years 1–9999).
func (x *Timestamp) CheckValid() error {
	if x == nil {
		return fmt.Errorf("timestamp: nil value")
	}
	return (&timestamppb.Timestamp{Seconds: x.Seconds, Nanos: x.Nanoseconds}).CheckValid()
}

// CheckValid checks canonical fractional seconds while retaining Core's full
// int64 seconds range, beyond the standard Protobuf duration range.
func (x *Duration) CheckValid() error {
	if x == nil {
		return fmt.Errorf("duration: nil value")
	}
	if x.Nanoseconds <= -1000000000 || x.Nanoseconds >= 1000000000 {
		return fmt.Errorf("duration: nanoseconds out of range: %d", x.Nanoseconds)
	}
	if x.Seconds > 0 && x.Nanoseconds < 0 || x.Seconds < 0 && x.Nanoseconds > 0 {
		return fmt.Errorf("duration: seconds and nanoseconds have different signs")
	}
	return nil
}

func (x *Value) checkScalar() error {
	if x == nil || x.Val == nil {
		return nil
	}
	if reflect.ValueOf(x.Val).IsNil() {
		return fmt.Errorf("value: nil oneof payload")
	}
	switch v := x.Val.(type) {
	case *Value_NegativeValue:
		if v.NegativeValue == 0 || v.NegativeValue > uint64(1)<<63 {
			return fmt.Errorf("value: negative magnitude out of int64 range: %d", v.NegativeValue)
		}
	case *Value_StringValue:
		if !utf8.ValidString(v.StringValue) {
			return fmt.Errorf("invalid UTF-8 in string: %q", v.StringValue)
		}
	}
	return nil
}

// CheckValid checks nested keys, payloads, and cycles. Nil values and absent
// payloads represent null.
func (x *Value) CheckValid() error {
	return checkValue(x, make(map[*Value]bool))
}

func (x *Object) CheckValid() error {
	return checkObject(x, make(map[*Value]bool))
}

func (x *Values) CheckValid() error {
	return checkValues(x, make(map[*Value]bool))
}

func checkValue(x *Value, visiting map[*Value]bool) error {
	if err := x.checkScalar(); err != nil || x == nil {
		return err
	}
	if visiting[x] {
		return fmt.Errorf("cyclic value")
	}
	visiting[x] = true
	defer delete(visiting, x)
	if err := checkObject(x.GetObject(), visiting); err != nil {
		return err
	}
	return checkValues(x.GetValuesValue(), visiting)
}

func checkObject(x *Object, visiting map[*Value]bool) error {
	if err := validateJSONKeys(x.GetVals()); err != nil {
		return err
	}
	for key, value := range x.GetVals() {
		if err := checkValue(value, visiting); err != nil {
			return fmt.Errorf("object field %q: %w", key, err)
		}
	}
	return nil
}

func checkValues(x *Values, visiting map[*Value]bool) error {
	for i, value := range x.GetVals() {
		if err := checkValue(value, visiting); err != nil {
			return fmt.Errorf("array element %d: %w", i, err)
		}
	}
	return nil
}
