package core

import (
	"encoding/json"
	"fmt"
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
			iter.ReportError("decode timestamp", "invalid JSON number")
			return
		}
		// Numeric timestamps are Unix seconds, independent of platform int size.
		number, err := strconv.ParseInt(a.ToString(), 10, 64)
		if err != nil {
			iter.ReportError("decode timestamp", err.Error())
			return
		}
		parsed := &Timestamp{Seconds: number}
		if err := parsed.CheckValid(); err != nil {
			iter.ReportError("decode timestamp", err.Error())
			return
		}
		ts.Seconds, ts.Nanoseconds = number, 0
	} else if a.ValueType() == jsoniter.StringValue {
		if err := ts.Parse(a.ToString()); err != nil {
			iter.ReportError("decode timestamp", err.Error())
		}
	} else if a.ValueType() == jsoniter.NilValue {
		ts.Seconds, ts.Nanoseconds = 0, 0
	} else {
		iter.ReportError("decode timestamp", "expected Unix seconds or timestamp string")
	}
}

func (codec *TimestampCodec) IsEmpty(ptr unsafe.Pointer) bool {
	ts := (*Timestamp)(ptr)
	return ts == nil || (ts.Seconds == 0 && ts.Nanoseconds == 0)
}

func (codec *TimestampCodec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	ts := (*Timestamp)(ptr)
	if err := ts.CheckValid(); err != nil {
		stream.Error = err
		return
	}
	stream.WriteVal(ts.Format())
}

func (x *Timestamp) MarshalJSON() ([]byte, error) {
	return marshalJSONCodec(x, &TimestampCodec{})
}

func (x *Timestamp) UnmarshalJSON(data []byte) error {
	if x == nil {
		return fmt.Errorf("timestamp.UnmarshalJSON: nil receiver")
	}
	value, err := unmarshalJSONCodec[Timestamp](data, &TimestampCodec{})
	if err == nil {
		x.Seconds, x.Nanoseconds = value.Seconds, value.Nanoseconds
	}
	return err
}
