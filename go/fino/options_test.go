package fino_test

import (
	"encoding/json"
	"testing"

	"github.com/fino-io/core/go/fino"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestOptionsPreservePresence(t *testing.T) {
	for _, test := range []struct {
		name    string
		message proto.Message
	}{
		{"database omitted", &fino.DBOptions{}},
		{"database explicit empty", &fino.DBOptions{
			Len: proto.Int64(0), Default: proto.String(""), Comment: proto.String(""),
			Index: proto.String(""), UniqueIndex: proto.String(""),
		}},
		{"database configured", &fino.DBOptions{Len: proto.Int64(128), Default: proto.String("'name'"), Index: proto.String("by_name")}},
		{"validation omitted", &fino.ValidateOptions{}},
		{"validation explicit zero", &fino.ValidateOptions{
			Len: proto.Int64(0), Min: proto.Int64(0), Max: proto.Int64(0), Oneof: proto.String(""),
		}},
		{"validation configured", &fino.ValidateOptions{Min: proto.Int64(-1), Max: proto.Int64(10), Oneof: proto.String("active disabled")}},
	} {
		t.Run(test.name, func(t *testing.T) {
			for _, codec := range []struct {
				name      string
				marshal   func(proto.Message) ([]byte, error)
				unmarshal func([]byte, proto.Message) error
			}{
				{"protobuf", proto.Marshal, proto.Unmarshal},
				{"protojson", protojson.Marshal, protojson.Unmarshal},
				{"json", func(m proto.Message) ([]byte, error) { return json.Marshal(m) }, func(b []byte, m proto.Message) error { return json.Unmarshal(b, m) }},
			} {
				t.Run(codec.name, func(t *testing.T) {
					encoded, err := codec.marshal(test.message)
					require.NoError(t, err)
					decoded := test.message.ProtoReflect().New().Interface()
					require.NoError(t, codec.unmarshal(encoded, decoded))
					require.True(t, proto.Equal(test.message, decoded), "presence or value lost: %s", encoded)
				})
			}
		})
	}
}

func TestFieldOptionsPreserveExplicitDefaults(t *testing.T) {
	options := &descriptorpb.FieldOptions{}
	proto.SetExtension(options, fino.E_Db, &fino.DBOptions{Default: proto.String(""), Index: proto.String("")})
	proto.SetExtension(options, fino.E_Validate, &fino.ValidateOptions{Min: proto.Int64(0), Max: proto.Int64(0)})
	encoded, err := proto.Marshal(options)
	require.NoError(t, err)
	decoded := &descriptorpb.FieldOptions{}
	require.NoError(t, proto.Unmarshal(encoded, decoded))
	db := proto.GetExtension(decoded, fino.E_Db).(*fino.DBOptions)
	require.NotNil(t, db.Default)
	require.Empty(t, *db.Default)
	require.NotNil(t, db.Index)
	validate := proto.GetExtension(decoded, fino.E_Validate).(*fino.ValidateOptions)
	require.NotNil(t, validate.Min)
	require.Zero(t, *validate.Min)
	require.NotNil(t, validate.Max)
	require.Zero(t, *validate.Max)
}
