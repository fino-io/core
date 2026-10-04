package core

import (
	"encoding/hex"
	"reflect"
	"strconv"
)

type ScalarLike interface {
	ToScalar() any
}

type ArrayLike interface {
	ToArray() any
}

type MapLike interface {
	ToMap() any
}

func ToString(value any) string {
	if bytes, ok := value.([]byte); ok {
		return hex.EncodeToString(bytes)
	}
	text, _ := formatScalar(value)
	return text
}

func formatScalar(value any) (string, bool) {
	v := reflect.ValueOf(value)
	if !v.IsValid() || (v.Kind() == reflect.Ptr && v.IsNil()) {
		return "", false
	}
	switch v := value.(type) {
	case ToStringConverter:
		return v.ToString(), true
	case Formatter:
		return v.Format(), true
	}

	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10), true
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'g', -1, v.Type().Bits()), true
	case reflect.Bool:
		return strconv.FormatBool(v.Bool()), true
	case reflect.String:
		return v.String(), true
	}
	return "", false
}
