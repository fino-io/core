package core

import (
	"math"
	"time"
)

func FromDuration(d time.Duration) *Duration {
	dur := &Duration{}
	return dur.FromDuration(d)
}

func NewDuration(sec float64) *Duration {
	dur := &Duration{}
	return dur.FromSeconds(sec)
}

func (x *Duration) FromDuration(d time.Duration) *Duration {
	if x != nil {
		sec := d / time.Second
		nsec := d % time.Second

		x.Seconds = int64(sec)
		x.Nanoseconds = int32(nsec)
	}
	return x
}

func (x *Duration) FromSeconds(sec float64) *Duration {
	if x != nil {
		x.Seconds = int64(sec)
		delta := sec - float64(x.Seconds)
		x.Nanoseconds = int32(math.Round(delta * float64(time.Second)))
	}
	return x
}

func (x *Duration) ToDuration() time.Duration {
	return time.Duration(x.Seconds)*time.Second + time.Duration(x.Nanoseconds)*time.Nanosecond
}

func (x *Duration) ToHours() float64 {
	if x != nil {
		return x.ToDuration().Hours()
	}
	return 0
}

func (x *Duration) ToMinutes() float64 {
	if x != nil {
		return x.ToDuration().Minutes()
	}
	return 0
}

func (x *Duration) ToSeconds() float64 {
	if x != nil {
		return x.ToDuration().Seconds()
	}
	return 0
}

func (x *Duration) ToNanoSeconds() int64 {
	if x != nil {
		return x.ToDuration().Nanoseconds()
	}
	return 0
}

func (x *Duration) Compare(d *Duration) int {
	if x != nil {
		if d != nil {
			if x.Seconds == d.Seconds {
				if x.Nanoseconds == d.Nanoseconds {
					return 0
				} else if x.Nanoseconds > d.Nanoseconds {
					return 1
				} else {
					return -1
				}
			} else if x.Seconds > d.Seconds {
				return 1
			}
			return -1
		} else {
			return 1
		}
	} else if d != nil {
		return -1
	}
	return 0
}
