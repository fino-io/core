package core

import (
	"fmt"
	"strconv"

	"google.golang.org/protobuf/proto"
)

func NewErrorCode(code int32) *ErrorCode {
	if ec, ok := errorCodeIndex[code]; ok {
		return proto.Clone(ec).(*ErrorCode)
	}
	return &ErrorCode{Code: code}
}

func ParseErrorCode(code string) (*ErrorCode, error) {
	ec := &ErrorCode{}
	err := ec.Parse(code)
	if err != nil {
		return nil, err
	}
	return ec, nil
}

func (x *ErrorCode) Parse(code string) error {
	if x == nil {
		return fmt.Errorf("ErrorCode.Parse: nil receiver")
	}
	v, err := strconv.ParseInt(code, 10, 32)
	if err != nil {
		return fmt.Errorf("failed to parse error code %w", err)
	}
	parsed := NewErrorCode(int32(v))
	proto.Reset(x)
	proto.Merge(x, parsed)
	return nil
}

func (x *ErrorCode) Format() string {
	if x != nil {
		return strconv.FormatInt(int64(x.Code), 10)
	}
	return ""
}
