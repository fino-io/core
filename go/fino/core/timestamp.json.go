package core

import (
	"encoding/json"
	"io"
	"strconv"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

func init() {
	RegisterJSONTypeDecoder("core.Timestamp", &TimestampCodec{})
	RegisterJSONTypeEncoder("core.Timestamp", &TimestampCodec{})
}

type TimestampCodec struct {
}

func (codec *TimestampCodec) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	a := iter.ReadAny()
	if iter.Error != nil && iter.Error != io.EOF {
		return
	}
	ts := (*Timestamp)(ptr)
	if a.ValueType() == jsoniter.NumberValue {
		if !json.Valid([]byte(a.ToString())) {
			iter.ReportError("Decode Timestamp", "invalid JSON number")
			return
		}
		// Numeric timestamps are Unix seconds, independent of platform int size.
		number, err := strconv.ParseInt(a.ToString(), 10, 64)
		if err != nil {
			iter.ReportError("Decode Timestamp", err.Error())
			return
		}
		ts.Seconds, ts.Nanoseconds = number, 0
	} else if a.ValueType() == jsoniter.StringValue {
		if err := ts.Parse(a.ToString()); err != nil {
			iter.ReportError("Decode Timestamp", err.Error())
		}
	} else if a.ValueType() == jsoniter.NilValue {
		ts.Seconds, ts.Nanoseconds = 0, 0
	} else {
		iter.ReportError("Decode Timestamp", "expected Unix seconds or timestamp string")
	}
}

func (codec *TimestampCodec) IsEmpty(ptr unsafe.Pointer) bool {
	ts := (*Timestamp)(ptr)
	return ts == nil || (ts.Seconds == 0 && ts.Nanoseconds == 0)
}

func (codec *TimestampCodec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	stream.WriteString((*Timestamp)(ptr).Format())
}
