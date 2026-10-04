package core

func init() {
	RegisterJSONValuesCodec(ValuesTypeFullName, func(x *Values) *[]*Value { return &x.Vals })
}
