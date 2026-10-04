package core

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"testing"

	_ "github.com/fino-io/core/go/fino"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
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

func TestFileContractRoundTrip(t *testing.T) {
	file := &File{Name: "report.txt", Mode: File_MODE_FILE, Info: &File_Info{Suffix: "txt", Size: 7}}
	directory := &File{Name: "reports", Mode: File_MODE_DIR, Files: []*File{file}}
	encoded, err := proto.Marshal(directory)
	require.NoError(t, err)
	decoded := &File{}
	require.NoError(t, proto.Unmarshal(encoded, decoded))
	require.True(t, proto.Equal(directory, decoded))
	for _, codec := range []struct {
		name      string
		marshal   func(proto.Message) ([]byte, error)
		unmarshal func([]byte, proto.Message) error
	}{
		{"protojson", protojson.Marshal, protojson.Unmarshal},
		{"json", func(m proto.Message) ([]byte, error) { return json.Marshal(m) }, func(b []byte, m proto.Message) error { return json.Unmarshal(b, m) }},
	} {
		t.Run(codec.name, func(t *testing.T) {
			encoded, err := codec.marshal(directory)
			require.NoError(t, err)
			require.NotContains(t, string(encoded), "isDir")
			decoded := &File{}
			require.NoError(t, codec.unmarshal(encoded, decoded))
			require.True(t, proto.Equal(directory, decoded))
			var fields map[string]any
			require.NoError(t, json.Unmarshal(encoded, &fields))
			child := fields["files"].([]any)[0].(map[string]any)
			require.Equal(t, "report.txt", child["name"])
			require.NotContains(t, child["info"].(map[string]any), "name")
		})
	}
	var mode File_Mode
	require.NoError(t, mode.Parse("MODE_FILE"))
	require.Equal(t, File_MODE_FILE, mode)
	require.Equal(t, "MODE_FILE", mode.Format())
}
