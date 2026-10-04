package core

import (
	"fmt"
	"maps"
	"reflect"
	"unicode/utf8"

	"github.com/fino-io/core/go/fino/core/strcase"
	jsoniter "github.com/json-iterator/go"
	"google.golang.org/protobuf/proto"
)

const ObjectTypeName = "Object"
const ObjectTypeFullName = "core.Object"

func NewObject() *Object {
	return &Object{Vals: make(map[string]*Value)}
}

// NewObjectFromMap constructs an Object from a general-purpose Go map.
// The map keys must be valid UTF-8.
// The map values are converted using NewValue.
func NewObjectFromMap(m map[string]any) (*Object, error) {
	x := &Object{Vals: make(map[string]*Value, len(m))}
	for k, v := range m {
		if !utf8.ValidString(k) {
			return nil, fmt.Errorf("invalid UTF-8 in object key: %q", k)
		}
		var err error
		x.Vals[k], err = NewValue(v)
		if err != nil {
			return nil, err
		}
	}
	return x, nil
}

func NewObjectFromKeyVals(kvs ...any) (*Object, error) {
	if len(kvs)%2 != 0 {
		return nil, fmt.Errorf("invalid number of key/value pairs: %d", len(kvs))
	}

	m := make(map[string]any, len(kvs)/2)
	for i := 0; i < len(kvs); i += 2 {
		k, ok := kvs[i].(string)
		if !ok {
			return nil, fmt.Errorf("invalid key/value pair at index %d", i)
		}
		m[k] = kvs[i+1]
	}

	return NewObjectFromMap(m)
}

func NewObjectFrom(val any) (*Object, error) {
	obj := NewObject()
	return obj, obj.From(val)
}

func NewObjectFromValues(v map[string]*Value) *Object {
	return &Object{Vals: v}
}

func MergeObjects(objs ...*Object) *Object {
	obj := NewObject()
	for _, o := range objs {
		obj.Merge(o)
	}
	return obj
}

// AsMap converts Object into a map[string]any using Value.AsInterface.
// The result is intended for JSON serialization or dynamic inspection,
// not for round-trip binary fidelity.
func (x *Object) AsMap() map[string]any {
	if x == nil || x.Vals == nil {
		return nil
	}

	f := x.GetVals()
	vs := make(map[string]any, len(f))
	for k, v := range f {
		vs[k] = v.AsInterface()
	}

	return vs
}

func (x *Object) To(val any) error {
	if x == nil {
		return nil
	}
	rv := reflect.ValueOf(val)
	if !rv.IsValid() || rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("Object.To: expected a non-nil pointer, got %T", val)
	}

	marshal, err := jsoniter.ConfigFastest.Marshal(x)
	if err != nil {
		return err
	}

	return jsoniter.ConfigFastest.Unmarshal(marshal, val)
}

// From replaces the object with the input's JSON representation, preserving
// JSON field names, omitempty, and registered codecs. JSON null clears the object;
// other non-object representations return an error without changing the receiver.
// Use Merge to combine objects and Clone to preserve Protobuf types and metadata.
func (x *Object) From(val any) error {
	if x == nil {
		return nil
	}

	data, err := jsoniter.ConfigFastest.Marshal(val)
	if err != nil {
		return err
	}

	var object Object
	if err := jsoniter.ConfigFastest.Unmarshal(data, &object); err != nil {
		return err
	}
	x.Vals = object.Vals
	return nil
}

// From2 is a compatibility alias for From.
// Deprecated: use From instead.
func (x *Object) From2(val any) error {
	return x.From(val)
}

func (x *Object) IsEmpty() bool {
	return len(x.GetVals()) == 0
}

func (x *Object) init() {
	if x != nil && x.Vals == nil {
		x.Vals = make(map[string]*Value)
	}
}

func (x *Object) SetValue(key string, val *Value) *Object {
	if x != nil {
		x.init()
		x.Vals[key] = val
	}
	return x
}

func (x *Object) SetBool(key string, val bool) *Object {
	return x.SetValue(key, NewBoolValue(val))
}

func (x *Object) SetBytes(key string, val []byte) *Object {
	return x.SetValue(key, NewBytesValue(val))
}

func (x *Object) SetInt(key string, val int) *Object {
	return x.SetValue(key, NewIntValue(val))
}

func (x *Object) SetInt32(key string, val int32) *Object {
	return x.SetValue(key, NewInt32Value(val))
}

func (x *Object) SetInt64(key string, val int64) *Object {
	return x.SetValue(key, NewInt64Value(val))
}

func (x *Object) SetUint(key string, val uint) *Object {
	return x.SetValue(key, NewUintValue(val))
}

func (x *Object) SetUint32(key string, val uint32) *Object {
	return x.SetValue(key, NewUint32Value(val))
}

func (x *Object) SetUint64(key string, val uint64) *Object {
	return x.SetValue(key, NewUint64Value(val))
}

func (x *Object) SetFloat32(key string, val float32) *Object {
	return x.SetValue(key, NewFloat32Value(val))
}

