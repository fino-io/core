package core

import (
	"cmp"
	"errors"
	"fmt"
	"math"
	"strconv"
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

// FromSeconds sets a finite, representable number of seconds.
// Use Scan when accepting external input that needs validation.
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

func (x *Duration) setSeconds(seconds float64) error {
	// The upper bound is exclusive: float64 rounds MaxInt64 up to 2^63.
	// These comparisons also reject NaN and infinities.
	if !(seconds >= math.MinInt64 && seconds < -float64(math.MinInt64)) {
		return fmt.Errorf("duration seconds out of range: %v", seconds)
	}
	x.FromSeconds(seconds)
	return nil
}

func (x *Duration) parseSeconds(value string) error {
	seconds, err := strconv.ParseInt(value, 10, 64)
	if err == nil {
		x.Seconds, x.Nanoseconds = seconds, 0
		return nil
	}
	if errors.Is(err, strconv.ErrRange) {
		return err
	}
	number, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}
	return x.setSeconds(number)
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
