package core

import (
	"cmp"
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
		x.Nanoseconds = int64ToInt32(int64(nsec))
	}
	return x
}

func (x *Duration) FromSeconds(sec float64) *Duration {
	if x != nil {
		x.Seconds = int64(sec)
		delta := sec - float64(x.Seconds)
		x.Nanoseconds = int32(math.Round(delta * float64(time.Second)))
		if x.Nanoseconds >= int32(time.Second) || x.Nanoseconds <= -int32(time.Second) {
			x.Seconds += int64(x.Nanoseconds) / int64(time.Second)
			x.Nanoseconds %= int32(time.Second)
		}
	}
	return x
}

func (x *Duration) ToDuration() time.Duration {
	return time.Duration(x.GetSeconds())*time.Second + time.Duration(x.GetNanoseconds())
}

func (x *Duration) ToHours() float64 {
	return x.ToSeconds() / 3600
}

func (x *Duration) ToMinutes() float64 {
	return x.ToSeconds() / 60
}

func (x *Duration) ToSeconds() float64 {
	return float64(x.GetSeconds()) + float64(x.GetNanoseconds())/float64(time.Second)
}

func (x *Duration) ToNanoSeconds() int64 {
	if x != nil {
		return x.ToDuration().Nanoseconds()
	}
	return 0
}

func (x *Duration) Compare(d *Duration) int {
	if x == nil || d == nil {
		if x == d {
			return 0
		}
		if x == nil {
			return -1
		}
		return 1
	}
	if order := cmp.Compare(x.Seconds, d.Seconds); order != 0 {
		return order
	}
	return cmp.Compare(x.Nanoseconds, d.Nanoseconds)
}
