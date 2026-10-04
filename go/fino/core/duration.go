package core

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"
)

func FromDuration(d time.Duration) *Duration {
	return (&Duration{}).FromDuration(d)
}

func NewDuration(sec float64) (*Duration, error) {
	dur := &Duration{}
	if err := dur.FromSeconds(sec); err != nil {
		return nil, err
	}
	return dur, nil
}

func (x *Duration) FromDuration(d time.Duration) *Duration {
	if x != nil {
		x.Seconds = int64(d / time.Second)
		x.Nanoseconds = int32(d % time.Second)
	}
	return x
}

// FromSeconds rounds to the nearest nanosecond and leaves x unchanged on error.
func (x *Duration) FromSeconds(sec float64) error {
	seconds := new(big.Rat).SetFloat64(sec)
	if seconds == nil {
		return fmt.Errorf("invalid duration seconds: %v", sec)
	}
	return x.setSeconds(seconds)
}

func (x *Duration) parseSeconds(value string) error {
	if x == nil {
		return fmt.Errorf("Duration: nil receiver")
	}
	if !json.Valid([]byte(value)) {
		return fmt.Errorf("invalid duration seconds: %q", value)
	}
	approximate, err := strconv.ParseFloat(value, 64)
	if err != nil || approximate < math.MinInt64 || approximate > float64(math.MaxInt64) {
		return fmt.Errorf("duration seconds out of range: %q", value)
	}
	if math.Abs(approximate) < 0.5/float64(time.Second) {
		x.Seconds, x.Nanoseconds = 0, 0
		return nil
	}
	seconds, ok := new(big.Rat).SetString(value)
	if !ok {
		return fmt.Errorf("invalid duration seconds: %q", value)
	}
	return x.setSeconds(seconds)
}

func (x *Duration) setSeconds(seconds *big.Rat) error {
	if x == nil {
		return fmt.Errorf("Duration: nil receiver")
	}
	seconds.Mul(seconds, big.NewRat(int64(time.Second), 1))
	nanos, remainder := new(big.Int), new(big.Int)
	nanos.QuoRem(seconds.Num(), seconds.Denom(), remainder)
	// Round half away from zero without passing large integers through float64.
	if remainder.Abs(remainder).Lsh(remainder, 1).Cmp(seconds.Denom()) >= 0 {
		nanos.Add(nanos, big.NewInt(int64(seconds.Sign())))
	}
	whole, fraction := new(big.Int), new(big.Int)
	whole.QuoRem(nanos, big.NewInt(int64(time.Second)), fraction)
	if !whole.IsInt64() {
		return fmt.Errorf("duration seconds out of range")
	}
	x.Seconds, x.Nanoseconds = whole.Int64(), int32(fraction.Int64())
	return nil
}

func (x *Duration) totalNanoseconds() *big.Int {
	nanos := new(big.Int).Mul(big.NewInt(x.GetSeconds()), big.NewInt(int64(time.Second)))
	return nanos.Add(nanos, big.NewInt(int64(x.GetNanoseconds())))
}

func (x *Duration) ToDuration() (time.Duration, error) {
	nanos := x.totalNanoseconds()
	if !nanos.IsInt64() {
		return 0, fmt.Errorf("duration exceeds time.Duration range: %s", x.Format())
	}
	return time.Duration(nanos.Int64()), nil
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

func (x *Duration) ToNanoSeconds() (int64, error) {
	duration, err := x.ToDuration()
	return int64(duration), err
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
	return x.totalNanoseconds().Cmp(d.totalNanoseconds())
}
