package core

type FromBoolConverter interface {
	FromBool(value bool) error
}

type ToBoolConverter interface {
	ToBool() bool
}

type FromIntConverter interface {
	FromInt(value int) error
}

type ToIntConverter interface {
	ToInt() int
}

type FromInt32Converter interface {
	FromInt32(value int32) error
}

type ToInt32Converter interface {
	ToInt32() int32
}

type FromInt64Converter interface {
	FromInt64(value int64) error
}

type ToInt64Converter interface {
	ToInt64() int64
}

type FromUintConverter interface {
	FromUint(value uint) error
}

type ToUintConverter interface {
	ToUint() uint
}

type FromUint32Converter interface {
	FromUint32(value uint32) error
}

type ToUint32Converter interface {
	ToUint32() uint32
}

type FromUint64Converter interface {
	FromUint64(value uint64) error
}

type ToUint64Converter interface {
	ToUint64() uint64
}

type FromFloat32Converter interface {
	FromFloat32(value float32) error
}

type ToFloat32Converter interface {
	ToFloat32() float32
}

type FromFloat64Converter interface {
	FromFloat64(value float64) error
}

type ToFloat64Converter interface {
	ToFloat64() float64
}

type FromStringConverter interface {
	FromString(value string) error
}

type ToStringConverter interface {
	ToString() string
}
