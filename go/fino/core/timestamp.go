package core

import (
	"fmt"
	"time"
)

func Now() *Timestamp {
	return FromTime(time.Now())
}

// FromTime converts time.Time to Timestamp.
func FromTime(t time.Time) *Timestamp {
	return (&Timestamp{}).FromTime(t)
}

func Since(t *Timestamp) *Duration {
	if t != nil {
		return FromDuration(time.Since(t.ToTime()))
	}
	return nil
}

func Until(t *Timestamp) *Duration {
	if t != nil {
		return FromDuration(time.Until(t.ToTime()))
	}
	return nil
}

func (x *Timestamp) FromTime(t time.Time) *Timestamp {
	if x != nil {
		x.Seconds = t.Unix()
		x.Nanoseconds = int32(t.Nanosecond())
	}
	return x
}

func (x *Timestamp) ToTime() time.Time {
	if x != nil {
		return time.Unix(x.Seconds, int64(x.Nanoseconds))
	}
	return time.Time{}
}

func (x *Timestamp) After(u *Timestamp) bool {
	return x != nil && u != nil && x.ToTime().After(u.ToTime())
}

func (x *Timestamp) Before(u *Timestamp) bool {
	return x != nil && u != nil && x.ToTime().Before(u.ToTime())
}

func (x *Timestamp) Equal(u *Timestamp) bool {
	return x != nil && u != nil && x.ToTime().Equal(u.ToTime())
}

func (x *Timestamp) Compare(u *Timestamp) int {
	if x == u {
		return 0
	}
	if x == nil {
		return -1
	}
	if u == nil {
		return 1
	}
	return x.ToTime().Compare(u.ToTime())
}

func (x *Timestamp) Date() *Date {
	if x != nil {
		year, month, day := x.ToTime().Date()
		return &Date{
			Year:  int64ToInt32(int64(year)),
			Month: int64ToInt32(int64(month)),
			Day:   int64ToInt32(int64(day)),
		}
	}
	return nil
}

func (x *Timestamp) Add(d *Duration) (*Timestamp, error) {
	if x == nil || d == nil {
		return nil, fmt.Errorf("timestamp.Add: nil timestamp or duration")
	}
	if err := x.CheckValid(); err != nil {
		return nil, err
	}
	duration, err := d.ToDuration()
	if err != nil {
		return nil, err
	}
	result := FromTime(x.ToTime().Add(duration))
	if err := result.CheckValid(); err != nil {
		return nil, err
	}
	return result, nil
}

func (x *Timestamp) AddDate(year, month, day int) *Timestamp {
	if x != nil {
		return FromTime(x.ToTime().AddDate(year, month, day))
	}
	return nil
}

func (x *Timestamp) Sub(u *Timestamp) *Duration {
	if x != nil {
		return FromDuration(x.ToTime().Sub(u.ToTime()))
	}
	return nil
}
