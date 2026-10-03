package core

import (
	"encoding/hex"
	"reflect"
	"strconv"
)

type StringLike interface {
	ToString() string
}

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
	switch v := value.(type) {
	case bool:
		return strconv.FormatBool(v)
	case []byte:
		return hex.EncodeToString(v)
	case string:
		return v
	case StringLike:
		return v.ToString()
	case Formatter:
		return v.Format()
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'g', -1, v.Type().Bits())
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.String:
		return v.String()
	}
	return ""
}
