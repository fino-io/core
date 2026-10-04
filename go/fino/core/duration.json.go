package core

import (
	"io"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

func init() {
	RegisterJSONTypeDecoder("core.Duration", &DurationCodec{})
	RegisterJSONTypeEncoder("core.Duration", &DurationCodec{})
}

type DurationCodec struct {
}

func (codec *DurationCodec) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	a := iter.ReadAny()
	if iter.Error != nil && iter.Error != io.EOF {
		return
	}
	duration := (*Duration)(ptr)
	if a.ValueType() == jsoniter.NumberValue {
		if err := duration.parseSeconds(a.ToString()); err != nil {
			iter.ReportError("Decode Duration", err.Error())
		}
	} else if a.ValueType() == jsoniter.StringValue {
		str := a.ToString()
		err := duration.Parse(str)
		if err != nil {
			iter.ReportError("Decode Duration", err.Error())
		}
	} else if a.ValueType() == jsoniter.NilValue {
		duration.Seconds, duration.Nanoseconds = 0, 0
	} else {
		iter.ReportError("Decode Duration", "expected seconds or duration string")
	}
}

func (codec *DurationCodec) IsEmpty(ptr unsafe.Pointer) bool {
	duration := (*Duration)(ptr)
	return duration == nil || (duration.Seconds == 0 && duration.Nanoseconds == 0)
}

func (codec *DurationCodec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	duration := (*Duration)(ptr)
	stream.WriteVal(duration.Format())
}
