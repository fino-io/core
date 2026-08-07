package core

import (
	"reflect"
	"testing"
)

func TestNewValue(t *testing.T) {
	tests := []struct {
		name    string
		v       any
		want    *Value
		wantErr bool
	}{
		{name: "nil", v: nil, want: NewNullValue(), wantErr: false},
		{name: "bool-true", v: true, want: NewBoolValue(true), wantErr: false},
		{name: "bool-false", v: false, want: NewBoolValue(false), wantErr: false},
		{name: "int", v: 13, want: NewIntValue(13), wantErr: false},
		{name: "int8", v: int8(13), want: NewInt8Value(13), wantErr: false},
		{name: "int16", v: int16(13), want: NewInt16Value(13), wantErr: false},
		{name: "int32", v: int32(13), want: NewInt32Value(13), wantErr: false},
		{name: "int64", v: int64(13), want: NewInt64Value(13), wantErr: false},
		{name: "uint", v: uint(13), want: NewUintValue(13), wantErr: false},
		{name: "uint8", v: uint8(13), want: NewUint8Value(13), wantErr: false},
		{name: "uint16", v: uint16(13), want: NewUint16Value(13), wantErr: false},
		{name: "uint32", v: uint32(13), want: NewUint32Value(13), wantErr: false},
		{name: "uint64", v: uint64(13), want: NewUint64Value(13), wantErr: false},
		{name: "float32", v: float32(13.1), want: NewFloat32Value(13.1), wantErr: false},
		{name: "float64", v: 13.1, want: NewFloat64Value(13.1), wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewValue(tt.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewValue() got = %v, want %v", got, tt.want)
			}
		})
	}
}
