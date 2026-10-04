package core

import (
	"reflect"
	"strconv"
	"strings"
	"testing"

	_ "github.com/fino-io/core/go/fino"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// Fino adapts protoc's snake_case struct tags to the contract's JSON names.
func TestGeneratedJSONTagsMatchDescriptors(t *testing.T) {
	checked := 0
	protoregistry.GlobalTypes.RangeMessages(func(message protoreflect.MessageType) bool {
		name := string(message.Descriptor().FullName())
		if !strings.HasPrefix(name, "fino.") || message.Descriptor().IsMapEntry() {
			return true
		}
		typ := reflect.TypeOf(message.New().Interface()).Elem()
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			tag := strings.Split(field.Tag.Get("protobuf"), ",")
			if len(tag) < 2 {
				continue
			}
			number, err := strconv.ParseInt(tag[1], 10, 32)
			require.NoError(t, err)
			descriptor := message.Descriptor().Fields().ByNumber(protoreflect.FieldNumber(number))
			require.NotNil(t, descriptor)
			require.Equal(t, descriptor.JSONName()+",omitempty", field.Tag.Get("json"), "%s.%s", name, field.Name)
			checked++
		}
		return true
	})
	require.Greater(t, checked, 0)
}