func (x *Object) SetFloat64(key string, val float64) *Object {
	return x.SetValue(key, NewFloat64Value(val))
}

func (x *Object) SetString(key string, val string) *Object {
	return x.SetValue(key, NewStringValue(val))
}

func (x *Object) SetObject(key string, val *Object) *Object {
	return x.SetValue(key, NewObjectValue(val))
}

func (x *Object) SetIntArray(key string, vals ...int) *Object {
	return x.SetValue(key, NewIntArrayValue(vals...))
}

func (x *Object) SetInt32Array(key string, vals ...int32) *Object {
	return x.SetValue(key, NewInt32ArrayValue(vals...))
}

func (x *Object) SetInt64Array(key string, vals ...int64) *Object {
	return x.SetValue(key, NewInt64ArrayValue(vals...))
}

func (x *Object) SetUintArray(key string, vals ...uint) *Object {
	return x.SetValue(key, NewUintArrayValue(vals...))
}

func (x *Object) SetUint32Array(key string, vals ...uint32) *Object {
	return x.SetValue(key, NewUint32ArrayValue(vals...))
}

func (x *Object) SetUint64Array(key string, vals ...uint64) *Object {
	return x.SetValue(key, NewUint64ArrayValue(vals...))
}

func (x *Object) SetFloat32Array(key string, vals ...float32) *Object {
	return x.SetValue(key, NewFloat32ArrayValue(vals...))
}

func (x *Object) SetFloat64Array(key string, vals ...float64) *Object {
	return x.SetValue(key, NewFloat64ArrayValue(vals...))
}

func (x *Object) SetStringArray(key string, vals ...string) *Object {
	return x.SetValue(key, NewStringArrayValue(vals...))
}

func (x *Object) SetObjectArray(key string, vals ...*Object) *Object {
	return x.SetValue(key, NewObjectArrayValue(vals...))
}

func (x *Object) GetValue(key string) *Value {
	return x.GetVals()[key]
}

func (x *Object) GetBool(key string) bool {
	return x.GetValue(key).GetBool()
}

func (x *Object) GetBytes(key string) []byte {
	return x.GetValue(key).GetBytes()
}

func (x *Object) GetInt(key string) int {
	return x.GetValue(key).GetInt()
}

func (x *Object) GetInt32(key string) int32 {
	return x.GetValue(key).GetInt32()
}

func (x *Object) GetInt64(key string) int64 {
	return x.GetValue(key).GetInt64()
}

func (x *Object) GetUint(key string) uint {
	return x.GetValue(key).GetUint()
}

func (x *Object) GetUint32(key string) uint32 {
	return x.GetValue(key).GetUint32()
}

func (x *Object) GetUint64(key string) uint64 {
	return x.GetValue(key).GetUint64()
}

func (x *Object) GetFloat32(key string) float32 {
	return x.GetValue(key).GetFloat32()
}

func (x *Object) GetFloat64(key string) float64 {
	return x.GetValue(key).GetFloat64()
}

func (x *Object) GetString(key string) string {
	return x.GetValue(key).GetString()
}

func (x *Object) GetObject(key string) *Object {
	return x.GetValue(key).GetObject()
}

func (x *Object) GetBoolArray(key string) []bool {
	return x.GetValue(key).GetBoolArray()
}

func (x *Object) GetIntArray(key string) []int {
	return x.GetValue(key).GetIntArray()
}

func (x *Object) GetInt64Array(key string) []int64 {
	return x.GetValue(key).GetInt64Array()
}

func (x *Object) GetUintArray(key string) []uint {
	return x.GetValue(key).GetUintArray()
}

func (x *Object) GetFloat64Array(key string) []float64 {
	return x.GetValue(key).GetFloat64Array()
}

func (x *Object) GetStringArray(key string) []string {
	return x.GetValue(key).GetStringArray()
}

func (x *Object) GetObjectArray(key string) []*Object {
	return x.GetValue(key).GetObjectArray()
}

func (x *Object) GetValueArray(key string) []*Value {
	return x.GetValue(key).GetValueArray()
}

// Merge copies entries from o. Nested values remain shared; use Clone to isolate them.
func (x *Object) Merge(o *Object) *Object {
	if x != nil && o != nil {
		x.init()
		maps.Copy(x.Vals, o.Vals)
	}
	return x
}

// Clone returns a deep copy, including nested objects, arrays, and byte slices.
func (x *Object) Clone() *Object {
	if x != nil {
		return proto.Clone(x).(*Object)
	}
	return x
}

func (x *Object) Delete(key string) *Object {
	delete(x.GetVals(), key)
	return x
}

func (x *Object) ToLowerCamelKeys() *Object {
	obj := NewObject()
	for k, v := range x.GetVals() {
		obj.SetValue(strcase.ToLowerCamel(k), v)
	}
	return obj
}

func (x *Object) ToSnakeKeys() *Object {
	obj := NewObject()
	for k, v := range x.GetVals() {
		obj.SetValue(strcase.ToSnake(k), v)
	}
	return obj
}
