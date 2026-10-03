package core

import (
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
		parsed := NewErrorCode(int32(v))
		x.Code = parsed.Code
		x.Name = parsed.Name
		x.Description = parsed.Description
		x.HttpStatusCode = parsed.HttpStatusCode
	}
	return nil
}

func (x *ErrorCode) Format() string {
	if x != nil {
		return strconv.FormatInt(int64(x.Code), 10)
	}
	return ""
}
