package core

import (
	"bytes"
	"fmt"
	"strconv"
)

func NewErrorCode(code int32) *ErrorCode {
	if ec, ok := errorCodeIndex[code]; ok {
		return &ErrorCode{
			Code:           ec.Code,
			Name:           ec.Name,
			Description:    ec.Description,
			HttpStatusCode: ec.HttpStatusCode,
		}
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
	if x != nil && len(code) > 0 {
		v, err := strconv.ParseInt(code, 10, 32)
		if err != nil {
			return fmt.Errorf("failed to parse error code %w", err)
		}
		x.Code = int32(v)
	}
	return nil
}

func (x *ErrorCode) Format() string {
	if x != nil {
		buffer := bytes.Buffer{}
		buffer.WriteString(strconv.FormatInt(int64(x.Code), 10))
		return buffer.String()
	}
	return ""
}
