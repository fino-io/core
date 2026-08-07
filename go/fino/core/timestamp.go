package core

import (
	"time"
)

func Now() *Timestamp {
	return FromTime(time.Now())
}

// FromTime covert form time.Time to Timestamp
func FromTime(t time.Time) *Timestamp {
	sec := t.Unix()
	return &Timestamp{
		Seconds:     sec,
		Nanoseconds: int32(t.UnixNano() - sec*1e9),
	}
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
		ft := FromTime(t)
		x.Seconds = ft.Seconds
		x.Nanoseconds = ft.Nanoseconds
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
	if x != nil && u != nil {
		if x != u {
			return x.ToTime().After(u.ToTime())
		}
	}
	return false
}

func (x *Timestamp) Before(u *Timestamp) bool {
	if x != nil && u != nil {
		if x != u {
			return x.ToTime().Before(u.ToTime())
		}
	}
	return false
}

func (x *Timestamp) Equal(u *Timestamp) bool {
	if x != nil && u != nil {
		if x == u {
			return true
		}
		return x.ToTime().Equal(u.ToTime())
	}
	return false
}

func (x *Timestamp) Compare(u *Timestamp) int {
	if x != nil {
		if u != nil {
			return x.ToTime().Compare(u.ToTime())
		} else {
			return 1
		}
	}

	if u != nil {
		return -1
	}
	return 0
}

func (x *Timestamp) Date() *Date {
	if x != nil {
		year, month, day := x.ToTime().Date()
		return &Date{
			Year:  int32(year),
			Month: int32(month),
			Day:   int32(day),
		}
	}
	return nil
}

func (x *Timestamp) Add(d *Duration) *Timestamp {
	if x != nil && d != nil {
		return FromTime(x.ToTime().Add(d.ToDuration()))
	}
	return nil
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
